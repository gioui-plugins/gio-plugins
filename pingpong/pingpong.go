package pingpong

import (
	"gioui.org/op"
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

func init() {
	plugin.Register(func(w *app.Window, handler *plugin.Plugin) plugin.Handler {
		return &pingPong{w: w, p: handler}
	})
}

type pingPong struct {
	w *app.Window
	p *plugin.Plugin
}

func (p *pingPong) TypeOp() []reflect.Type      { return []reflect.Type{reflect.TypeOf(&PingOp{})} }
func (p *pingPong) TypeCommand() []reflect.Type { return []reflect.Type{reflect.TypeOf(PingCmd{})} }
func (p *pingPong) TypeEvent() []reflect.Type   { return nil }

func (p *pingPong) Op(op interface{}) {
	if ping, ok := op.(*PingOp); ok {
		defer _PingOpPool.Release(ping)
		p.p.SendEvent(ping.Tag, PongEvent{Text: ping.Text})
	}
}

func (p *pingPong) Execute(cmd interface{}) {
	if ping, ok := cmd.(PingCmd); ok {
		p.p.SendEvent(ping.Tag, PongEvent{Text: ping.Text})
	}
}

func (p *pingPong) Event(evt event.Event) {}

// PingOp is an operation that sends a PongEvent to the given tag.
type PingOp struct {
	Tag  event.Tag
	Text string
}

var _PingOpPool = plugin.NewOpPool[PingOp]()

// Add writes the operation to the given op.Ops.
func (o PingOp) Add(ops *op.Ops) {
	opc := _PingOpPool.Get()
	*opc = o
	plugin.WriteOp(ops, opc)
}

// PingCmd is a command that sends a PongEvent to the given tag.
type PingCmd struct {
	Tag  event.Tag
	Text string
}

func (o PingCmd) ImplementsCommand() {}

// PongEvent is the event that is sent back once PingOp is invoked.
type PongEvent struct {
	Text string
}

func (p PongEvent) ImplementsEvent() {}

// Filter returns true if the event is a PongEvent.
type Filter struct {
	Target event.Tag
}

func (f Filter) ImplementsFilter() {}

func (f Filter) Tag() event.Tag {
	return f.Target
}

func (f Filter) Matches(e event.Event) bool {
	switch e.(type) {
	case PongEvent:
		return true
	default:
		return false
	}
}
