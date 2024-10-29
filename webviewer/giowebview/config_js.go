package giowebview

import (
	"syscall/js"

	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
	"github.com/inkeliz/go_inkwasm/inkwasm"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) webview.Config {
	return webview.Config{Element: inkwasm.NewObjectFromSyscall(js.Global().Get("document").Get("body")), RunOnMain: w.Run}
}

// UpdateConfigFromViewEvent updates a webview.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *webview.Config, w *app.Window, e app.ViewEvent) {
	config.Element = inkwasm.NewObjectFromSyscall(js.Global().Get("document").Get("body"))
	config.RunOnMain = w.Run
}

// UpdateConfigFromFrameEvent updates a webview.Config based on app.FrameEvent.
func UpdateConfigFromFrameEvent(config *webview.Config, w *app.Window, evt app.FrameEvent) {
	config.RunOnMain = w.Run
	config.PxPerDp = evt.Metric.PxPerDp
}
