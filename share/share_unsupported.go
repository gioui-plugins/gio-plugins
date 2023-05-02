// SPDX-License-Identifier: Unlicense OR MIT

//go:build !windows && !js && !ios && !android && !darwin
// +build !windows,!js,!ios,!android,!darwin

package share

type driver struct{}

func attachDriver(house *Share, config Config) {}

func configureDriver(driver *driver, config Config) {}

func (e *driver) shareText(title, text string) error {
	return ErrNotAvailable
}

func (e *driver) shareWebsite(title, description, url string) error {
	return ErrNotAvailable
}
