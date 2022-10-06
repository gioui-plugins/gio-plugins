package webview

import (
	"github.com/inkeliz/go_inkwasm/inkwasm"
)

// Config is the configuration for a WebView.
type Config struct {
	// Element is the parent element of the webview.
	Element inkwasm.Object

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())

	// PxPerDp represents how many pixels per each dp.
	PxPerDp float32
}
