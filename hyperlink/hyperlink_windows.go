//go:build windows
// +build windows

package hyperlink

import (
	"gioui.org/io/event"
	"golang.org/x/sys/windows"
	"net/url"
)

type hyperlink struct{}

func (*hyperlinkPlugin) listenEvents(_ event.Event) {}

func (*hyperlinkPlugin) open(u *url.URL) error {
	return windows.ShellExecute(0, nil, windows.StringToUTF16Ptr(u.String()), nil, nil, windows.SW_SHOWNORMAL)
}
