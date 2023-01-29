package webview

// Event is the marker interface for events.
type Event interface {
	ImplementsEvent()
}

// NavigationEvent is issued when the URL is changed
// in the WebView.
type NavigationEvent struct {
	URL string
}

func (NavigationEvent) ImplementsEvent() {}

// TitleEvent is issued when the Title of the website
// is changed in the WebView.
type TitleEvent struct {
	Title string
}

func (TitleEvent) ImplementsEvent() {}
