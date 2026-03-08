//go:build windows

package pushnotification

// Config holds the configuration for Push on Windows.
type Config struct {
	// HWND is the HWND of the window.
	HWND uintptr
	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(func())
}
