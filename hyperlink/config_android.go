package hyperlink

// Config is the configuration for a Hyperlink.
type Config struct {
	// View is the Android View.
	View uintptr

	// VM is the Java VM.
	VM uintptr

	// Context is the Android Context.
	Context uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
