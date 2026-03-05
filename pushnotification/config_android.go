//go:build android

package pushnotification

// Config holds the configuration for Push on Android.
type Config struct {
	View      uintptr
	VM        uintptr
	Context   uintptr
	RunOnMain func(f func())

	// AndroidFirebaseConfig holds the configuration for Push on Android.
	AndroidFirebaseConfig
}
