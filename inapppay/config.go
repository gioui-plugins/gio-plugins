//go:build !android && !ios && !darwin
// +build !android,!ios,!darwin

package inapppay

import "gioui.org/app"

type Config struct{}

func NewConfigFromViewEvent(w *app.Window, event app.ViewEvent) Config {
	return Config{}
}

func UpdateConfigFromViewEvent(c *Config, w *app.Window, event app.ViewEvent) {}
