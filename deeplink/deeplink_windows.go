package deeplink

import (
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"golang.org/x/sys/windows/registry"
	"log"
	"net"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"unsafe"
)

// That is process on init, to match the behaviour of the other platforms, that is not necessary on Windows,
// we could have something similar to `flag` package (i.e. flag.Parse()), but I think it's better to force
// the same behaviour across all platforms.
func init() {
	if err := checkLinkOnStartup(); err != nil {
		_Deeplink.updateError(err)
	}
}

func init() {
	if err := listenScheme(); err != nil {
		_Deeplink.updateError(err)
		return
	}

	installSchemes()
}

func checkLinkOnStartup() error {
	if len(os.Args) <= 1 {
		return nil
	}

	// Extract the URL from the command-line argument
	u := os.Args[1]

	uri, err := url.Parse(u)
	if err != nil {
		return nil
	}

	// Check if the URL scheme is supported
	found := false
	for _, v := range schemeList {
		if v == uri.Scheme {
			found = true
		}
	}

	if !found {
		return nil
	}

	if isProcessRunning() {
		if err := broadcastScheme(u); err == nil {
			// Was able to send the message to the running instance, so exit here.
			os.Exit(0)
			return nil
		}
	}

	_Deeplink.updateURL(u)
	return nil
}

func isProcessRunning() bool {
	handle, err := syscall.CreateToolhelp32Snapshot(syscall.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		log.Fatal(err)
		return false
	}
	defer syscall.CloseHandle(handle)

	var entry syscall.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	err = syscall.Process32First(handle, &entry)
	if err != nil {
		return false
	}

	executablePath, err := os.Executable()
	if err != nil {
		return false
	}

	executableBase := filepath.Base(executablePath)
	for {
		if syscall.UTF16ToString(entry.ExeFile[:]) == executableBase {
			return true
		}

		err = syscall.Process32Next(handle, &entry)
		if err != nil {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}
			return false
		}
	}

	return false
}

func socketPathCurrent() string {
	p := os.Getpid()
	return socketPath() + "-" + strconv.Itoa(p)
}

func socketPath() string {
	b := sha512.Sum512_224([]byte(schemes))
	return filepath.Join(os.TempDir(), "deeplink-"+string(base64.RawURLEncoding.EncodeToString(b[:])))
}

func broadcastScheme(scheme string) error {
	matches, err := filepath.Glob(socketPath() + "*")
	if err != nil {
		return err
	}

	done := false
	for _, v := range matches {
		client, err := net.Dial("unix", v)
		if err != nil {
			continue
		}

		success := true
		if _, err = client.Write([]byte(scheme)); err != nil {
			success = false
		}

		if err := client.Close(); err != nil {
			success = false
		}

		if success {
			done = true
		}
	}

	if done {
		return nil
	}
	return errors.New("no running instance found")
}

func listenScheme() error {
	// Start a UNIX domain socket server to receive messages from subsequent instances
	ln, err := net.Listen("unix", socketPathCurrent())
	if err != nil {
		return err
	}

	// Handle messages received from subsequent instances
	go func() {
		buffer := make([]byte, 65536)
		for {
			conn, err := ln.Accept()
			if err != nil {
				continue
			}
			receiveScheme(conn, buffer)
		}
		ln.Close()
	}()

	return nil
}

func receiveScheme(conn net.Conn, buffer []byte) {
	defer conn.Close()

	// Read the received message (deep link URL) from the connection
	n, err := conn.Read(buffer)
	if err != nil {
		return
	}

	// Extract the deep link URL from the received message
	_Deeplink.updateURL(string(buffer[:n]))
}

func installSchemes() {
	for _, v := range schemeList {
		installScheme(v)
	}
}

func installScheme(scheme string) error {
	path, err := os.Executable()
	if err != nil {
		return err
	}

	key, _, err := registry.CreateKey(registry.CURRENT_USER, `Software\\Classes\\`+scheme, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer key.Close()

	err = key.SetStringValue("", "URL:"+scheme+" Protocol")
	if err != nil {
		return err
	}

	err = key.SetStringValue("URL Protocol", "")
	if err != nil {
		return err
	}

	defaultIconKey, _, err := registry.CreateKey(key, `DefaultIcon`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer defaultIconKey.Close()

	err = defaultIconKey.SetStringValue("", `"`+path+`",1`)
	if err != nil {
		return err
	}

	shellOpenCommandKey, _, err := registry.CreateKey(key, `shell\\open\\command`, registry.ALL_ACCESS)
	if err != nil {
		return err
	}
	defer shellOpenCommandKey.Close()

	err = shellOpenCommandKey.SetStringValue("", `"`+path+`" "%1"`)
	if err != nil {
		return err
	}

	return nil
}
