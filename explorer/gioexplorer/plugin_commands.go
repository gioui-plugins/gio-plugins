package gioexplorer

import (
	"errors"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
	"reflect"
	"strings"
)

var wantCommands = []reflect.Type{
	reflect.TypeOf(OpenFileCmd{}),
	// reflect.TypeOf(&OpenDirectoryOp{}),
	reflect.TypeOf(SaveFileCmd{}),
	// reflect.TypeOf(&SaveDirectoryOp{}),
}

// OpenFileCmd opens the file selector and returns the selected file.
// The Mimetype may filter the files that can be selected.
type OpenFileCmd struct {
	Tag      event.Tag
	Mimetype []mimetype.MimeType
}

func (o OpenFileCmd) ImplementsCommand() {}

func (o OpenFileCmd) execute(p *explorerPlugin) {
	go func() {
		res, err := p.client.OpenFile(o.Mimetype)
		if err != nil {
			if errors.Is(err, ErrUserDecline) {
				p.plugin.SendEvent(o.Tag, CancelEvent{})
			} else {
				p.plugin.SendEvent(o.Tag, ErrorEvent{error: err})
			}
			return
		}

		p.plugin.SendEvent(o.Tag, OpenFileEvent{File: res})
	}()
}

// SaveFileCmd opens the file-picker to save a file, the file is created if it doesn't exist, or replace existent file.
// The Filename is a suggestion for the file name, the user can change it.
type SaveFileCmd struct {
	Tag      event.Tag
	Filename string
	Mimetype mimetype.MimeType
}

func (o SaveFileCmd) ImplementsCommand() {}

func (o SaveFileCmd) execute(p *explorerPlugin) {
	go func() {
		if strings.HasPrefix(o.Filename, "."+o.Mimetype.Extension) {
			o.Filename = o.Filename[:len(o.Filename)-(1+len(o.Mimetype.Extension))]
		}

		res, err := p.client.SaveFile(o.Filename, o.Mimetype)
		if err != nil {
			if errors.Is(err, ErrUserDecline) {
				p.plugin.SendEvent(o.Tag, CancelEvent{})
			} else {
				p.plugin.SendEvent(o.Tag, ErrorEvent{error: err})
			}
			return
		}

		p.plugin.SendEvent(o.Tag, SaveFileEvent{File: res})
	}()
}
