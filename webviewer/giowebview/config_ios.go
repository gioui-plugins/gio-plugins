package giowebview

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) webview.Config {
	evt, ok := e.(app.UIKitViewEvent)
	if !ok {
		return webview.Config{}
	}

	return webview.Config{View: evt.ViewController, RunOnMain: w.Run}
}
