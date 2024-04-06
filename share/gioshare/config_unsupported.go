//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package gioshare

import (
	"gioui.org/app"
"github.com/gioui-plugins/gio-plugins/share"
)

// NewConfigFromViewEvent creates a share.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent, title string) share.Config {
	return share.Config{}
}

func UpdateConfigFromViewEvent(config *share.Config, w *app.Window, evt app.ViewEvent) {}

func UpdateConfigFromFrameEvent(config *share.Config, w *app.Window, evt app.FrameEvent) {}
