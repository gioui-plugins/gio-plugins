package gioshare

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"
)

var wantEvents = []reflect.Type{
	reflect.TypeOf(plugin.ViewEvent{}),
	reflect.TypeOf(app.DestroyEvent{}),
	reflect.TypeOf(app.FrameEvent{}),
}

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	Error error
}

func (ErrorEvent) ImplementsEvent() {}
