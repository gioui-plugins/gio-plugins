//go:build darwin && !ios
// +build darwin,!ios

package hyperlink

import (
	"net/url"
	"os/exec"
)

type driver struct{}

func attachDriver(house *Hyperlink, config Config) {}

func configureDriver(driver *driver, config Config) {}

func (*driver) open(u *url.URL) error {
	return exec.Command("open", u.String()).Run()
}
