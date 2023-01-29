package webview

import (
	"errors"
)

var (
	// ErrNotInstalled is returned when the webview is not installed.
	// Usually when WebView2 is not installed, see installview package for
	// more information of how install it.
	ErrNotInstalled = errors.New("webview is not installed")
	// ErrInvalidURL is returned when the URL is invalid.
	ErrInvalidURL = errors.New("invalid URL")
	// ErrInvalidSize is returned when the size is invalid.
	ErrInvalidSize = errors.New("invalid size")
	// ErrInvalidAuthProxy is returned when the proxy is invalid.
	ErrInvalidAuthProxy = errors.New("proxy with user and password is not supported")
	// ErrInvalidProxy is returned when the proxy is invalid.
	ErrInvalidProxy = errors.New("invalid proxy")
	// ErrInvalidCert is returned when the certificate is invalid.
	ErrInvalidCert = errors.New("invalid certificate")
	// ErrInvalidOptionChange is returned when an option is changed after the webview is created.
	ErrInvalidOptionChange = errors.New("invalid option change")

	// ErrNotSupported is returned when a feature is not supported.
	ErrNotSupported = errors.New("feature not supported")

	// ErrJavascriptManagerNotSupported is returned when the javascript manager is not supported.
	ErrJavascriptManagerNotSupported = errors.New("javascript manager not supported")
	// ErrJavascriptCallbackDuplicate is returned when the javascript callback is duplicated.
	ErrJavascriptCallbackDuplicate = errors.New("javascript callback is duplicated")
	// ErrJavascriptCallbackInvalidName is returned when the javascript name is too long.
	ErrJavascriptCallbackInvalidName = errors.New("javascript name is too long")
	// ErrInvalidJavascript is returned when the javascript is invalid.
	ErrInvalidJavascript = errors.New("invalid javascript")
)
