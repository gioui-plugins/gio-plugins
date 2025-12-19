package gioauth

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
		reflect.TypeOf(plugin.EndFrameEvent{}),
	}
	wantGlobalEvents = []reflect.Type{
		reflect.TypeOf(app.URLEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent auth.ErrorEvent

// AuthEvent is issued as response for ListenOp or RequestOp.
type AuthEvent auth.AuthenticatedEvent

func (e ErrorEvent) ImplementsEvent() {}
func (e AuthEvent) ImplementsEvent()  {}
