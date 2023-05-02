package share

// Config is the configuration for a WebView.
type Config struct {
	// HWND is the handle to the window.
	HWND uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
