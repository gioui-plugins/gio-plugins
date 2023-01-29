package webviewer

import (
	"gioui.org/app"
	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) webview.Config {
	return webview.Config{View: jni.Class(evt.View), VM: jni.JVMFor(app.JavaVM()), Context: jni.Object(app.AppContext()), RunOnMain: w.Run}
}
