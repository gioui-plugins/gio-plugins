package safehouse

// Config is the configuration for a WebView.
type Config struct {
	// App is the name of the app, which is used, in some
	// OSes to identify the app who creates the
	// credentials.
	App string

	// HWND is the handle to the window.
	HWND uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
