package giohyperlink

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
		reflect.TypeOf(app.ConfigEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	Error error
}

// ImplementsEvent implements event.Event.
func (e ErrorEvent) ImplementsEvent() {}
