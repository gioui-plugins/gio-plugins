package giowebview

import "gioui.org/io/event"

// Filter is a filter for webview events.
type Filter struct {
	Target event.Tag
}

// ImplementsFilter implements event.Filter interface.
func (f Filter) ImplementsFilter() {}

// Tag implements plugin.Filter interface.
func (f Filter) Tag() event.Tag {
	return f.Target
}

// Matches implements plugin.Filter interface.
func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case StorageEvent:
		return true
	case CookiesEvent:
		return true
	case MessageEvent:
		return true
	case NavigationEvent:
		return true
	case TitleEvent:
		return true
	case ErrorEvent:
		return true
	case ClearCacheEvent:
		return true
	case SetCookieEvent:
		return true
	case SetCookieArrayEvent:
		return true
	case RemoveCookieEvent:
		return true
	case SetStorageEvent:
		return true
	case RemoveStorageEvent:
		return true
	default:
		return false
	}
}
