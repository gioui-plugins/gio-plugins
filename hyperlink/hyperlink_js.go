//go:build js

package hyperlink

import (
	"net/url"
	"syscall/js"
)

var (
	_document = js.Global().Get("document")
	_body     = js.Global().Get("document").Get("body")
	_actives  = make([]js.Value, 0, 32)
)

type driver struct {
	closeFn js.Func
}

func attachDriver(house *Hyperlink, config Config) {
	d := driver{}
	d.closeFn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		configureDriver(&d, Config{Blur: true})
		return nil
	})

	js.Global().Get("document").Call("addEventListener", "visibilitychange", d.closeFn)

	house.driver = d
}

func configureDriver(driver *driver, config Config) {
	if config.Blur {
		for i := 0; i < len(_actives); i++ {
			_body.Call("removeChild", _actives[i])
		}
		_actives = _actives[:0]
	}
}

func (d *driver) closeFunc() js.Func {
	return d.closeFn
}

func (d *driver) open(u *url.URL) error {
	if ok := js.Global().Call("open", u.String(), "_blank", "noreferrer,noopener").Truthy(); !ok {
		// If there's a error let's use the hacky way:
		// It will create a "fullscreen <a>", which clicking will
		// open the URL.
		// Generally, it will need two clicks to open the URL.

		// We remove this `a` when the app lost focus (based on Page Visibility API, which Gio relies on),
		// or when the user clicks on the `a` element, if the `onclick` works.
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
		a.Get("style").Set("z-index", "2147483647")
		a.Set("onclick", d.closeFunc())
		_body.Call("appendChild", a)
		_actives = append(_actives, a)
	}

	return nil
}
