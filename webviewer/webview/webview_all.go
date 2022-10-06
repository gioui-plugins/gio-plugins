//go:build ios || darwin || windows || js || android

package webview

import (
	"net/url"
	"sync"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

// webview implements the WebView interface.
// The driver is different for each platform.
type webview struct {
	handle internal.Handle
	driver *driver

	mutex     sync.Mutex
	scheduler internal.Scheduler
	fan       internal.Fan[Event]

	lastPos [4]float32
	visible bool
	closed  bool

	javascriptManager JavascriptManager
	dataManager       DataManager
}

func newWebview(config Config) (*webview, error) {
	w := &webview{driver: &driver{config: config}}
	w.handle = internal.NewHandle(w)

	if err := w.driver.attach(w); err != nil {
		return nil, err
	}

	return w, nil
}

// Configure implements the WebView interface.
func (w *webview) Configure(config Config) {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.closed {
		return
	}

	w.driver.configure(w, config)
}

// Resize implements the WebView interface.
func (w *webview) Resize(size Point, offset Point) {
	pos := [4]float32{offset.X, offset.Y, size.X, size.Y}
	if w.lastPos == pos {
		return
	}
	w.lastPos = pos

	w.scheduler.Run(idMethodResize, func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		if w.closed {
			return
		}

		w.driver.resize(w, pos)
	})
}

// Navigate implements the WebView interface.
func (w *webview) Navigate(url *url.URL) {
	w.scheduler.Run(idMethodNavigate, func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		if w.closed {
			return
		}

		w.driver.navigate(w, url)
	})
}

// Close implements the WebView interface.
func (w *webview) Close() {
	w.scheduler.Run(idMethodClose, func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()
		if w.closed {
			return
		}

		w.driver.close(w)
		w.fan.Close()
		w.closed = true
	})
}

// Events implements the WebView interface.
func (w *webview) Events() chan Event { return w.fan.Add() }

// DataManager implements the WebView interface.
func (w *webview) DataManager() DataManager { return w.dataManager }

// JavascriptManager implements the WebView interface.
func (w *webview) JavascriptManager() JavascriptManager { return w.javascriptManager }
