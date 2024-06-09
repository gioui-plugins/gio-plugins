package giosafedata

import (
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"

	"github.com/gioui-plugins/gio-plugins/safedata"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	Error error
}

// SecretsEvent is issued as response for ReadSecretCmd and ListSecretCmd.
type SecretsEvent struct {
	Secrets []safedata.Secret
}

func (e ErrorEvent) ImplementsEvent()   {}
func (e SecretsEvent) ImplementsEvent() {}
