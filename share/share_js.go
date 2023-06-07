// SPDX-License-Identifier: Unlicense OR MIT

package share

import (
	"syscall/js"
)

type driver struct{}

func attachDriver(house *Share, config Config) {
	house.driver = driver{}
}

func configureDriver(driver *driver, config Config) {}

func (e *driver) shareText(title, text string) error {
	obj := js.Global().Get("Object").New()
	obj.Set("title", title)
	obj.Set("text", text)
	return e.showDialog(obj)
}

func (e *driver) shareWebsite(title, description, url string) error {
	obj := js.Global().Get("Object").New()
	obj.Set("title", title)
	obj.Set("text", description)
	obj.Set("url", url)
	return e.showDialog(obj)
}

func (e *driver) showDialog(obj js.Value) error {
	navigator := js.Global().Get("navigator")
	if !navigator.Get("share").Truthy() {
		return ErrNotAvailable
	}
	if !navigator.Call("share", obj).Truthy() {
		return ErrNotAvailable
	}
	return nil
}
