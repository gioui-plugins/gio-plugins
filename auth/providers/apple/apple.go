package apple

import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"net/url"
)

const IdentifierApple = providers.Identifier("apple")

// Provider is the settings for the provider, Apple.
type Provider struct {
	// ServiceIdentifier  is the OAuth 2.0 Client ID, this is required.
	// You MUST obtain this from the "Apple Certificates, Identifiers & Profiles", in the "Identifiers" section,
	// then "Services IDs".
	// The ServiceIdentifier is the "Identifier" field.
	ServiceIdentifier string

	// RedirectURL is the URL to redirect to after the user has authenticated.
	// You MUST register this URL in your Service ID.
	// It MUST be the current URL, on JS.
	RedirectURL string

	// Scheme is the scheme for which will be used, notice that Apple doesn't
	// allow to use custom schemes as RedirectURL. So, you MUST redirect to
	// a custom scheme from the RedirectURL.
	SchemeURL string
}

// URL returns the URL for the current platform.
func (c *Provider) URL(nonce string) string {
	return c.webURL(nonce)
}

// ClientID returns the client ID for the current platform.
func (c *Provider) ClientID() string {
	return c.ServiceIdentifier
}

// webURL is used in Android and JS.
func (c *Provider) webURL(nonce string) string {
	v := url.Values{}
	v.Set("scope", "email name")
	v.Set("response_type", "code id_token")
	v.Set("state", c.Identifier().State())
	v.Set("redirect_uri", c.RedirectURL)
	v.Set("client_id", c.ClientID())
	v.Set("nonce", nonce)
	v.Set("response_mode", "form_post")

	return "https://appleid.apple.com/auth/authorize?" + v.Encode()
}

// Scheme returns the scheme for Windows, Android and macOS (non-AppStore).
func (c *Provider) Scheme() (out string) {
	return c.SchemeURL
}

func (c *Provider) Identifier() providers.Identifier {
	return IdentifierApple
}
