package plugin

// ViewEvent is issued when the view is ready.
type ViewEvent struct{}

func (ViewEvent) ImplementsEvent() {}

// EndFrameEvent is issued at the end of a frame.
type EndFrameEvent struct{}

func (EndFrameEvent) ImplementsEvent() {}
