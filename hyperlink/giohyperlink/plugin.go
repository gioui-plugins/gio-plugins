package giohyperlink

import (
	"github.com/gioui-plugins/gio-plugins/hyperlink"
	"reflect"

	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/plugin"

	"gioui.org/io/event"
)

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &hyperlinkPlugin{window: w, plugin: handler}
	})
}

type hyperlinkPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	config hyperlink.Config
	client *hyperlink.Hyperlink
}

// TypeOp implements plugin.Handler.
func (h *hyperlinkPlugin) TypeOp() []reflect.Type { return nil }

// TypeCommand implements plugin.Handler.
func (h *hyperlinkPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements plugin.Handler.
func (h *hyperlinkPlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op implements plugin.Handler.
func (h *hyperlinkPlugin) Op(op interface{}) {}

// Execute implements plugin.Handler.
func (h *hyperlinkPlugin) Execute(op any) {
	switch op := op.(type) {
	case OpenCmd:
		if err := h.client.OpenWith(op.URI, op.PreferredPackage); err != nil {
			h.plugin.SendEvent(op.Tag, ErrorEvent{err})
		}
	}
}

// Event implements plugin.Handler.
func (h *hyperlinkPlugin) Event(evt event.Event) {
	if h == nil {
		return
	}

	switch evt := evt.(type) {
	case app.ViewEvent:
		UpdateConfigFromViewEvent(&h.config, h.window, evt)
		if h.client == nil {
			h.client = hyperlink.NewHyperlink(h.config)
		} else {
			h.client.Configure(h.config)
		}
	case app.ConfigEvent:
		UpdateConfigFromConfigEvent(&h.config, h.window, evt)
		if h.client != nil {
			h.client.Configure(h.config)
		}
	}
}
