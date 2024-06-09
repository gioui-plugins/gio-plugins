package giosafedata

import "gioui.org/io/event"

// Filter is used to filter events.
type Filter struct {
	Source event.Tag
}

// ImplementsFilter implements plugin.Filter and event.Filter.
func (f Filter) ImplementsFilter() {}

// Tag implements plugin.Filter.
func (f Filter) Tag() event.Tag { return f.Source }

// Matches implements plugin.Filter.
func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case ErrorEvent:
		return true
	case SecretsEvent:
		return true
	default:
		return false
	}
}
