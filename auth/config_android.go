//go:build android

package auth

// Config is the configuration for a GoogleAuth.
type Config struct {
	// VM is the Java VM.
	VM uintptr

	// View is the Android View.
	View uintptr

	// Context is the Android Context.
	Context uintptr

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
