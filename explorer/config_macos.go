//go:build darwin && !ios

package explorer

// Config is the configuration for a Explorer.
type Config struct {
	// View is a CFTypeRef for the NSView for the window.
	View uintptr

	// Layer is a CFTypeRef of the CALayer of View.
	Layer uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
