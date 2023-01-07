package webview

import (
	"fmt"
	"net/url"

	"github.com/inkeliz/go_inkwasm/inkwasm"
)

type driver struct {
	config Config
	iframe inkwasm.Object
}

func (r *driver) attach(w *webview) error {
	defer w.scheduler.SetRunner(r.config.RunOnMain)

	r.iframe = createElement("iframe")
	setStyleDisplay(r.iframe, "none")
	setStylePosition(r.iframe, "fixed")
	setStyleZIndex(r.iframe, "1")
	setStyleBorder(r.iframe, "none")
	prepend(r.config.Element, r.iframe)

	return nil
}

func (r *driver) configure(w *webview, config Config) {
	r.config = config
}

func (r *driver) resize(w *webview, pos [4]float32) {
	scale := r.config.PxPerDp

	if pos[2] == 0 && pos[3] == 0 {
		if w.visible {
			setStyleDisplay(r.iframe, "none")
			w.visible = false
		}
	} else {
		pos := [4]int{int((pos[0] + 0.5) / scale), int((pos[1] + 0.5) / scale), int((pos[2] + 0.5) / scale), int((pos[3] + 0.5) / scale)}
		for i, v := range []func(object inkwasm.Object, v string){setStyleLeft, setStyleTop, setStyleWidth, setStyleHeight} {
			v(r.iframe, fmt.Sprintf("%dpx", pos[i]))
		}

		if !w.visible {
			setStyleDisplay(r.iframe, "block")
			w.visible = true
		}
	}
}

func (r *driver) navigate(w *webview, url *url.URL) {
	setSrc(r.iframe, url.String())
}

func (r *driver) close(w *webview) {
	removeChild(r.config.Element, r.iframe)
}
