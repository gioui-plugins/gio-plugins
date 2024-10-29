//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package giowebview

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

// NewConfigFromViewEvent creates a webview.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) webview.Config {
	return webview.Config{}
}

// UpdateConfigFromViewEvent updates a webview.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *webview.Config, w *app.Window, e app.ViewEvent) {}

// UpdateConfigFromFrameEvent updates a webview.Config based on app.FrameEvent.
func UpdateConfigFromFrameEvent(config *webview.Config, w *app.Window, evt app.FrameEvent) {}
