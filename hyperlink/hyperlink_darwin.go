//go:build darwin && !ios
// +build darwin,!ios

package hyperlink

import (
	"gioui.org/io/event"
	"net/url"
	"os/exec"
)

type hyperlink struct{}

func (*hyperlinkPlugin) listenEvents(_ event.Event) {}

func (*hyperlinkPlugin) open(u *url.URL) error {
	return exec.Command("open", u.String()).Run()
}
