// SPDX-License-Identifier: Unlicense OR MIT

package share

import (
	"errors"
)

var (
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")

	// ErrNotAvailableAction is return when the current Shareable item isn't supported.
	ErrNotAvailableAction = errors.New("current shareable item not supported")
)

type Share struct {
	// share holds OS-Specific content, it varies for each OS.
	driver
}

func NewShare(config Config) *Share {
	s := new(Share)
	attachDriver(s, config)
	return s
}

func (s *Share) Configure(config Config) {
	configureDriver(&s.driver, config)
}

func (s *Share) Text(title, desc string) error {
	return s.shareText(title, desc)
}

func (s *Share) Website(title, desc, link string) error {
	return s.shareWebsite(title, desc, link)
}
