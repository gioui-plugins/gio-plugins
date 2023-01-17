//go:build darwin && !ios

package giosafedata

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/safedata"
)

// NewConfigFromViewEvent creates a safedata.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, evt app.ViewEvent, title string) safedata.Config {
	return safedata.Config{App: title}
}
