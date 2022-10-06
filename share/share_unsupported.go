// SPDX-License-Identifier: Unlicense OR MIT

//go:build !windows && !js && !ios && !android && !darwin
// +build !windows,!js,!ios,!android,!darwin

package share

import (
	"gioui.org/app"
	"gioui.org/io/event"
)

type share struct{}

func newShare(w *app.Window) *share {
	return new(share)
}

func (e *sharePlugin) listenEvents(_ event.Event) {

}

func (e *sharePlugin) shareShareable(shareable Shareable) error {
	return ErrNotAvailable
}
