package gioauth

import (
	"gioui.org/io/transfer"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"reflect"

	"gioui.org/io/event"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	wantOps = []reflect.Type{
		reflect.TypeOf(new(RequestOp)),
		reflect.TypeOf(new(ListenOp)),
	}
)

// RequestOp is used to request an authentication, and it will be responded with an AuthEvent to ListenOp,
// the Tag is used to receive errors, and the Provider is used to identify the provider.
type RequestOp struct {
	Tag      event.Tag
	Provider providers.Identifier
	Nonce    string
}

var poolRequestOp = plugin.NewOpPool[RequestOp]()

// Add adds the operation to op.Ops.
func (o RequestOp) Add(ops *op.Ops) {
	poolRequestOp.WriteOp(ops, o)
}

// ListenOp is used to receive an authentication, and it will be responded with an AuthEvent,
// to the same tag.
type ListenOp struct {
	Tag event.Tag
}

var poolListenOp = plugin.NewOpPool[ListenOp]()

// Add adds the operation to op.Ops.
func (o ListenOp) Add(ops *op.Ops) {
	transfer.SchemeOp{Tag: &DefaultProviders}.Add(ops)
	poolListenOp.WriteOp(ops, o)
}
