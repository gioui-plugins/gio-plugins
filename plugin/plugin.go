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

type Plugin struct {
	window *app.Window

	eventsCustomNextMutex    sync.Mutex
	eventsCustomCurrentMutex sync.Mutex

	// double buffered events
	eventsCustomNext    map[event.Tag][]event.Event
	eventsCustomCurrent map[event.Tag][]event.Event

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

func NewPlugin(w *app.Window) *Plugin {
	h := &Plugin{
		window:           w,
		Plugins:          make([]Handler, len(registeredPlugins)),
		visited:          make(map[uintptr]struct{}, 128),
		RedirectOp:       make(map[reflect.Type][]int, 128),
		RedirectCommands: make(map[reflect.Type][]int, 128),
		RedirectEvent:    make(map[reflect.Type][]int, 128),
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

func (l *Plugin) SendEvent(tag event.Tag, data event.Event) {
	l.eventsCustomNextMutex.Lock()
	defer l.eventsCustomNextMutex.Unlock()

	if l.eventsCustomNext == nil {
		l.eventsCustomNext = make(map[event.Tag][]event.Event, 128)
	}

	if l.eventsCustomNext[tag] == nil {
		l.eventsCustomNext[tag] = make([]event.Event, 0, 128)
	}

	l.eventsCustomNext[tag] = append(l.eventsCustomNext[tag], data)

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

func (l *Plugin) Event(filters ...event.Filter) (event.Event, bool) {
	l.eventsCustomCurrentMutex.Lock()
	defer l.eventsCustomCurrentMutex.Unlock()

	for _, filter := range filters {
		f, ok := filter.(Filter)
		if !ok {
			continue
		}

		tag := f.Tag()
		for _, evt := range l.eventsCustomCurrent[tag] {
			if !f.Matches(evt) {
				continue
			}

			copy(l.eventsCustomCurrent[tag], l.eventsCustomCurrent[tag][1:])
			l.eventsCustomCurrent[tag] = l.eventsCustomCurrent[tag][:len(l.eventsCustomCurrent[tag])-1]

			return evt, true
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

func (l *Plugin) Enabled() bool {
	return true
}

func (l *Plugin) Focused(tag event.Tag) bool {
	return l.OriginalSource.Focused(tag)
}

func (l *Plugin) Frame(ops *op.Ops) {
	l.OriginalFrame(ops)

	for _, index := range l.RedirectEvent[reflect.TypeOf(EndFrameEvent{})] {
		l.Plugins[index].Event(EndFrameEvent{})
	}

	for i := range l.visited {
		delete(l.visited, i)
	}

	if len(l.RedirectOp) > 0 {
		l.Op((*unsafeOps)(unsafe.Pointer(&ops.Internal)))
	}

	l.eventsCustomNextMutex.Lock()
	l.eventsCustomCurrentMutex.Lock()
	for v := range l.eventsCustomCurrent {
		l.eventsCustomCurrent[v] = l.eventsCustomCurrent[v][:0]
	}
	l.eventsCustomNextMutex.Unlock()
	l.eventsCustomCurrentMutex.Unlock()
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
		l.eventsCustomNextMutex.Unlock()
		l.eventsCustomCurrentMutex.Unlock()

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
