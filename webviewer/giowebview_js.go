package webviewer

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
