//go:build windows
// +build windows

package hyperlink

import (
	"net/url"

	"golang.org/x/sys/windows"
)

type driver struct{}

func attachDriver(house *Hyperlink, config Config) {}

func configureDriver(driver *driver, config Config) {}

func (*driver) open(u *url.URL) error {
	return windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(u.String()), nil, nil, windows.SW_SHOWNORMAL)
}
