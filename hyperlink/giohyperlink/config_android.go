package giohyperlink

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/hyperlink"
)

// NewConfigFromViewEvent creates a share.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) hyperlink.Config {
	r := hyperlink.Config{}
	UpdateConfigFromViewEvent(&r, w, evt)
	return r
}

func UpdateConfigFromViewEvent(config *hyperlink.Config, w *app.Window, evt app.ViewEvent) {
	config.VM = app.JavaVM()
	config.Context = app.AppContext()
	config.View = evt.View
	config.RunOnMain = w.Run
}
