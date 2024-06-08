package gioauth

import (
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"

	"gioui.org/app"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
		reflect.TypeOf(app.FrameEvent{}),
		reflect.TypeOf(plugin.EndFrameEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent auth.ErrorEvent

// AuthEvent is issued as response for ListenOp or RequestOp.
type AuthEvent auth.AuthenticatedEvent

func (e ErrorEvent) ImplementsEvent() {}
func (e AuthEvent) ImplementsEvent()  {}
