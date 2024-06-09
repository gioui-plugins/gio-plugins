package giosafedata

import (
	"errors"
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/safedata"
)

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &safedataPlugin{
			window: w,
			plugin: handler,
		}
	})
}

var (
	// ErrNotReady may occur when try to open a URL before the initialization is done.
	ErrNotReady = errors.New("some needed library was not loaded yet, make use that you are using ListenEvents()")

	// ErrInvalidURL occur when provide an invalid URL, like a non http/https URL.
	ErrInvalidURL = errors.New("given url is invalid")
)

type safedataPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	client *safedata.SafeData
}

// TypeOp implements plugin.Handler.
func (p *safedataPlugin) TypeOp() []reflect.Type { return nil }

// TypeCommand implements plugin.Handler.
func (p *safedataPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements plugin.Handler.
func (p *safedataPlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op implements plugin.Handler.
func (p *safedataPlugin) Op(op interface{}) {}

// Execute implements plugin.Handler.
func (p *safedataPlugin) Execute(op interface{}) {
	switch op := op.(type) {
	case WriteSecretCmd:
		op.execute(p)
	case ReadSecretCmd:
		op.execute(p)
	case DeleteSecretCmd:
		op.execute(p)
	case ListSecretCmd:
		op.execute(p)
	}
}

// Event implements plugin.Handler.
func (p *safedataPlugin) Event(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		config := NewConfigFromViewEvent(p.window, evt, safedata.DefaultAppName)
		if p.client == nil {
			p.client = safedata.NewSafeData(config)
		} else {
			p.client.Configure(config)
		}
	}
}
