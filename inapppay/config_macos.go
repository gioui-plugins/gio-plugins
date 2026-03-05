//go:build darwin && !ios

package inapppay

// Config holds the configuration for the InAppPay.
type Config struct {
	// View is the view of the window.
	View uintptr
	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(func())
}
