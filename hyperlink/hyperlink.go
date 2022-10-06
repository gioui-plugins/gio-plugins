package hyperlink

import (
	"errors"
	"net/url"
	"reflect"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"

	"gioui.org/io/event"
)

var (
	wantOps = []reflect.Type{
		reflect.TypeOf(&OpenOp{}),
	}
	wantEvents = []reflect.Type{
		reflect.TypeOf(app.ViewEvent{}),
		reflect.TypeOf(system.StageEvent{}),
	}
)

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &hyperlinkPlugin{window: w, plugin: handler}
	})
}

var (
	// ErrNotReady may occur when try to open a URL before the initialization is done.
	ErrNotReady = errors.New("some needed library was not loaded yet, make use that you are using ListenEvents()")
	// ErrInvalidURL occur when provide an invalid URL, like a non http/https URL.
	ErrInvalidURL = errors.New("given url is invalid")
)

var (
	// InsecureIgnoreScheme will remove any attempt to validate the URL
	// It's "false" by default. Set it to "true" if you are using a custom scheme (like "myapp://").
	InsecureIgnoreScheme bool
)

type hyperlinkPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	// hyperlink is OS-specific data
	hyperlink
}

// TypeOp implements plugin.Handler.
func (h *hyperlinkPlugin) TypeOp() []reflect.Type {
	return wantOps
}

// TypeEvent implements plugin.Handler.
func (h *hyperlinkPlugin) TypeEvent() []reflect.Type {
	return wantEvents
}

// ListenOps implements plugin.Handler.
func (h *hyperlinkPlugin) ListenOps(op interface{}) {
	switch op := op.(type) {
	case *OpenOp:
		defer openOpPool.Release(op)

		if op.URI == nil || op.URI.Scheme == "" || ((op.URI.Scheme != "http" && op.URI.Scheme != "https") && InsecureIgnoreScheme == false) {
			h.plugin.SendEvent(op.Tag, event.Event(ErrorEvent{ErrInvalidURL}))
			return
		}

		if err := h.open(op.URI); err != nil {
			h.plugin.SendEvent(op.Tag, event.Event(ErrorEvent{err}))
		}
	}
}

// ListenEvents implements plugin.Handler.
func (h *hyperlinkPlugin) ListenEvents(evt event.Event) {
	h.listenEvents(evt)
}

// OpenOp is an operation that will open a URL.
type OpenOp struct {
	Tag event.Tag
	URI *url.URL
}

var openOpPool = plugin.NewOpPool[OpenOp]()

// Add adds an OpenOp to the queue.
// It will open the given URL at the end of the frame.
func (o OpenOp) Add(ops *op.Ops) {
	openOpPool.WriteOp(ops, o)
}

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	error
}

// ImplementsEvent implements event.Event.
func (e ErrorEvent) ImplementsEvent() {}
