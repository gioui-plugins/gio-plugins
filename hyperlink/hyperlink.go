package hyperlink

import (
	"errors"
	"net/url"
)

var (
	// ErrNotReady may occur when try to open a URL before the initialization is done.
	ErrNotReady = errors.New("some needed library was not loaded yet, make use that you are using ListenEvents()")
	// ErrInvalidURL occur when provide an invalid URL, like a non http/https URL.
	ErrInvalidURL = errors.New("given url is invalid")
)

var (
	// InsecureIgnoreScheme will remove any attempt to validate the URL
	// It's "false" by default. Set it to "true" if you are using a custom scheme (like "myapp://").
	InsecureIgnoreScheme bool
)

// Hyperlink is the main struct, which holds the driver.
type Hyperlink struct {
	// driver holds OS-Specific content, it varies for each OS.
	driver
}

// NewHyperlink creates a new hyperlink.
func NewHyperlink(config Config) *Hyperlink {
	r := new(Hyperlink)
	attachDriver(r, config)
	return r
}

// Configure reconfigures the driver.
func (h *Hyperlink) Configure(config Config) {
	configureDriver(&h.driver, config)
}

// Open opens the given URL.
// It will return ErrInvalidURL if the URL doesn't use http or https.
//
// If you want to ignore the scheme, set InsecureIgnoreScheme to true, or use OpenUnsafe.
func (h *Hyperlink) Open(uri *url.URL) error {
	return h.OpenWith(uri, "")
}

// OpenWith opens the given URL with some preferred package.
// It will return ErrInvalidURL if the URL doesn't use http or https.
//
// The preferredPackage is the name of the package (on Android), such as "com.android.chrome".
//
// If you want to ignore the scheme, set InsecureIgnoreScheme to true, or use OpenUnsafe.
func (h *Hyperlink) OpenWith(uri *url.URL, preferredPackage string) error {
	if uri == nil || uri.Scheme == "" || ((uri.Scheme != "http" && uri.Scheme != "https") && InsecureIgnoreScheme == false) {
		return ErrInvalidURL
	}
	return h.driver.open(uri, preferredPackage)
}

// OpenUnsafe is the same as Open, but it doesn't validate the URL.
//
// That may crash or cause unexpected behavior if you use a non http/https URL,
// so use it with caution.
func (h *Hyperlink) OpenUnsafe(uri *url.URL) error {
	return h.driver.open(uri, "")
}
