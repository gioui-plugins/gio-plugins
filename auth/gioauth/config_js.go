package gioauth

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/hyperlink/giohyperlink"
)

// NewConfigFromViewEvent creates an auth.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) auth.Config {
	r := auth.Config{}
	UpdateConfigFromViewEvent(&r, w, e)
	return r
}

// UpdateConfigFromViewEvent updates an auth.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *auth.Config, w *app.Window, e app.ViewEvent) {
	config.Config = giohyperlink.NewConfigFromViewEvent(w, e)
}
