package gioauth

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"
)

// DefaultProviders is a list of providers that are used by default,
// you must append your own providers to this list.
//
// It MUST be modified before the plugin is initialized, so you
// need to do it on init() function or before init Gio.
var DefaultProviders []providers.Provider

var intName = plugin.NewIntName("gioauth")

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &authPlugin{window: w, plugin: handler, startupURL: startURL()}
	})
}

type authPlugin struct {
	window *app.Window
	plugin *plugin.Plugin
	config auth.Config

	startupURL string
	client     *auth.Auth
}

// TypeOp implements plugin.Handler.
func (p *authPlugin) TypeOp() []reflect.Type { return nil }

// TypeCommand implements plugin.Handler.
func (p *authPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements plugin.Handler.
func (p *authPlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op implements plugin.Handler.
func (p *authPlugin) Op(op interface{}) {}

// Execute implements plugin.Handler.
func (p *authPlugin) Execute(cmd interface{}) {
	switch o := cmd.(type) {
	case OpenCmd:
		if err := p.client.Open(o.Provider, o.Nonce); err != nil {
			p.plugin.SendEventUntagged(intName, ErrorEvent(auth.ErrorEvent{Error: err}))
		}
	}
}

// Event implements plugin.Handler.
func (p *authPlugin) Event(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		if p.client == nil {
			p.config = NewConfigFromViewEvent(p.window, evt)
			p.client = auth.NewAuth(p.config, DefaultProviders...)

			evts := p.client.Events()
			go func() {
				for evt := range evts {
					switch evt := evt.(type) {
					case auth.ErrorEvent:
						p.plugin.SendEventUntagged(intName, ErrorEvent(evt))
					case auth.AuthenticatedEvent:
						p.plugin.SendEventUntagged(intName, AuthEvent(evt))
					}
				}
			}()
		} else {
			UpdateConfigFromViewEvent(&p.config, p.window, evt)
			p.client.Configure(p.config)
		}
	case plugin.EndFrameEvent:
		if p.startupURL != "" {
			if err := p.client.ProcessCustomSchemeCallback(p.startupURL); err != nil {
				p.plugin.SendEventUntagged(intName, ErrorEvent(auth.ErrorEvent{Error: err}))
			}
			p.startupURL = ""
		}
	}
}

func (p *authPlugin) TypeGlobalEvent() []reflect.Type {
	return wantGlobalEvents
}

func (p *authPlugin) GlobalEvent(evt event.Event) {
	/*
		switch evt := evt.(type) {

			case app.URLEvent:
				if p.client != nil {
					if err := p.client.ProcessCustomSchemeCallback(evt.URL.String()); err != nil {
						p.plugin.SendEventUntagged(intName, ErrorEvent(auth.ErrorEvent{Error: err}))
					}
				} else {
					p.startupURL = evt.URL.String()
				}
		}
	*/
}
