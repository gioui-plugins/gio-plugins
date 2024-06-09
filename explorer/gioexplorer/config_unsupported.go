//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package gioexplorer

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/explorer"
)

// NewConfigFromViewEvent creates a explorer.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent, title string) explorer.Config {
	return explorer.Config{}
}

// UpdateConfigFromViewEvent updates explorer.Config based on app.ViewEvent.
func UpdateConfigFromViewEvent(config *explorer.Config, w *app.Window, e app.ViewEvent) {}
