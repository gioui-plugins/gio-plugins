//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package safedata

// Config is the configuration for a SafeData.
//
// Each OS contains their own settings and options,
// check each config_* file for more details.
type Config struct{}
