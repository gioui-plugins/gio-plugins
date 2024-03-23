package google

import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"net/url"
	"runtime"
)

const IdentifierGoogle = providers.Identifier("google")

// Provider is the settings for the provider, Google.
type Provider struct {
	// WebClientID is the OAuth 2.0 Client ID, this is required.
	// You MUST obtain this from the Google API Console, create one "Web App".
	WebClientID string

	// DesktopClientID is the OAuth 2.0 Client ID, this is required.
	// You MUST obtain this from the Google API Console, create one "Desktop App".
	DesktopClientID string

	// RedirectURL is the URL to redirect to after the user has authenticated.
	// You MUST register this URL in the Google API Console.
	// It MUST be the current URL, on JS.
	RedirectURL string
}

// URL returns the URL for the current platform.
func (c *Provider) URL(nonce string) string {
	switch runtime.GOOS {
	case "android", "js":
		return c.webURL(nonce)
	default:
		return c.desktopURL(nonce)
	}
}

// ClientID returns the client ID for the current platform.
func (c *Provider) ClientID() string {
	switch runtime.GOOS {
	case "android", "js":
		return c.WebClientID
	default:
		return c.DesktopClientID
	}
}

// webURL is used in Android and JS.
func (c *Provider) webURL(nonce string) string {
	v := url.Values{}
	v.Set("scope", "openid email profile")
	v.Set("response_type", "code id_token")
	v.Set("state", c.Identifier().State())
	v.Set("url", "https://oauth2.example.com/token")
	v.Set("redirect_uri", c.RedirectURL)
	v.Set("client_id", c.WebClientID)
	v.Set("nonce", nonce)
	v.Set("response_mode", "query")

	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

// desktopURL is used in Windows, macOS, iOS.
func (c *Provider) desktopURL(nonce string) string {
	v := url.Values{}
	v.Set("scope", "openid email profile")
	v.Set("response_type", "code id_token")
	v.Set("state", c.Identifier().State())
	v.Set("url", "https://oauth2.example.com/token")
	v.Set("redirect_uri", c.Scheme()+":/googleauth")
	v.Set("client_id", c.DesktopClientID)
	v.Set("nonce", nonce)
	v.Set("response_mode", "query")

	return "https://accounts.google.com/o/oauth2/v2/auth?" + v.Encode()
}

// Scheme returns the scheme for the DesktopClientID, used in Windows, macOS, iOS.
func (c *Provider) Scheme() (out string) {
	last := len(c.DesktopClientID)
	for i := len(c.DesktopClientID) - 1; i > 0; i-- {
		if c.DesktopClientID[i] == '.' {
			out += c.DesktopClientID[i+1:last] + "."
			last = i
		}
	}
	return out + c.DesktopClientID[0:last]
}

func (c *Provider) Identifier() providers.Identifier {
	return IdentifierGoogle
}
