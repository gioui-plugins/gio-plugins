//go:build android

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
	evt, ok := event.(app.AndroidViewEvent)
	if !ok {
		return
	}
	c.View = evt.View
	c.VM = app.JavaVM()
	c.Context = app.AppContext()
	c.RunOnMain = w.Run

	for _, ext := range extra {
		switch eee := ext.(type) {
		case pushnotification.AndroidFirebaseConfig:
			c.AndroidFirebaseConfig = eee
		}
	}
}
