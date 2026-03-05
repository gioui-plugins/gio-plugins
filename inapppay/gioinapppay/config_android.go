//go:build android

package gioinapppay

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/inapppay"
)

// NewConfigFromViewEvent creates a new Config from a ViewEvent.
func NewConfigFromViewEvent(w *app.Window, event app.ViewEvent) inapppay.Config {
	c := inapppay.Config{}
	UpdateConfigFromViewEvent(&c, w, event)
	return c
}

// UpdateConfigFromViewEvent updates the Config from a ViewEvent.
func UpdateConfigFromViewEvent(c *inapppay.Config, w *app.Window, event app.ViewEvent) {
	evt, ok := event.(app.AndroidViewEvent)
	if !ok {
		return
	}
	c.View = evt.View
	c.VM = app.JavaVM()
	c.Context = app.AppContext()
	c.RunOnMain = w.Run
}
