package gioauth

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/auth"
)

// NewConfigFromViewEvent creates a auth.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) auth.Config {
	r := auth.Config{}
	UpdateConfigFromViewEvent(&r, w, e)
	return r
}

// UpdateConfigFromViewEvent updates a auth.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *auth.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.UIKitViewEvent)
	if !ok {
		return
	}

	config.View = evt.ViewController
	config.RunOnMain = w.Run
}
