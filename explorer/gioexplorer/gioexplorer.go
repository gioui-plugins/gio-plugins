package gioexplorer

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"sync"

	"github.com/gioui-plugins/gio-plugins/explorer"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	wantOps = []reflect.Type{
		reflect.TypeOf(&OpenFileOp{}),
		// reflect.TypeOf(&OpenDirectoryOp{}),
		reflect.TypeOf(&SaveFileOp{}),
		// reflect.TypeOf(&SaveDirectoryOp{}),
	}
	wantEvents = []reflect.Type{
		//reflect.TypeOf(app.ViewEvent{}),
	}
)

func init() {
	plugin.Register(func(w *app.Window, p *plugin.Plugin) plugin.Handler {
		return &explorerPlugin{window: w, plugin: p}
	})
}

var (
	// ErrUserDecline is returned when the user doesn't select the file.
	ErrUserDecline = errors.New("user exited the file selector without selecting a file")
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")
)

type explorerPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	config explorer.Config
	client *explorer.Explorer
}

func (e *explorerPlugin) TypeOp() []reflect.Type {
	return wantOps
}

func (e *explorerPlugin) TypeEvent() []reflect.Type {
	return wantEvents
}

func (e *explorerPlugin) ListenOps(op interface{}) {
	switch o := op.(type) {
	case *OpenFileOp:
		o.execute(e)
	case *SaveFileOp:
		o.execute(e)
	}
}

func (e *explorerPlugin) ListenEvents(evt event.Event) {
	if e == nil {
		return
	}
	switch evt := evt.(type) {
	case app.ViewEvent:
		UpdateConfigFromViewEvent(&e.config, e.window, evt)
		if e.client == nil {
			e.client = explorer.NewExplorer(e.config)
		} else {
			e.client.Configure(e.config)
		}
	}
}

// OpenFileOp opens the file selector and returns the selected file.
// The Mimetype may filter the files that can be selected.
type OpenFileOp struct {
	Tag      event.Tag
	Mimetype []mimetype.MimeType
}

var openFileOpPool = plugin.NewOpPool[OpenFileOp]()

// Add adds the operation into the queue.
func (o OpenFileOp) Add(ops *op.Ops) {
	oc := openFileOpPool.Get()
	oc.Tag = o.Tag
	oc.Mimetype = o.Mimetype
	if l := len(o.Mimetype) - cap(oc.Mimetype); l > 0 {
		oc.Mimetype = append(oc.Mimetype, make([]mimetype.MimeType, l)...)
	}
	oc.Mimetype = oc.Mimetype[:len(o.Mimetype)]
	copy(oc.Mimetype, o.Mimetype)
	plugin.WriteOp(ops, oc)
}

func (o *OpenFileOp) execute(p *explorerPlugin) {
	go func() {
		defer openFileOpPool.Release(o)

		res, err := p.client.OpenFile(o.Mimetype)
		if err != nil {
			if err == ErrUserDecline {
				p.plugin.SendEvent(o.Tag, CancelEvent{})
			} else {
				p.plugin.SendEvent(o.Tag, ErrorEvent{error: err})
			}
			return
		}

		p.plugin.SendEvent(o.Tag, OpenFileEvent{File: res})
	}()
}

// SaveFileOp opens the file-picker to save a file, the file is created if it doesn't exist, or replace existent file.
// The Filename is a suggestion for the file name, the user can change it.
type SaveFileOp struct {
	Tag      event.Tag
	Filename string
	Mimetype mimetype.MimeType
}

var saveFileOpPool = plugin.NewOpPool[SaveFileOp]()

// Add adds the event into the queue.
func (o SaveFileOp) Add(ops *op.Ops) {
	saveFileOpPool.WriteOp(ops, o)
}

func (o *SaveFileOp) execute(p *explorerPlugin) {
	go func() {
		defer saveFileOpPool.Release(o)

		if strings.HasPrefix(o.Filename, "."+o.Mimetype.Extension) {
			o.Filename = o.Filename[:len(o.Filename)-(1+len(o.Mimetype.Extension))]
		}

		res, err := p.client.SaveFile(o.Filename, o.Mimetype)
		if err != nil {
			if err == ErrUserDecline {
				p.plugin.SendEvent(o.Tag, CancelEvent{})
			} else {
				p.plugin.SendEvent(o.Tag, ErrorEvent{error: err})
			}
			return
		}
		p.plugin.SendEvent(o.Tag, SaveFileEvent{File: res})
	}()
}

// OpenFileEvent is sent as response to OpenFileOp.
type OpenFileEvent struct {
	File io.ReadCloser
}

func (e OpenFileEvent) ImplementsFilter() {
}

// SaveFileEvent is sent as response to SaveFileOp.
type SaveFileEvent struct {
	File io.WriteCloser
}

func (e SaveFileEvent) ImplementsFilter() {
}

// ErrorEvent is issued when error occurs.
type ErrorEvent struct {
	error
}

func (e ErrorEvent) ImplementsFilter() {
}

// CancelEvent is sent when the user cancels the file selector.
type CancelEvent struct{}

func (e CancelEvent) ImplementsFilter() {
}

func (OpenFileEvent) ImplementsEvent() {}
func (SaveFileEvent) ImplementsEvent() {}
func (ErrorEvent) ImplementsEvent()    {}
func (CancelEvent) ImplementsEvent()   {}

var stringBuilderPool = sync.Pool{New: func() any { return &strings.Builder{} }}
