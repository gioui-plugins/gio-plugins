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

// UpdateConfigFromViewEvent updates a webview.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *webview.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.UIKitViewEvent)
	if !ok {
		return
	}

	config.View = evt.ViewController
	config.RunOnMain = w.Run
}

// UpdateConfigFromFrameEvent updates a webview.Config based on app.FrameEvent.
func UpdateConfigFromFrameEvent(config *webview.Config, w *app.Window, evt app.FrameEvent) {
	config.RunOnMain = w.Run
	config.PxPerDp = evt.Metric.PxPerDp
}
