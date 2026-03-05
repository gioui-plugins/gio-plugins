package inapppay

import (
	"errors"
	"sync"
)

var (
	// ErrNotAvailable is returned when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")

	// ErrProviderNotFound is returned when the store provider is not found.
	ErrProviderNotFound = errors.New("store provider not found")

	// ErrNotConfigured is returned when InAppPay is not configured.
	ErrNotConfigured = errors.New("inapppay not configured")

	// ErrUserCancelled is returned when the user cancels the purchase.
	ErrUserCancelled = errors.New("user cancelled")
)

// InAppPay is the main struct for In-App Purchases.
type InAppPay struct {
	driver *driver

	eventsMutex sync.Mutex
	eventsChan  []chan Event
}

// NewInAppPay creates a new InAppPay instance.
// Provide config if needed on creation (often just ViewEvent later).
func NewInAppPay(config Config) *InAppPay {
	iap := &InAppPay{}
	attachDriver(iap, config)
	return iap
}

// Configure updates the configuration (e.g. new View).
func (p *InAppPay) Configure(config Config) {
	configureDriver(p.driver, config)
}

// ListProducts requests the list of products from the stores.
// Results will be delivered via Events().
//
// That only returns one-time purchasable products.
// Use ListMemberships for subscriptions.
func (p *InAppPay) ListProducts(productIDs []string) error {
	return p.driver.listProducts(productIDs)
}

// Purchase initiates a payment for a specific product.
// An arbitrary identifier (developerPayload) can be provided.
//
// That only works if the provided productID is one-time purchasable.
func (p *InAppPay) Purchase(productID string, customPayload string, isPersonalizedPrice bool) error {
	return p.driver.purchase(productID, customPayload, isPersonalizedPrice)
}

// ListMemberships requests the list of products from the stores.
// Results will be delivered via Events().
//
// That only returns subscriptions.
// Use ListProducts for one-time purchasable products.
func (p *InAppPay) ListMemberships(productIDs []string) error {
	return errors.New("not implemented: ListMemberships for subscriptions")
}

// Subscribe initiates a subscription for a specific product.
// An arbitrary identifier (developerPayload) can be provided.
//
// That only works if the provided productID is a subscription.
// Use Purchase for one-time purchasable products.
func (p *InAppPay) Subscribe(productID string, developerPayload string) error {
	return errors.New("not implemented: Subscribe for subscriptions")
}

// Events returns a channel to receive events (ProductDetails, PaymentResult).
func (p *InAppPay) Events() <-chan Event {
	p.eventsMutex.Lock()
	defer p.eventsMutex.Unlock()

	c := make(chan Event, 8)
	p.eventsChan = append(p.eventsChan, c)
	return c
}

func (p *InAppPay) sendResponse(event Event) {
	p.eventsMutex.Lock()
	defer p.eventsMutex.Unlock()

	for _, c := range p.eventsChan {
		select {
		case c <- event:
		default:
		}
	}
}
