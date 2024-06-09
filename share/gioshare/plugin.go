package gioshare

import (
	"github.com/gioui-plugins/gio-plugins/share"
	"reflect"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

func init() {
	plugin.Register(func(w *app.Window, plugin *plugin.Plugin) plugin.Handler {
		return &sharePlugin{plugin: plugin, window: w}
	})
}

type sharePlugin struct {
	mutex  sync.Mutex
	window *app.Window
	plugin *plugin.Plugin

	config share.Config
	client *share.Share
}

// TypeOp implements plugin.Handler.
func (e *sharePlugin) TypeOp() []reflect.Type { return nil }

// TypeCommand implements plugin.Handler.
func (e *sharePlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements plugin.Handler.
func (e *sharePlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op implements plugin.Handler.
func (e *sharePlugin) Op(op interface{}) {}

func (e *sharePlugin) Execute(op interface{}) {
	var (
		err error
		tag event.Tag
	)
	switch o := op.(type) {
	case WebsiteCmd:
		err = e.client.Website(o.Title, o.Text, o.Link)
		tag = o.Tag
	case TextCmd:
		err = e.client.Text(o.Title, o.Text)
		tag = o.Tag
	}

	if err != nil {
		e.plugin.SendEvent(tag, ErrorEvent{Error: err})
	}
}

// Event implements plugin.Handler.
func (e *sharePlugin) Event(evt event.Event) {
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
