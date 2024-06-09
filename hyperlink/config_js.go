//go:build js

package hyperlink

// Config is the configuration for a Hyperlink.
type Config struct {
	// Blur must be true if the current window is not in focus.
	Blur bool
}
