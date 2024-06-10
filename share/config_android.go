package share

// Config is the configuration for a Share.
type Config struct {
	// VM is the Java VM.
	VM uintptr

	View uintptr

	// Context is the Android Context.
	Context uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
