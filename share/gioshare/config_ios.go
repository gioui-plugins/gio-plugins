package gioshare

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/share"
)

// NewConfigFromViewEvent creates a share.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent) share.Config {
	r := share.Config{}
	UpdateConfigFromViewEvent(&r, w, evt)
	return r
}

func UpdateConfigFromViewEvent(config *share.Config, w *app.Window, e app.ViewEvent) {
	evt, ok := e.(app.UIKitViewEvent)
	if !ok {
		return
	}

	config.View = evt.ViewController
	config.RunOnMain = w.Run
}

func UpdateConfigFromFrameEvent(config *share.Config, w *app.Window, evt app.FrameEvent) {
	config.PxPerDp = evt.Metric.PxPerDp
}
