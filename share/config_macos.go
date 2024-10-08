//go:build darwin && !ios

package share

// Config is the configuration for a Share.
type Config struct {
	// View is a CFTypeRef for the NSView for the window.
	View uintptr

	// Layer is a CFTypeRef of the CALayer of View.
	Layer uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())

	// Size is the size of window (in pixels).
	Size [2]float32

	// PxPerDp represents how many pixels per each dp.
	PxPerDp float32
}
