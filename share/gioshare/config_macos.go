//go:build darwin && !ios

package gioshare

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"github.com/gioui-plugins/gio-plugins/share"
)

// NewConfigFromViewEvent creates a share.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) share.Config {
	r := share.Config{}
	UpdateConfigFromViewEvent(&r, w, evt)
	return r
}

func UpdateConfigFromViewEvent(config *share.Config, w *app.Window, evt app.ViewEvent) {
	config.View = evt.View
	config.Layer = evt.Layer
	config.RunOnMain = w.Run
}

func UpdateConfigFromFrameEvent(config *share.Config, w *app.Window, evt system.FrameEvent) {
	config.PxPerDp = evt.Metric.PxPerDp
	config.Size = [2]float32{float32(evt.Size.X), float32(evt.Size.Y)}
}
