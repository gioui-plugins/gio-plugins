package webview

import (
	"time"
)

// DataManager is a data manager for the webview.
type DataManager interface {
	CookieManager
	StorageManager
}

// CookieManager can access and modify cookies from Webview.
//
// Cookies might be shared between Webview instances, under
// the same app or device.
type CookieManager interface {
	// Cookies returns the cookies for the current page or all browser.
	Cookies(fn DataLooper[CookieData]) (err error)
	// AddCookie adds a cookie/local-storage item.
	// The cookie must be a valid cookie string, without semi-colon, space, or comma.
	AddCookie(c CookieData) error
	// RemoveCookie removes a cookie/local-storage item.
	RemoveCookie(c CookieData) error
}

// StorageManager can access and modify LocalStorage/SessionStorage and other storage devices from Webview.
type StorageManager interface {
	// LocalStorage returns the local storage for the current page.
	LocalStorage(fn DataLooper[StorageData]) (err error)
	// AddLocalStorage adds a local storage item.
	AddLocalStorage(c StorageData) error
	// RemoveLocalStorage removes a local storage item.
	RemoveLocalStorage(c StorageData) error

	// SessionStorage returns the session storage for the current page.
	SessionStorage(fn DataLooper[StorageData]) (err error)
	// AddSessionStorage adds a session storage item.
	AddSessionStorage(c StorageData) error
	// RemoveSessionStorage removes a session storage item.
	RemoveSessionStorage(c StorageData) error
}

// DataLooper receives one or more data chunks.
// The pointer to the data is valid until the next call to DataLooper,
// you must use/copy it before the end of the call.
//
// The function must return true to continue receiving data.
type DataLooper[T CookieData | StorageData] func(d *T) (next bool)

// StorageData is a LocalStorage data.
type StorageData struct {
	Key   string
	Value string
}

// CookieData is a cookie data.
type CookieData struct {
	Name     string
	Value    string
	Domain   string
	Path     string
	Expires  time.Time
	Features CookieFeatures
}

// CookieFeatures is a set of cookie features.
type CookieFeatures uint64

const (
	// CookieSecure is a secure cookie.
	CookieSecure CookieFeatures = 1 << iota
	// CookieHTTPOnly is a http only cookie.
	CookieHTTPOnly
)

// IsSecure returns true if the cookie is secure.
func (f CookieFeatures) IsSecure() bool { return f&CookieSecure == CookieSecure }

// IsHTTPOnly returns true if the cookie is http only.
func (f CookieFeatures) IsHTTPOnly() bool { return f&CookieHTTPOnly == CookieHTTPOnly }
