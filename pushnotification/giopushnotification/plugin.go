package giopushnotification

import (
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/pushnotification"
)

// DefaultProviders is a list of providers that are used by default,
// you must append your own providers to this list.
//
// It MUST be modified before the plugin is initialized, so you
// need to do it on init() function or before init Gio.
var DefaultProviders []pushnotification.ExternalConfig

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &notificationPlugin{w: w, plugin: handler}
	})
}

// Plugin implements plugin.Handler.
type notificationPlugin struct {
	w      *app.Window
	plugin *plugin.Plugin

	config pushnotification.Config
	client *pushnotification.Push
}

// TypeOp returns the operations that are redirected to this plugin.
func (p *notificationPlugin) TypeOp() []reflect.Type { return []reflect.Type{} }

// TypeCommand returns the commands that are redirected to this plugin.
func (p *notificationPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent returns the events that are redirected to this plugin.
func (p *notificationPlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op handles the operation.
func (p *notificationPlugin) Op(op interface{}) {}

// Execute executes the command.
func (p *notificationPlugin) Execute(cmd interface{}) {
	switch c := cmd.(type) {
	case RequestTokenCmd:
		go func() {
			token, err := p.client.RequestToken()
			if err != nil {
				p.plugin.SendEvent(c.Tag, ErrorEvent{Error: err})
				return
			}
			p.plugin.SendEvent(c.Tag, TokenReceivedEvent(token))
		}()
	}
}

// Event handles the event.
func (p *notificationPlugin) Event(evt event.Event) {
	switch e := evt.(type) {
	case app.ViewEvent:
		UpdateConfigFromViewEvent(&p.config, p.w, e, DefaultProviders)
		if p.client == nil {
			p.client = pushnotification.NewPush(p.config)
		} else {
			p.client.Configure(p.config)
		}
	}
}

// Push returns the underlying Push instance.
func (p *notificationPlugin) Push() *pushnotification.Push {
	return p.client
}
