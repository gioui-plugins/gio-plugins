package gioauth

import (
	"gioui.org/io/system"
	"gioui.org/io/transfer"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"reflect"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var DefaultProviders []providers.Provider

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &authPlugin{
			window: w,
			plugin: handler,
		}
	})
}

type authPlugin struct {
	window *app.Window
	plugin *plugin.Plugin
	config auth.Config

	startupURL string
	client     *auth.Auth

	mutex sync.Mutex

	current []event.Tag
	targets []event.Tag
	ready   bool
}

// TypeOp implements plugin.Handler.
func (p *authPlugin) TypeOp() []reflect.Type { return wantOps }

// TypeEvent implements plugin.Handler.
func (p *authPlugin) TypeEvent() []reflect.Type { return wantEvents }

// ListenOps implements plugin.Handler.
func (p *authPlugin) ListenOps(op interface{}) {
	switch o := op.(type) {
	case *RequestOp:
		defer poolRequestOp.Release(o)

		if err := p.client.Open(o.Provider, o.Nonce); err != nil {
			p.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		}

	case *ListenOp:
		defer poolListenOp.Release(o)
		p.current = append(p.current, o.Tag)
	}
}

// ListenEvents implements plugin.Handler.
func (p *authPlugin) ListenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case system.FrameEvent:
		for _, evt := range evt.Queue.Events(&DefaultProviders) {
			switch evt := evt.(type) {
			case transfer.URLEvent:
				if p.ready && p.client != nil {
					p.client.ProcessCustomSchemeCallback(evt.URL.String())
				} else {
					p.startupURL = evt.URL.String()
				}
			}
		}
	case app.ViewEvent:
		if p.client == nil {
			p.config = NewConfigFromViewEvent(p.window, evt)
			p.client = auth.NewAuth(p.config, DefaultProviders...)
			go func() {
				for evt := range p.client.Events() {
					var resp event.Event
					switch evt := evt.(type) {
					case auth.ErrorEvent:
						resp = ErrorEvent(evt)
					case auth.AuthenticatedEvent:
						resp = AuthEvent(evt)
					default:
						continue
					}

					p.mutex.Lock()
					for _, tag := range p.targets {
						p.plugin.SendEvent(tag, resp)
					}
					p.mutex.Unlock()
				}
			}()
		} else {
			UpdateConfigFromViewEvent(&p.config, p.window, evt)
			p.client.Configure(p.config)
		}

	case plugin.EndFrameEvent:
		p.mutex.Lock()
		defer p.mutex.Unlock()

		if len(p.current) > cap(p.targets) {
			p.targets = make([]event.Tag, 0, len(p.current))
		}

		p.targets = p.targets[:len(p.current)]
		copy(p.targets, p.current)
		p.current = p.current[:0]

		if !p.ready && len(p.targets) > 0 && p.client != nil {
			p.ready = true

			if p.startupURL != "" {
				p.client.ProcessCustomSchemeCallback(p.startupURL)
				p.startupURL = ""
			}
		}
	}
}
