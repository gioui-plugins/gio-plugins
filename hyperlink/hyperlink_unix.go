//go:build (linux && !android) || openbsd || freebsd || netbsd || dragonfly
// +build linux,!android openbsd freebsd netbsd dragonfly

package hyperlink

import (
	"net/url"
	"os/exec"

	"gioui.org/io/event"
)

type hyperlink struct{}

func (*hyperlinkPlugin) listenEvents(_ event.Event) {}

func (*hyperlinkPlugin) open(u *url.URL) error {
	return exec.Command("xdg-open", u.String()).Run()
}
