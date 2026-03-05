//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package giopushnotification

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/pushnotification"
)

// NewConfigFromViewEvent creates a new Config from a ViewEvent.
func NewConfigFromViewEvent(w *app.Window, event app.ViewEvent, extra []pushnotification.ExternalConfig) pushnotification.Config {
	return pushnotification.Config{}
}

// UpdateConfigFromViewEvent updates the Config from a ViewEvent.
func UpdateConfigFromViewEvent(c *pushnotification.Config, w *app.Window, event app.ViewEvent, extra []pushnotification.ExternalConfig) {
}
