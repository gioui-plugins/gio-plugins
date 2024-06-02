//go:build darwin && !ios

package gioshare

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/share"
)

// NewConfigFromViewEvent creates a share.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent) share.Config {
	r := share.Config{}
	UpdateConfigFromViewEvent(&r, w, e)
	return r
}

func UpdateConfigFromViewEvent(config *share.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.AppKitViewEvent)
	if !ok {
		return
	}

	config.View = evt.View
	config.Layer = evt.Layer
	config.RunOnMain = w.Run
}

func UpdateConfigFromFrameEvent(config *share.Config, w *app.Window, e app.FrameEvent) {
	config.PxPerDp = e.Metric.PxPerDp
	config.Size = [2]float32{float32(e.Size.X), float32(e.Size.Y)}
}
