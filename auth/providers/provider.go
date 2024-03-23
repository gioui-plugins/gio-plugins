package providers

import "runtime"

type Identifier string

func (i Identifier) State() string {
	return string(i) + "-" + runtime.GOOS
}

type Provider interface {
	URL(nonce string) string
	ClientID() string
	Scheme() string
	Identifier() Identifier
}
