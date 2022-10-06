//go:build darwin && !ios

package webviewer

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) webview.Config {
	return webview.Config{View: evt.View, Layer: evt.Layer, RunOnMain: w.Run}
}
