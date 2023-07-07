package deeplink

import "net/url"

type Event interface {
	ImplementsEvent()
}

// Linked is issued when the app is opened from a deeplink, or when the app is already running and a deeplink is opened,
// but NewDeepLink wasn't called yet.
type Linked struct {
	URL   *url.URL
	IsOld bool
}

func (Linked) ImplementsEvent() {}
