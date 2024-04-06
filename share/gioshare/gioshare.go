package gioshare

import (
	"errors"
	"reflect"
	"sync"

	"github.com/gioui-plugins/gio-plugins/share"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	wantEvents = []reflect.Type{
		// reflect.TypeOf(app.ViewEvent{}),
		reflect.TypeOf(app.DestroyEvent{}),
		reflect.TypeOf(app.FrameEvent{}),
	}
	wantOp = []reflect.Type{
		reflect.TypeOf(new(WebsiteOp)),
		reflect.TypeOf(new(TextOp)),
	}
)

func init() {
	plugin.Register(func(w *app.Window, plugin *plugin.Plugin) plugin.Handler {
		return &sharePlugin{plugin: plugin, window: w}
	})
}

var (
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")

	// ErrNotAvailableAction is return when the current Shareable item isn't supported.
	ErrNotAvailableAction = errors.New("current shareable item not supported")
)

type sharePlugin struct {
	mutex  sync.Mutex
	window *app.Window
	plugin *plugin.Plugin

	config share.Config
	client *share.Share
}

// ListenEvents implements plugin.Handler.
func (e *sharePlugin) ListenEvents(evt event.Event) {
	if e == nil {
		return
	}
	switch evt := evt.(type) {
	case app.FrameEvent:
		UpdateConfigFromFrameEvent(&e.config, e.window, evt)
		if e.client != nil {
			e.client.Configure(e.config)
		}
	case app.ViewEvent:
		UpdateConfigFromViewEvent(&e.config, e.window, evt)
		if e.client == nil {
			e.client = share.NewShare(e.config)
		} else {
			e.client.Configure(e.config)
		}
	}
}

// ListenOps implements plugin.Handler.
func (e *sharePlugin) ListenOps(op interface{}) {
	var (
		err error
		tag event.Tag
	)
	switch o := op.(type) {
	case *WebsiteOp:
		defer websiteOpPool.Release(o)
		err = e.client.Website(o.Title, o.Text, o.Link)
		tag = o.Tag
	case *TextOp:
		defer textOpPool.Release(o)
		err = e.client.Text(o.Title, o.Text)
		tag = o.Tag
	}

	if err != nil {
		e.plugin.SendEvent(tag, ErrorEvent{error: err})
	}
}

// TypeEvent implements plugin.Handler.
func (e *sharePlugin) TypeEvent() []reflect.Type {
	return wantEvents
}

// TypeOp implements plugin.Handler.
func (e *sharePlugin) TypeOp() []reflect.Type {
	return wantOp
}

type Shareable interface {
	ImplementsShareable()
}

// TextOp represents the text to be shared.
type TextOp struct {
	Tag   event.Tag
	Title string
	Text  string
}

var textOpPool = plugin.NewOpPool[TextOp]()

func (o TextOp) Add(ops *op.Ops) {
	textOpPool.WriteOp(ops, o)
}

// WebsiteOp represents the website/link to be shared.
type WebsiteOp struct {
	Tag   event.Tag
	Title string
	Text  string
	Link  string
}

var websiteOpPool = plugin.NewOpPool[WebsiteOp]()

// Add adds the WebsiteOp to the queue.
// The WebsiteOp will be executed on the next frame, showing the share dialog.
func (o WebsiteOp) Add(ops *op.Ops) {
	websiteOpPool.WriteOp(ops, o)
}

func (o TextOp) ImplementsShareable()    {}
func (o WebsiteOp) ImplementsShareable() {}

// ErrorEvent is issued when an error occurs.
type ErrorEvent struct {
	error
}

func (ErrorEvent) ImplementsEvent() {}
