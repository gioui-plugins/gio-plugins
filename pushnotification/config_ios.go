//go:build ios

package pushnotification

// Config holds the configuration for Push on iOS.
type Config struct {
	View 	  uintptr
	RunOnMain func(func())
}