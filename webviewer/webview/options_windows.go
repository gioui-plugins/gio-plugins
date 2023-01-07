package webview

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"os"
	"strings"

	"golang.org/x/sys/windows"
)

func (r *driver) setProxy() {
	options.Lock()
	defer options.Unlock()

	if options.proxy.ip == "" && options.proxy.port == "" {
		return
	}

	proxy := fmt.Sprintf("%s:%s", options.proxy.ip, options.proxy.port)
	if strings.Index(options.proxy.ip, `:`) >= 0 && !strings.HasPrefix(options.proxy.ip, `[`) && !strings.HasPrefix(options.proxy.ip, `]`) {
		proxy = fmt.Sprintf("[%s]:%s", options.proxy.ip, options.proxy.port)
	}

	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", fmt.Sprintf(`%s --proxy-server="%s"`, os.Getenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS"), proxy))
}

func (r *driver) setCerts() {
	options.Lock()
	defer options.Unlock()

	if options.certs == nil {
		return
	}

	var jcerts string
	h := sha256.New()
	for _, c := range options.certs {
		h.Write(c.RawSubjectPublicKeyInfo)
		jcerts += base64.StdEncoding.EncodeToString(h.Sum(nil)) + ","
		h.Reset()
	}

	os.Setenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", fmt.Sprintf(`%s --ignore-certificate-errors-spki-list="%s"`, os.Getenv("WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS"), jcerts))
}

func (r *driver) setDir() {
	options.Lock()
	defer options.Unlock()

	if len(options.folder) == 0 {
		return
	}

	f, err := windows.UTF16FromString(options.folder)
	if err != nil {
		return
	}
	r.dir = f
}
