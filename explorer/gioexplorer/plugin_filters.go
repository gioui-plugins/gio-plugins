package gioexplorer

import "gioui.org/io/event"

type Filter struct {
	Target event.Tag
}

// ImplementsFilter implements the event.Filter interface.
func (f Filter) ImplementsFilter() {}

// Tag implements the plugin.Filter interface.
func (f Filter) Tag() event.Tag { return f.Target }

// Matches implements the plugin.Filter interface.
func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case OpenFileEvent:
		return true
	case SaveFileEvent:
		return true
	case ErrorEvent:
		return true
	case CancelEvent:
		return true
	default:
		return false
	}
}
