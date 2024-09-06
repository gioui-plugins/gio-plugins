//go:build (linux && !android) || openbsd || freebsd || netbsd || dragonfly

package hyperlink

import (
	"net/url"
	"os/exec"
)

type driver struct{}

func attachDriver(house *Hyperlink, config Config) {}

func configureDriver(driver *driver, config Config) {}

func (*driver) open(u *url.URL, preferredPackage string) error {
	return exec.Command("xdg-open", u.String()).Run()
}
