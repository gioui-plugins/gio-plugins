//go:build darwin && !ios

package gioexplorer

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/explorer"
)

// NewConfigFromViewEvent creates a explorer.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) explorer.Config {
	r := explorer.Config{}
	UpdateConfigFromViewEvent(&r, w, evt)
	return r
}

func UpdateConfigFromViewEvent(config *explorer.Config, w *app.Window, evt app.ViewEvent) {
	config.View = evt.View
	config.Layer = evt.Layer
	config.RunOnMain = w.Run
}
