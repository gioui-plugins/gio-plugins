package auth

import "github.com/gioui-plugins/gio-plugins/auth/providers"

type Event interface {
	ImplementsEvent()
}

// AuthenticatedEvent is the response from the user.
type AuthenticatedEvent struct {
	// Provider is the providers.Identifier from Provider used to authenticate.
	Provider providers.Identifier

	// Code is the code from the user.
	// This is only available if the user has authenticated.
	Code string

	// IDToken is the ID Token (JWT/OpenConnect) from the user.
	// This is only available if the user has authenticated.
	IDToken string
}

// ErrorEvent is the error from the user.
type ErrorEvent struct {
	Error error
}

func (a AuthenticatedEvent) ImplementsEvent() {}
func (a ErrorEvent) ImplementsEvent()         {}
