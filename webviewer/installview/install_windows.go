package installview

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var DownloadURL = `https://go.microsoft.com/fwlink/p/?LinkId=2124703`

func (i *Installer) install() error {
	if i.Data == nil {
		return fmt.Errorf("no data")
	}

	randName := make([]byte, 32)
	rand.Read(randName)

	path := filepath.Join(os.TempDir(), hex.EncodeToString(randName)+".exe")

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { os.Remove(path) }()

	if _, err := file.Write(i.Data); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	return exec.Command("powershell.exe", "-Command", "Start-Process cmd '/c "+path+" /install' -WindowStyle hidden -Wait").Run()
}
