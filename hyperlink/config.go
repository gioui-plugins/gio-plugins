//go:build !android && !js

package hyperlink

// Config is the configuration for a WebView.
//
// Each OS contains their own settings and options,
// check each config_* file for more details.
type Config struct{}
