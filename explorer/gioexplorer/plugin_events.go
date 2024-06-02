package gioexplorer

import (
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"io"
	"reflect"
)

var wantEvents = []reflect.Type{
	reflect.TypeOf(plugin.ViewEvent{}),
}

// OpenFileEvent is sent as response to OpenFileOp.
type OpenFileEvent struct {
	Tag  event.Tag
	File io.ReadCloser
}

// SaveFileEvent is sent as response to SaveFileOp.
type SaveFileEvent struct {
	Tag  event.Tag
	File io.WriteCloser
}

// ErrorEvent is issued when error occurs.
type ErrorEvent struct {
	error
}

// CancelEvent is sent when the user cancels the file selector.
type CancelEvent struct{}

func (OpenFileEvent) ImplementsEvent() {}
func (SaveFileEvent) ImplementsEvent() {}
func (ErrorEvent) ImplementsEvent()    {}
func (CancelEvent) ImplementsEvent()   {}
