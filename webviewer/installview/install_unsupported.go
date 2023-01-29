//go:build !windows

package installview

var DownloadURL = ""

func (i *Installer) install() error {
	return nil
}
