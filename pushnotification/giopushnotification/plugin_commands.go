package giopushnotification

import (
	"reflect"

	"gioui.org/io/event"
)

var wantCommands = []reflect.Type{
	reflect.TypeOf(RequestTokenCmd{}),
}

// RequestTokenCmd is used to request a token.
type RequestTokenCmd struct {
	Tag event.Tag
}

func (RequestTokenCmd) ImplementsCommand() {}
