package webview

type cookieManager struct {
	*webview
}

func newCookieManager(w *webview) *cookieManager {
	return new(cookieManager)
}

// Cookies implements the CookieManager interface.
func (*cookieManager) Cookies(_ DataLooper[CookieData]) (err error) {
	return ErrNotSupported
}

// AddCookie implements the CookieManager interface.
func (*cookieManager) AddCookie(_ CookieData) error {
	return ErrNotSupported
}

// RemoveCookie implements the CookieManager interface.
func (*cookieManager) RemoveCookie(_ CookieData) error {
	return ErrNotSupported
}
