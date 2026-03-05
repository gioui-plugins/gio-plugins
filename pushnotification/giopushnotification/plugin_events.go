package giopushnotification

import (
	"reflect"

	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/pushnotification"
)

var (
	wantCommand = []reflect.Type{
		reflect.TypeOf(RequestTokenCmd{}),
	}
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
	}
)

// TokenReceivedEvent wrapper.
type TokenReceivedEvent pushnotification.Token

// ErrorEvent wrapper.
type ErrorEvent struct {
	Error error
}

func (e TokenReceivedEvent) ImplementsEvent() {}
func (e ErrorEvent) ImplementsEvent()         {}
