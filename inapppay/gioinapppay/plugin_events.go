package gioinapppay

import (
	"reflect"

	"github.com/gioui-plugins/gio-plugins/inapppay"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	wantEvents = []reflect.Type{
		reflect.TypeOf(plugin.ViewEvent{}),
		reflect.TypeOf(plugin.EndFrameEvent{}),
	}
)

// ErrorEvent is issued when an error occurs.
type ErrorEvent inapppay.ErrorEvent

// ProductDetailsEvent is issued when products are loaded.
type ProductDetailsEvent inapppay.ProductDetailsEvent

// PaymentResultEvent is issued when a purchase completes/fails.
type PaymentResultEvent inapppay.PaymentResultEvent

func (e ErrorEvent) ImplementsEvent()          {}
func (e ProductDetailsEvent) ImplementsEvent() {}
func (e PaymentResultEvent) ImplementsEvent()  {}
