//go:build ios

package auth

// Config is the configuration for Auth.
type Config struct {
	// View is a CFTypeRef for the UIViewController for the window.
	View uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
