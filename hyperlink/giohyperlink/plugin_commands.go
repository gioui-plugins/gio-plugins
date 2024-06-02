package giohyperlink

import (
	"gioui.org/io/event"
	"net/url"
	"reflect"
)

var (
	wantCommands = []reflect.Type{
		reflect.TypeOf(OpenCmd{}),
	}
)

// OpenCmd is an operation that will open a URL.
// It will issue an ErrorEvent if an error occurs.
type OpenCmd struct {
	Tag event.Tag
	URI *url.URL
}

func (o OpenCmd) ImplementsCommand() {}
