package gioauth

import (
	"gioui.org/io/system"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"

	"gioui.org/app"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(app.ViewEvent{}),
		reflect.TypeOf(system.FrameEvent{}),
		reflect.TypeOf(plugin.EndFrameEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent auth.ErrorEvent

// AuthEvent is issued as response for ListenOp or RequestOp.
type AuthEvent auth.AuthenticatedEvent

func (e ErrorEvent) ImplementsEvent() {}
func (e AuthEvent) ImplementsEvent()  {}
