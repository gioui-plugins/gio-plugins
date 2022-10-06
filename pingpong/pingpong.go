package pingpong

import (
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
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

func (p *pingPong) TypeOp() []reflect.Type    { return []reflect.Type{reflect.TypeOf(&PingOp{})} }
func (p *pingPong) TypeEvent() []reflect.Type { return nil }

func (p *pingPong) ListenOps(op interface{}) {
	if ping, ok := op.(*PingOp); ok {
		defer pingOpPool.Release(ping)
		p.p.SendEvent(ping.Tag, PongEvent{Text: ping.Text})
	}
}

func (p *pingPong) ListenEvents(evt event.Event) {}

// PingOp is an operation that sends a PongEvent to the given tag.
type PingOp struct {
	Tag  event.Tag
	Text string
}

var pingOpPool = plugin.NewOpPool[PingOp]()

// Add writes the operation to the given op.Ops.
func (o PingOp) Add(ops *op.Ops) { pingOpPool.WriteOp(ops, o) }

// PongEvent is the event that is sent back once PingOp is invoked.
type PongEvent struct {
	Text string
}

func (p PongEvent) ImplementsEvent() {}
