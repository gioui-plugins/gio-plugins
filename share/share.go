// SPDX-License-Identifier: Unlicense OR MIT

package share

import (
	"errors"
	"reflect"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(app.ViewEvent{}),
		reflect.TypeOf(system.DestroyEvent{}),
		reflect.TypeOf(system.FrameEvent{}),
	}
	wantOp = []reflect.Type{
		reflect.TypeOf(new(WebsiteOp)),
		reflect.TypeOf(new(TextOp)),
	}
)

func init() {
	plugin.Register(func(w *app.Window, plugin *plugin.Plugin) plugin.Handler {
		return &sharePlugin{plugin: plugin, share: newShare(w)}
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
	plugin *plugin.Plugin

	// share holds OS-Specific content, it varies for each OS.
	share
}

// ListenEvents implements plugin.Handler.
func (e *sharePlugin) ListenEvents(evt event.Event) {
	if e == nil {
		return
	}
	e.listenEvents(evt)
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
		err = e.shareWebsite(*o)
		tag = o.Tag
	case *TextOp:
		err = e.shareText(*o)
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
