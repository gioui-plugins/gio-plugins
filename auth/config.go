//go:build !windows && !darwin && !windows && !ios && !android && !js

package auth

// Config is the configuration for Auth.
//
// Each OS contains their own settings and options,
// check each config_* file for more details.
type Config struct{}
