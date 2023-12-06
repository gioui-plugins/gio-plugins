package gioauth

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/auth"
)

// NewConfigFromViewEvent creates an auth.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) auth.Config {
	r := auth.Config{}
	UpdateConfigFromViewEvent(&r, w, evt)
	return r
}

// UpdateConfigFromViewEvent updates an auth.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *auth.Config, w *app.Window, evt app.ViewEvent) {}
