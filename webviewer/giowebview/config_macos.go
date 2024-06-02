//go:build darwin && !ios

package giowebview

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) webview.Config {
	evt, ok := e.(app.AppKitViewEvent)
	if !ok {
		return webview.Config{}
	}

	return webview.Config{View: evt.View, Layer: evt.Layer, RunOnMain: w.Run}
}
