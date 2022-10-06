package webview

import (
	"crypto/x509"
	"net/url"
	"sync"
)

var options = struct {
	sync.Mutex
	proxy struct{ ip, port string }
	certs []*x509.Certificate
}{}

// SetProxy sets the HTTP proxy to use for the webview.
// It's applied for all viewers, and must be called before creating any WebView.
//
// Only supported by Windows and Android.
func SetProxy(u *url.URL) error {
	options.Lock()
	defer options.Unlock()

	if options.proxy.ip != "" && options.proxy.port != "" {
		return ErrInvalidOptionChange
	}

	if u.User != nil {
		return ErrInvalidAuthProxy
	}

	options.proxy.ip = u.Hostname()
	options.proxy.port = u.Port()

	return nil
}

// SetCustomCertificates sets additional certificates to use for the webview.
// It's applied for all viewers, and must be called before creating any WebView.
//
// Only supported by Windows and Android.
func SetCustomCertificates(certs []*x509.Certificate) error {
	options.Lock()
	defer options.Unlock()

	if options.certs != nil {
		return ErrInvalidOptionChange
	}

	options.certs = certs
	return nil
}
