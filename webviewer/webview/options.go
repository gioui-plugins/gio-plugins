package webview

import (
	"crypto/x509"
	"net/url"
	"sync"
)

var options = struct {
	sync.Mutex
	proxy  struct{ ip, port string }
	certs  []*x509.Certificate
	folder string
	debug  int8
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

// SetDirectory sets the folder to use for the webview.
// It's applied for all viewers, and must be called before creating any WebView.
//
// Only supported by Windows.
func SetDirectory(folder string) error {
	options.Lock()
	defer options.Unlock()

	if options.folder != "" {
		return ErrInvalidOptionChange
	}

	options.folder = folder
	return nil
}

// SetDebug enables the debug (such as the inspector) for the webview.
// It's applied for all viewers, and must be called before creating any WebView.
//
// Only supported by macOS, iOS and Android.
func SetDebug(enable bool) error {
	options.Lock()
	defer options.Unlock()

	if options.debug != 0 {
		return ErrInvalidOptionChange
	}

	if enable {
		options.debug = 1
	} else {
		options.debug = -1
	}
	return nil
}
