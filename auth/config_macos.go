//go:build darwin && !ios

package auth

// Config is the configuration for a GoogleAuth.
type Config struct {
	// View is a CFTypeRef for the NSView for the window.
	View uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
