package giosafedata

import (
	"reflect"

	"github.com/gioui-plugins/gio-plugins/safedata"
)

var wantEvents = []reflect.Type{
	//reflect.TypeOf(app.ViewEvent{}),
}

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	Error error
}

// SecretsEvent is issued as response for ReadSecretOp and ListSecretOp.
type SecretsEvent struct {
	Secrets []safedata.Secret
}

func (e ErrorEvent) ImplementsEvent()   {}
func (e SecretsEvent) ImplementsEvent() {}
