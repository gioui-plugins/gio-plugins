package plugin

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/op"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"
)

type queue struct {
	taggedEvents   map[event.Tag][]event.Event
	untaggedEvents map[uint64][]event.Event
}

func newQueue() *queue {
	return &queue{
		taggedEvents:   make(map[event.Tag][]event.Event, 128),
		untaggedEvents: make(map[uint64][]event.Event, 128),
	}
}

// Plugin is the main interface for the plugins.
type Plugin struct {
	window *app.Window

	eventsCustomNextMutex    sync.Mutex
	eventsCustomCurrentMutex sync.Mutex

	// double buffered events
	eventsCustomNext    *queue
	eventsCustomCurrent *queue

	eventsPool []event.Event

	RedirectEvent    map[reflect.Type][]int
	RedirectOp       map[reflect.Type][]int
	RedirectCommands map[reflect.Type][]int

	visited map[uintptr]struct{}

	Plugins     []Handler
	Invalidated atomic.Bool

	OriginalFrame  func(ops *op.Ops)
	OriginalSource input.Source
}

// NewPlugin creates a new plugin.
func NewPlugin(w *app.Window) *Plugin {
	h := &Plugin{
		window:              w,
		Plugins:             make([]Handler, len(registeredPlugins)),
		visited:             make(map[uintptr]struct{}, 128),
		RedirectOp:          make(map[reflect.Type][]int, 128),
		RedirectCommands:    make(map[reflect.Type][]int, 128),
		RedirectEvent:       make(map[reflect.Type][]int, 128),
		eventsCustomNext:    newQueue(),
		eventsCustomCurrent: newQueue(),
	}

	for index, pf := range registeredPlugins {
		h.Plugins[index] = pf(w, h)

		for _, redirOp := range h.Plugins[index].TypeOp() {
			if h.RedirectOp[redirOp] == nil {
				h.RedirectOp[redirOp] = make([]int, 0, 4)
			}
			h.RedirectOp[redirOp] = append(h.RedirectOp[redirOp], index)
		}

		for _, redirCmd := range h.Plugins[index].TypeCommand() {
			if h.RedirectCommands[redirCmd] == nil {
				h.RedirectCommands[redirCmd] = make([]int, 0, 4)
			}
			h.RedirectCommands[redirCmd] = append(h.RedirectCommands[redirCmd], index)
		}

		for _, redirEvent := range h.Plugins[index].TypeEvent() {
			if h.RedirectEvent[redirEvent] == nil {
				h.RedirectEvent[redirEvent] = make([]int, 0, 4)
			}
			h.RedirectEvent[redirEvent] = append(h.RedirectEvent[redirEvent], index)
		}
	}

	return h
}

// SendEvent sends an event with a tag.
func (l *Plugin) SendEvent(tag event.Tag, data event.Event) {
	l.eventsCustomNextMutex.Lock()
	defer l.eventsCustomNextMutex.Unlock()

	if l.eventsCustomNext.taggedEvents[tag] == nil {
		l.eventsCustomNext.taggedEvents[tag] = make([]event.Event, 0, 128)
	}

	l.eventsCustomNext.taggedEvents[tag] = append(l.eventsCustomNext.taggedEvents[tag], data)

	if !l.Invalidated.Load() {
		l.window.Invalidate()
		l.Invalidated.Store(true)
	}
}

// SendEventUntagged sends an event without a tag.
func (l *Plugin) SendEventUntagged(tag uint64, data event.Event) {
	l.eventsCustomNextMutex.Lock()
	defer l.eventsCustomNextMutex.Unlock()

	if l.eventsCustomNext.untaggedEvents[tag] == nil {
		l.eventsCustomNext.untaggedEvents[tag] = make([]event.Event, 0, 128)
	}

	l.eventsCustomNext.untaggedEvents[tag] = append(l.eventsCustomNext.untaggedEvents[tag], data)

	if !l.Invalidated.Load() {
		l.window.Invalidate()
		l.Invalidated.Store(true)
	}
}

/*
func (l *Plugin) Event(t ...event.Filter) (event.Event, bool) {
	if l == nil {
		return nil, false
	}

	if evt, ok := l.event(t...); ok {
		return evt, true
	}

	return l.OriginalSource.Event(t...)
}
*/

// Event returns the first event that matches the filters.
func (l *Plugin) Event(filters ...event.Filter) (event.Event, bool) {
	l.eventsCustomCurrentMutex.Lock()
	defer l.eventsCustomCurrentMutex.Unlock()

	for _, filter := range filters {
		switch f := filter.(type) {
		case Filter:
			tag := f.Tag()
			for _, evt := range l.eventsCustomCurrent.taggedEvents[tag] {
				if !f.Matches(evt) {
					continue
				}

				copy(l.eventsCustomCurrent.taggedEvents[tag], l.eventsCustomCurrent.taggedEvents[tag][1:])
				l.eventsCustomCurrent.taggedEvents[tag] = l.eventsCustomCurrent.taggedEvents[tag][:len(l.eventsCustomCurrent.taggedEvents[tag])-1]

				return evt, true
			}
		case UntaggedFilter:
			tag := f.Name()
			for _, evt := range l.eventsCustomCurrent.untaggedEvents[tag] {
				if !f.Matches(evt) {
					continue
				}

				copy(l.eventsCustomCurrent.untaggedEvents[tag], l.eventsCustomCurrent.untaggedEvents[tag][1:])
				l.eventsCustomCurrent.untaggedEvents[tag] = l.eventsCustomCurrent.untaggedEvents[tag][:len(l.eventsCustomCurrent.untaggedEvents[tag])-1]

				return evt, true
			}
		}
	}

	return nil, false
}

/*
func (l *Plugin) Execute(c input.Command) {
	if ok := l.execute(c); ok {
		return
	}
	l.OriginalSource.Execute(c)
}
*/

// Execute executes the command.
func (l *Plugin) Execute(c input.Command) bool {
	t := reflect.TypeOf(c)
	if _, ok := l.RedirectCommands[t]; !ok {
		return false
	}

	for _, index := range l.RedirectCommands[t] {
		l.Plugins[index].Execute(c)
	}

	return true
}

// Enabled returns if the plugin is enabled.
func (l *Plugin) Enabled() bool {
	return true
}

// Focused returns if the plugin is focused.
func (l *Plugin) Focused(tag event.Tag) bool {
	return l.OriginalSource.Focused(tag)
}

// Frame is the function called at the end of the frame.
func (l *Plugin) Frame(ops *op.Ops) {
	l.OriginalFrame(ops)

	for i := range l.visited {
		delete(l.visited, i)
	}

	if len(l.RedirectOp) > 0 {
		l.Op((*unsafeOps)(unsafe.Pointer(&ops.Internal)))
	}

	// Must be after processing ops
	for _, index := range l.RedirectEvent[reflect.TypeOf(EndFrameEvent{})] {
		l.Plugins[index].Event(EndFrameEvent{})
	}

	l.eventsCustomNextMutex.Lock()
	l.eventsCustomCurrentMutex.Lock()
	for v := range l.eventsCustomCurrent.taggedEvents {
		l.eventsCustomCurrent.taggedEvents[v] = l.eventsCustomCurrent.taggedEvents[v][:0]
	}
	for v := range l.eventsCustomCurrent.untaggedEvents {
		l.eventsCustomCurrent.untaggedEvents[v] = l.eventsCustomCurrent.untaggedEvents[v][:0]
	}
	l.eventsCustomCurrentMutex.Unlock()
	l.eventsCustomNextMutex.Unlock()
}

// unsafeOps is a copy of internal/ops/ops.go
type unsafeOps struct {
	version     uint32
	data        []byte
	refs        []interface{}
	stringRefs  []string
	nextStateID uint32
	multipOp    bool
	macroStack  [2]uint32
	stacks      [4][2]uint32
}

var (
	internalOps = op.Ops{}.Internal
	typeOps     = reflect.TypeOf(&internalOps)
)

// Op is the function that processes the ops.
func (l *Plugin) Op(o *unsafeOps) {
	if _, ok := l.visited[uintptr(unsafe.Pointer(o))]; ok {
		return
	}
	l.visited[uintptr(unsafe.Pointer(o))] = struct{}{}

	for i := range o.refs {
		if reflect.TypeOf(o.refs[i]) == typeOps {
			o2 := *(**unsafeOps)(unsafe.Add(unsafe.Pointer(&o.refs[i]), unsafe.Sizeof(uintptr(0))))
			l.Op(o2)
		} else {
			for _, index := range l.RedirectOp[reflect.TypeOf(o.refs[i])] {
				l.Plugins[index].Op(o.refs[i])
			}
		}
	}
}

// ProcessEventFromGio processes the event from Gio.
func (l *Plugin) ProcessEventFromGio(evt event.Event) event.Event {
	for _, index := range l.RedirectEvent[reflect.TypeOf(evt)] {
		l.Plugins[index].Event(evt)
	}

	switch e := evt.(type) {
	case app.FrameEvent:
		l.Invalidated.Store(false)

		l.eventsCustomNextMutex.Lock()
		l.eventsCustomCurrentMutex.Lock()
		l.eventsCustomNext, l.eventsCustomCurrent = l.eventsCustomCurrent, l.eventsCustomNext
		l.eventsCustomCurrentMutex.Unlock()
		l.eventsCustomNextMutex.Unlock()

		l.OriginalFrame = e.Frame
		e.Frame = l.Frame

		return e
	case app.ViewEvent:
		for _, index := range l.RedirectEvent[reflect.TypeOf(ViewEvent{})] {
			l.Plugins[index].Event(e)
		}
		return e
	default:
		return evt
	}
}
