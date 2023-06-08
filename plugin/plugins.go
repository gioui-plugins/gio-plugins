package plugin

import (
	"reflect"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/pointer"
	"gioui.org/op"
	"gioui.org/op/clip"
)

// Handler is the interface that represents the Plugin.
type Handler interface {
	TypeOp() []reflect.Type
	TypeEvent() []reflect.Type

	ListenOps(op interface{})
	ListenEvents(evt event.Event)
}

var registeredPlugins []func(w *app.Window, handler *Plugin) Handler

// Register registers the Handler, it will be called when the window is created.
// You MUST call Register during init() function, otherwise it will not work or
// may cause unexpected behavior.
func Register(plugin func(w *app.Window, handler *Plugin) Handler) {
	registeredPlugins = append(registeredPlugins, plugin)
}

type handlerFunc struct {
	typeOp       []reflect.Type
	typeEvent    []reflect.Type
	listenOps    func(op interface{})
	listenEvents func(evt event.Event)
}

// TypeOp implements Handler interface.
func (p *handlerFunc) TypeOp() []reflect.Type { return p.typeOp }

// TypeEvent implements Handler interface.
func (p *handlerFunc) TypeEvent() []reflect.Type { return p.typeEvent }

// ListenOps implements Handler interface.
func (p *handlerFunc) ListenOps(op interface{}) { p.listenOps(op) }

// ListenEvents implements Handler interface.
func (p *handlerFunc) ListenEvents(evt event.Event) { p.listenEvents(evt) }

// NewHandlerFunc returns a Handler that calls the given functions.
func NewHandlerFunc(ops []reflect.Type, events []reflect.Type, listenOp func(op interface{}), listenEvents func(event.Event)) func(w *app.Window, h *Plugin) Handler {
	if listenOp == nil {
		listenOp = func(interface{}) {}
	}
	if listenEvents == nil {
		listenEvents = func(event.Event) {}
	}
	return func(w *app.Window, h *Plugin) Handler {
		return &handlerFunc{
			typeOp:       ops,
			typeEvent:    events,
			listenOps:    listenOp,
			listenEvents: listenEvents,
		}
	}
}

// WriteOp writes the given op into the op.Ops queue.
func WriteOp(op *op.Ops, c any) {
	defer clip.Rect{}.Push(op).Pop()
	pointer.InputOp{Tag: c}.Add(op)
}

// OpPool is a pool of specific type of op.
type OpPool[T any] struct {
	pool  sync.Pool
	empty T
}

// NewOpPool returns a new OpPool.
// That is useful to avoid memory allocation, you MUST
// call Release() after you done using the op.
func NewOpPool[T any]() OpPool[T] {
	return OpPool[T]{pool: sync.Pool{New: func() any { return new(T) }}}
}

// Get returns a new op from the pool.
func (x *OpPool[T]) Get() *T {
	return x.pool.Get().(*T)
}

// WriteOp adds the given data into the op.Ops queue.
func (x *OpPool[T]) WriteOp(op *op.Ops, data T) {
	cmd := x.pool.New()
	*cmd.(*T) = data

	WriteOp(op, cmd)
}

// Release releases the given data, so it can be reused.
func (x *OpPool[T]) Release(data *T) {
	*data = x.empty
	x.pool.Put(data)
}
