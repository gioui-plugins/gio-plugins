package share

// Config is the configuration for a Share.
type Config struct {
	// View is a CFTypeRef for the UIViewController for the window.
	View uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())

	// PxPerDp represents how many pixels per each dp.
	PxPerDp float32
}
