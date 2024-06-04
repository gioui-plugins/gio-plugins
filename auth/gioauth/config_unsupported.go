//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package gioauth

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/auth"
)

// NewConfigFromViewEvent creates an auth.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent, title string) auth.Config {
	return auth.Config{}
}

// UpdateConfigFromViewEvent updates an auth.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *auth.Config, w *app.Window, evt app.ViewEvent) {}
