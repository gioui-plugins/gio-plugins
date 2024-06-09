package gioauth

import (
	"gioui.org/io/event"
)

// Filter is used to receive an authentication, and it will be responded with an AuthEvent.
type Filter struct{}

// ImplementsFilter implements event.Filter.
func (f Filter) ImplementsFilter() {}

// Name implements plugin.UntaggedFilter
func (f Filter) Name() uint64 { return intName }

// Matches implements plugin.UntaggedFilter
func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case AuthEvent:
		return true
	case ErrorEvent:
		return true
	default:
		return false
	}
}
