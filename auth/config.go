//go:build !windows && !darwin && !windows && !ios && !android

package auth

// Config is the configuration for a GoogleAuth.
//
// Each OS contains their own settings and options,
// check each config_* file for more details.
type Config struct{}
