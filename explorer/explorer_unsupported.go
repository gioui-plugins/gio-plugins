// SPDX-License-Identifier: Unlicense OR MIT

//go:build !windows && !android && !js && !darwin && !ios
// +build !windows,!android,!js,!darwin,!ios

package explorer

import (
	"io"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type explorer struct{}

func newExplorer(w *app.Window) *explorer {
	return new(explorer)
}

func (e *explorerPlugin) listenEvents(_ event.Event) {}

func (e *explorerPlugin) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	return nil, ErrNotAvailable
}

func (e *explorerPlugin) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	return nil, ErrNotAvailable
}
