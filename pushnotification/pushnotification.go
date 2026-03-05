package pushnotification

import (
	"errors"

	"golang.org/x/sync/singleflight"
)

var (
	// ErrNotAvailable is returned when the push notification service is not available.
	ErrNotAvailable = errors.New("push notification not available")
	// ErrNotConfigured is returned when the plugin is not configured.
	ErrNotConfigured = errors.New("push notification not configured")
)

type Platform string

const (
	PlatformWindows Platform = "windows"
	PlatformAndroid Platform = "android"
	PlatformIOS     Platform = "ios"
	PlatformMacOS   Platform = "macos"
	PlatformWeb     Platform = "web"
)

// Token represents a push notification token, used to send push notifications.
type Token struct {
	Token    string   `json:"token"`
	Platform Platform `json:"platform"`
}

// WebPushToken represents a Web Push token, that is the value of Token if Platform is PlatformWeb.
// Make sure to validate the Endpoint in your server to protect against SSRF.
type WebPushToken struct {
	Endpoint string      `json:"endpoint"`
	Keys     WebPushKeys `json:"keys"`
}

// WebPushKeys represents the keys for a WebPushToken.
type WebPushKeys struct {
	Auth   string `json:"auth"`
	P256dh string `json:"p256dh"`
}

// Push handles push notification interactions.
type Push struct {
	*driver
	mutex singleflight.Group
}

// NewPush creates a new Push instance.
func NewPush(config Config) *Push {
	p := &Push{}
	attachDriver(p, config)
	return p
}

// Configure updates the configuration.
func (p *Push) Configure(config Config) {
	configureDriver(p.driver, config)
}

// RequestToken requests the push token from the underlying system.
func (p *Push) RequestToken() (Token, error) {
	v, err, _ := p.mutex.Do("requestToken", func() (interface{}, error) {
		if p.driver == nil {
			return Token{}, ErrNotConfigured
		}
		token, err := p.driver.requestToken()
		return token, err
	})
	return v.(Token), err
}

type idriver interface {
	requestToken() (Token, error)
}
