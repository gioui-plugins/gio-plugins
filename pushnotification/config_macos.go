//go:build darwin && !ios

package pushnotification

// Config holds the configuration for Push on macOS.
type Config struct {
	View uintptr
	RunOnMain func(func())
}