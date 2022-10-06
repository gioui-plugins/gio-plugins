// SPDX-License-Identifier: Unlicense OR MIT

package share

import (
	"syscall/js"

	"gioui.org/app"
	"gioui.org/io/event"
)

type share struct{}

func newShare(w *app.Window) share {
	return share{}
}

func (e *sharePlugin) listenEvents(_ event.Event) {}

func (e *sharePlugin) shareText(op TextOp) error {
	obj := js.Global().Get("Object").New()
	obj.Set("text", op.Text)
	obj.Set("title", op.Title)
	return e.showDialog(obj)
}

func (e *sharePlugin) shareWebsite(op WebsiteOp) error {
	obj := js.Global().Get("Object").New()
	obj.Set("text", op.Text)
	obj.Set("title", op.Title)
	obj.Set("url", op.Link)
	return e.showDialog(obj)
}

func (e *sharePlugin) showDialog(obj js.Value) error {
	navigator := js.Global().Get("navigator")
	if !navigator.Get("share").Truthy() {
		return ErrNotAvailable
	}
	if !navigator.Call("share", obj).Truthy() {
		return ErrNotAvailable
	}
	return nil
}
