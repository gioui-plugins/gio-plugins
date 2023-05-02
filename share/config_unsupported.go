//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package share

// Config is the configuration for a WebView.
type Config struct{}
