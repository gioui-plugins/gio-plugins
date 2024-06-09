package giowebview

import (
	"gioui.org/app"
	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) webview.Config {
	evt, ok := e.(app.AndroidViewEvent)
	if !ok {
		return webview.Config{}
	}

	return webview.Config{View: jni.Class(evt.View), VM: jni.JVMFor(app.JavaVM()), Context: jni.Object(app.AppContext()), RunOnMain: w.Run}
}
