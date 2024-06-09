package providers

import "runtime"

// Identifier is the identifier for the provider.
// Must be unique for each provider.
type Identifier string

// State returns the "state" for the provider.
func (i Identifier) State() string {
	return string(i) + "-" + runtime.GOOS
}

// Provider is the interface for the OpenID Connect provider.
type Provider interface {
	// URL returns the URL to open the provider.
	URL(nonce string) string

	// ClientID returns the client ID for the provider.
	ClientID() string

	// Scheme returns the scheme for the provider.
	Scheme() string

	// Identifier returns the identifier for the provider.
	Identifier() Identifier
}
