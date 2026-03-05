//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package pushnotification

// Config holds the configuration for Push.
type Config struct{}
