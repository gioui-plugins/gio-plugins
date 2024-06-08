package gioauth

import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"reflect"
)

var (
	wantCommands = []reflect.Type{
		reflect.TypeOf(OpenCmd{}),
	}
)

// OpenCmd is used to request an authentication, and it will be responded with an AuthEvent to ListenOp,
// the Tag is used to receive errors, and the Provider is used to identify the provider.
type OpenCmd struct {
	// Provider is the provider to use.
	Provider providers.Identifier

	// Nonce is a random string that will be returned in the AuthEvent.
	Nonce string
}

func (OpenCmd) ImplementsCommand() {}
