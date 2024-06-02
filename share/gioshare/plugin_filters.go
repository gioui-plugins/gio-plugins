package gioshare

import "gioui.org/io/event"

type Filter struct {
	Target event.Tag
}

func (f Filter) ImplementsFilter() {}

func (f Filter) Tag() event.Tag {
	return f.Target
}

func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case ErrorEvent:
		return true
	default:
		return false
	}
}
