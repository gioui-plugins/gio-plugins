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
	config.HWND = evt.(app.Win32ViewEvent).HWND
	config.RunOnMain = w.Run
}
