package plugin

import (
	"encoding/binary"
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
	"gioui.org/op/clip"
	"golang.org/x/crypto/blake2b"
	"reflect"
)

// Handler is the interface that represents the Plugin.
type Handler interface {
	// TypeOp returns the list of Ops that the plugin can handle.
	// Op are data that are encoded into op.Ops queue from op.Ops.Add.
	TypeOp() []reflect.Type

	// TypeCommand returns the list of commands that the plugin can handle.
	// Command are data that are sent via Execute method, from input.Router.Execute.
	TypeCommand() []reflect.Type

	// TypeEvent returns the list of events that the plugin can handle and
	// are interested in.
	// Event are data that received from app.Window.Event.
	TypeEvent() []reflect.Type

	// Op is called when an op is sent to the plugin.
	Op(op interface{})

	// Execute is called when a command is sent to the plugin.
	Execute(cmd interface{})

	// Event is called when an event is sent to the plugin.
	Event(evt event.Event)
}

// Filter is used to filter events, extends event.Filter.
type Filter interface {
	event.Filter

	// Tag returns the event.Tag that the Filter is interested in.
	Tag() event.Tag

	// Matches returns true if the event matches the Filter.
	Matches(event.Event) bool
}

// UntaggedFilter is used to filter events, without using event.Tag.
type UntaggedFilter interface {
	event.Filter

	// Matches returns true if the event matches the Filter.
	Matches(event.Event) bool

	// Name returns the name of the filter.
	Name() uint64
}

// NewIntName returns a uint64 from the given name,
// it is used to generate a unique name for the untagged filter and events.
func NewIntName(name string) uint64 {
	h, _ := blake2b.New(8, nil)
	h.Write([]byte(name))
	return binary.BigEndian.Uint64(h.Sum(nil))
}

var registeredPlugins []func(w *app.Window, handler *Plugin) Handler

// Register registers the Handler, it will be called when the window is created.
//
// You MUST call Register during init() function, otherwise it will not work or
// may cause unexpected behavior.
func Register(plugin func(w *app.Window, handler *Plugin) Handler) {
	registeredPlugins = append(registeredPlugins, plugin)
}

type handlerFunc struct {
	typeOp       []reflect.Type
	typeCmd      []reflect.Type
	typeEvent    []reflect.Type
	listenOp     func(op interface{})
	listenCmd    func(cmd interface{})
	listenEvents func(evt event.Event)
}

// TypeOp implements Handler interface.
func (p *handlerFunc) TypeOp() []reflect.Type { return p.typeOp }

// TypeCommand implements Handler interface.
func (p *handlerFunc) TypeCommand() []reflect.Type { return p.typeCmd }

// TypeEvent implements Handler interface.
func (p *handlerFunc) TypeEvent() []reflect.Type { return p.typeEvent }

// Op implements Handler interface.
func (p *handlerFunc) Op(op interface{}) { p.listenOp(op) }

// Execute implements Handler interface.
func (p *handlerFunc) Execute(op interface{}) { p.listenCmd(op) }

// Event implements Handler interface.
func (p *handlerFunc) Event(evt event.Event) { p.listenEvents(evt) }

// NewHandlerFunc returns a Handler that calls the given functions.
func NewHandlerFunc(ops, cmd, evt []reflect.Type, listenOp, listenCmd func(op interface{}), listenEvents func(event.Event)) func(w *app.Window, h *Plugin) Handler {
	if listenCmd == nil {
		listenCmd = func(interface{}) {}
	}
	if listenEvents == nil {
		listenEvents = func(event.Event) {}
	}
	return func(w *app.Window, h *Plugin) Handler {
		return &handlerFunc{
			typeOp:       ops,
			typeCmd:      cmd,
			typeEvent:    evt,
			listenOp:     listenOp,
			listenCmd:    listenCmd,
			listenEvents: listenEvents,
		}
	}
}

// WriteOp writes the given op into the op.Ops queue.
func WriteOp(ops *op.Ops, c any) {
	defer clip.Rect{}.Push(ops).Pop()
	event.Op(ops, c)
}
