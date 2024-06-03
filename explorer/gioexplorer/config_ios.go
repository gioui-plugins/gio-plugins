package gioexplorer

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/explorer"
)

// NewConfigFromViewEvent creates a explorer.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) explorer.Config {
	r := explorer.Config{}
	UpdateConfigFromViewEvent(&r, w, e)
	return r
}

// UpdateConfigFromViewEvent updates explorer.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *explorer.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.UIKitViewEvent)
	if !ok {
		return
	}

	config.View = evt.ViewController
	config.RunOnMain = w.Run
}
