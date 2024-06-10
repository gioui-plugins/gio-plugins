package auth

import (
	"errors"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"net/url"
	"sync"
)

var (
	// ErrProviderNotFound  is returned when the provider is not found,
	// you need to create Auth with the given providers.Provider, before
	// calling Open.
	ErrProviderNotFound = errors.New("provider not found")

	// ErrNotConfigured is returned when the Auth is not configured,
	// check if Config contains a valid View and equivalent settings.
	ErrNotConfigured = errors.New("auth not configured")

	// ErrProviderNotAllowed is returned when the provider is not allowed.
	ErrProviderNotAllowed = errors.New("provider not allowed, check your configuration and/or the provider's documentation. If you are using Google, check your SHA-1 fingerprint and package name. If you are using Apple, check your bundle ID")

	// ErrUserCancelled is returned when the user cancels the authentication.
	ErrUserCancelled = errors.New("user cancelled")
)

// Auth is the main struct, which holds the driver.
type Auth struct {
	// driver holds OS-Specific content, it varies for each OS.
	driver

	providersMutex sync.Mutex
	providers      map[providers.Identifier]providers.Provider

	eventsMutex sync.Mutex
	eventsChan  []chan Event
}

type idriver interface {
	// Open authenticates the user.
	open(provider providers.Provider, nonce string) error
}

// NewAuth returns a new Auth struct.
func NewAuth(config Config, provider ...providers.Provider) *Auth {
	house := &Auth{
		providers: make(map[providers.Identifier]providers.Provider),
	}
	for _, p := range provider {
		house.AddProvider(p)
	}
	attachDriver(house, config)
	return house
}

// AddProvider adds a new provider to the Auth
func (a *Auth) AddProvider(provider providers.Provider) {
	a.providersMutex.Lock()
	defer a.providersMutex.Unlock()

	if provider == nil {
		return
	}
	if a.providers == nil {
		a.providers = make(map[providers.Identifier]providers.Provider)
	}
	a.providers[provider.Identifier()] = provider
}

// Configure configures the GoogleAuth.
func (a *Auth) Configure(config Config) {
	configureDriver(&a.driver, config)
}

// Open displays the authentication window.
// It will discard any previous authentication, if any.
//
// Nonce is a random string, which is used to prevent replay attacks,
// however it can be an image/result of a hash function, which you can
// later use to verify the authenticity of the response, using the
// pre-image of the hash, sent by the user.
func (a *Auth) Open(provider providers.Identifier, nonce string) error {
	a.providersMutex.Lock()
	p, ok := a.providers[provider]
	a.providersMutex.Unlock()
	if !ok {
		return ErrProviderNotFound
	}
	return a.driver.open(p, nonce)
}

// ProcessCustomSchemeCallback processes the URL received as Custom Scheme.
//
// In general, that is used internally, but you need to call it manually if you receive the URL from the OS,
// as custom scheme URLs are not automatically processed.
func (a *Auth) ProcessCustomSchemeCallback(uri string) error {
	u, err := url.Parse(uri)
	if err != nil {
		return err
	}

	if u.RawQuery == "" && u.Fragment != "" {
		u.RawQuery = u.Fragment
	}

	values := u.Query()
	currentState := values.Get("state")
	for _, p := range a.providers {
		if currentState != p.Identifier().State() {
			continue
		}

		a.sendResponse(AuthenticatedEvent{
			Provider: p.Identifier(),
			Code:     values.Get("code"),
			IDToken:  values.Get("id_token"),
		})
	}

	return nil
}

// Events returns the response from the user.
func (a *Auth) Events() <-chan Event {
	a.eventsMutex.Lock()
	defer a.eventsMutex.Unlock()

	c := make(chan Event, 8)
	a.eventsChan = append(a.eventsChan, c)
	return c
}

func (a *Auth) sendResponse(event Event) {
	a.eventsMutex.Lock()
	defer a.eventsMutex.Unlock()

	for _, c := range a.eventsChan {
		select {
		case c <- event:
		}
	}
}
