package giosafedata

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/safedata"
)

// NewConfigFromViewEvent creates a safedata.Config based on app.ViewEvent.
func NewConfigFromViewEvent(w *app.Window, e app.ViewEvent, title string) safedata.Config {
	dir, _ := app.DataDir()
	return safedata.Config{App: title, VM: app.JavaVM(), Context: app.AppContext(), Folder: dir, RunOnMain: w.Run}
}
