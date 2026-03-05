//go:build js && wasm

package pushnotification

// Config holds the configuration for Push on Web.
type Config struct {
	RunOnMain func(func())

	// BrowserConfig holds the configuration for Push on Web.
	BrowserConfig
}
