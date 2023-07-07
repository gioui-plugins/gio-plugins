package deeplink

import (
	"fmt"
	"github.com/gioui-plugins/gio-plugins/deeplink/internal"
	"net/url"
	"sync"
)

var _Deeplink = new(DeepLink)

type DeepLink struct {
	fan   internal.Fan[Event]
	mutex sync.Mutex

	lastURL *url.URL

	lastError error
}

// NewDeepLink creates a new DeepLink.
// The driver is different for each platform.
func NewDeepLink() *DeepLink {
	return _Deeplink
}

// Events returns a channel that receives events.
func (e *DeepLink) Events() <-chan Event {
	ch := e.fan.Add(1)
	e.mutex.Lock()
	defer e.mutex.Unlock()
	if e.lastURL != nil {
		ch <- Linked{URL: e.lastURL, IsOld: true}
	}
	return ch
}

// LastURL returns the last URL.
func (e *DeepLink) LastURL() *url.URL {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return e.lastURL
}

// Schemes returns the list of supported schemes.
func (e *DeepLink) Schemes() []string {
	return schemeList
}

// updateURL updates the last URL that was opened.
func (e *DeepLink) updateURL(u string) {
	fmt.Println("updateURL", u)
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.lastURL, _ = url.Parse(u)
	e.fan.Send(Linked{URL: e.lastURL, IsOld: false})
}

// updateError updates the last error that was opened.
func (e *DeepLink) updateError(err error) {
	if err == nil {
		return
	}
	e.mutex.Lock()
	defer e.mutex.Unlock()
	e.lastError = err
}
