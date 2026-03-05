package gioinapppay

import (
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/inapppay"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var intName = plugin.NewIntName("gioinapppay")

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &inAppPayPlugin{window: w, plugin: handler}
	})
}

type inAppPayPlugin struct {
	window *app.Window
	plugin *plugin.Plugin
	config inapppay.Config
	client *inapppay.InAppPay
}

func (p *inAppPayPlugin) TypeOp() []reflect.Type          { return nil }
func (p *inAppPayPlugin) TypeCommand() []reflect.Type     { return wantCommands }
func (p *inAppPayPlugin) TypeEvent() []reflect.Type       { return wantEvents }
func (p *inAppPayPlugin) TypeGlobalEvent() []reflect.Type { return nil }
func (p *inAppPayPlugin) Op(op interface{})               {}
func (p *inAppPayPlugin) GlobalEvent(evt event.Event)     {}

func (p *inAppPayPlugin) Execute(cmd interface{}) {
	if p.client == nil {
		return
	}
	switch c := cmd.(type) {
	case ListProductsCmd:
		if err := p.client.ListProducts(c.ProductIDs); err != nil {
			p.plugin.SendEventUntagged(intName, ErrorEvent{Error: err})
		}
	case PurchaseCmd:
		if err := p.client.Purchase(c.ProductID, c.CustomPayload, c.IsPersonalizedPrice); err != nil {
			p.plugin.SendEventUntagged(intName, ErrorEvent{Error: err})
		}
	}
}

func (p *inAppPayPlugin) Event(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		if p.client == nil {
			p.config = NewConfigFromViewEvent(p.window, evt)
			p.client = inapppay.NewInAppPay(p.config)

			evts := p.client.Events()
			go func() {
				for evt := range evts {
					switch evt := evt.(type) {
					case inapppay.ProductDetailsEvent:
						p.plugin.SendEventUntagged(intName, ProductDetailsEvent(evt))
					case inapppay.PaymentResultEvent:
						p.plugin.SendEventUntagged(intName, PaymentResultEvent(evt))
					case inapppay.ErrorEvent:
						p.plugin.SendEventUntagged(intName, ErrorEvent(evt))
					}
				}
			}()
		} else {
			UpdateConfigFromViewEvent(&p.config, p.window, evt)
			p.client.Configure(p.config)
		}
	}
}
