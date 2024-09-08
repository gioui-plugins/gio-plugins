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

// UpdateConfigFromViewEvent updates a webview.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *webview.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.AndroidViewEvent)
	if !ok {
		return
	}

	config.View = jni.Class(evt.View)
	config.VM = jni.JVMFor(app.JavaVM())
	config.Context = jni.Object(app.AppContext())
	config.RunOnMain = w.Run
}

// UpdateConfigFromFrameEvent updates a webview.Config based on app.FrameEvent.
func UpdateConfigFromFrameEvent(config *webview.Config, w *app.Window, evt app.FrameEvent) {
	config.RunOnMain = func(f func()) {
		w.Run(f)
	}
	config.PxPerDp = evt.Metric.PxPerDp
}
