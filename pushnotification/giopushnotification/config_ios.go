//go:build ios

package giopushnotification

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/pushnotification"
)

// NewConfigFromViewEvent creates a new Config from a ViewEvent.
func NewConfigFromViewEvent(w *app.Window, event app.ViewEvent, extra []pushnotification.ExternalConfig) pushnotification.Config {
	c := pushnotification.Config{}
	UpdateConfigFromViewEvent(&c, w, event, extra)
	return c
}

// UpdateConfigFromViewEvent updates the Config from a ViewEvent.
func UpdateConfigFromViewEvent(c *pushnotification.Config, w *app.Window, event app.ViewEvent, extra []pushnotification.ExternalConfig) {
	evt, ok := event.(app.UIKitViewEvent)
	if !ok {
		return
	}
	c.View = evt.ViewController
	c.RunOnMain = w.Run
}
