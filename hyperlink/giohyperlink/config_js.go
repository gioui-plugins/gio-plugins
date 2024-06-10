//go:build !android

package giohyperlink

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/hyperlink"
)

// NewConfigFromViewEvent creates a hyperlink.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) hyperlink.Config {
	r := hyperlink.Config{}
	UpdateConfigFromViewEvent(&r, w, e)
	return r
}

// UpdateConfigFromViewEvent updates a hyperlink.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *hyperlink.Config, w *app.Window, e app.ViewEvent) {}

// UpdateConfigFromConfigEvent creates a hyperlink.Config based on app.ConfigEvent.
func UpdateConfigFromConfigEvent(config *hyperlink.Config, _ *app.Window, e app.ConfigEvent) {
	config.Blur = e.Config.Focused != true
}
