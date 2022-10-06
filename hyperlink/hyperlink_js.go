//go:build js
// +build js

package hyperlink

import (
	"gioui.org/io/event"
	"gioui.org/io/system"
	"net/url"
	"syscall/js"
)

var (
	_document = js.Global().Get("document")
	_body     = js.Global().Get("document").Get("body")
)

type hyperlink struct{}

func (*hyperlinkPlugin) listenEvents(event event.Event) {
	if _, ok := event.(system.StageEvent); ok {
		links := _body.Call("querySelectorAll", "a.hyperlink")
		if !links.Truthy() {
			return
		}
		for i := 0; i < links.Length(); i++ {
			_body.Call("removeChild", links.Index(0))
		}
	}
}

func (*hyperlink) open(u *url.URL) error {
	if ok := js.Global().Call("open", u.String(), "_blank", "noreferrer,noopener").Truthy(); !ok {
		// If there's a error let's use the hacky way:
		// It will create a "fullscreen <a>", which clicking will
		// open the URL.
		// Generally, it will need two clicks to open the URL.

		// We can't hook into `a` (adding `a.addEvenetListener("click")` will make it fail again,
		// not sure why.
		// We remove this `a` when the app lost focus (based on Page Visibility API, which Gio relies on).
		a := _document.Call("createElement", "a")
		a.Set("href", u.String())
		a.Set("target", "_blank")
		a.Set("rel", "noreferrer,noopener")
		a.Set("innerText", " ")
		a.Get("classList").Call("add", "hyperlink")
		a.Get("style").Set("display", "block")
		a.Get("style").Set("width", "100vw")
		a.Get("style").Set("height", "100vh")
		a.Get("style").Set("position", "fixed")
		a.Get("style").Set("top", "0")
		a.Get("style").Set("z-index", "100")
		_body.Call("appendChild", a)
	}

	return nil
}
