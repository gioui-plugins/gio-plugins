package plugin

import (
	"gioui.org/io/input"
	"reflect"
	"sync"
	"sync/atomic"
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
)

var handlers = new(sync.Map) // map[app.Window]handler

type Plugin struct {
	window *app.Window
	queue  input.Source

	eventsMutex  sync.Mutex
	eventsCustom map[event.Tag][]event.Event
	eventsPool   []event.Event

	redirectEvent map[reflect.Type][]int
	redirectOp    map[reflect.Type][]int

	visited map[uintptr]struct{}

	plugins     []Handler
	invalidated atomic.Bool
}

func newHandler(w *app.Window) *Plugin {
	h := &Plugin{
		window:        w,
		eventsCustom:  make(map[event.Tag][]event.Event, 128),
		plugins:       make([]Handler, len(registeredPlugins)),
		redirectOp:    make(map[reflect.Type][]int, 128),
		redirectEvent: make(map[reflect.Type][]int, 128),
		visited:       make(map[uintptr]struct{}, 2048),
	}
	for index, pf := range registeredPlugins {
		h.plugins[index] = pf(w, h)

		for _, redirOp := range h.plugins[index].TypeOp() {
			if h.redirectOp[redirOp] == nil {
				h.redirectOp[redirOp] = make([]int, 0, 4)
			}
			h.redirectOp[redirOp] = append(h.redirectOp[redirOp], index)
		}

		for _, redirEvent := range h.plugins[index].TypeEvent() {
			if h.redirectEvent[redirEvent] == nil {
				h.redirectEvent[redirEvent] = make([]int, 0, 4)
			}
			h.redirectEvent[redirEvent] = append(h.redirectEvent[redirEvent], index)
		}
	}

	return h
}

func (l *Plugin) SendEvent(tag event.Tag, data event.Event) {
	l.eventsMutex.Lock()
	defer l.eventsMutex.Unlock()

	if l.eventsCustom == nil {
		l.eventsCustom = make(map[event.Tag][]event.Event, 128)
	}
	if l.eventsCustom[tag] == nil {
		l.eventsCustom[tag] = make([]event.Event, 0, 128)
	}
	l.eventsCustom[tag] = append(l.eventsCustom[tag], data)

	if l.invalidated.Load() {
		l.window.Invalidate()
	}
}

func (l *Plugin) Events(t event.Tag) {
	l.eventsMutex.Lock()
	defer l.eventsMutex.Unlock()

	for {
		_, ok := l.queue.Event(EndFrameEvent{})
		if !ok {
			break
		}
		//evtsCustom, _ := l.eventsCustom[t]
		//
		//switch {
		//case len(evtsGio) > 0 && len(evtsCustom) > 0:
		//	l.eventsPool = l.eventsPool[:0]
		//
		//	l.eventsPool = append(l.eventsPool, evtsGio...)
		//	l.eventsPool = append(l.eventsPool, evtsCustom...)
		//
		//	l.eventsCustom[t] = l.eventsCustom[t][:0]
		//
		//	return l.eventsPool
		//case len(evtsGio) > 0:
		//	return evtsGio
		//case len(evtsCustom) > 0:
		//	l.eventsCustom[t] = l.eventsCustom[t][:0]
		//	return evtsCustom
		//default:
		//	return nil
		//}
	}

}

type unsafeOps struct {
	version     int
	data        []byte
	refs        []interface{}
	nextStateID int
	multipOp    bool
}

var (
	internalOps = op.Ops{}.Internal
	typeOps     = reflect.TypeOf(&internalOps)
)

func (l *Plugin) processFrameEvent(o *unsafeOps) {
	if _, ok := l.visited[uintptr(unsafe.Pointer(o))]; ok {
		return
	}
	l.visited[uintptr(unsafe.Pointer(o))] = struct{}{}

	for i := range o.refs {
		if reflect.TypeOf(o.refs[i]) == typeOps {
			o2 := *(**unsafeOps)(unsafe.Add(unsafe.Pointer(&o.refs[i]), unsafe.Sizeof(uintptr(0))))
			l.processFrameEvent(o2)
		} else {
			for _, index := range l.redirectOp[reflect.TypeOf(o.refs[i])] {
				l.plugins[index].ListenOps(o.refs[i])
			}
		}
	}
}

func Install(w *app.Window, evt event.Event) {
	var h *Plugin
	li, ok := handlers.Load(w)
	if !ok {
		h = newHandler(w)
		handlers.Store(w, h)
	} else {
		h = li.(*Plugin)
	}

	for _, index := range h.redirectEvent[reflect.TypeOf(evt)] {
		h.plugins[index].ListenEvents(evt)
	}

	switch evt.(type) {
	case app.FrameEvent:
		ref := *(**app.FrameEvent)(unsafe.Add(unsafe.Pointer(&evt), unsafe.Sizeof(uintptr(0))))
		h.invalidated.Store(false)

		//q := ref.Queue
		//h.queue = q
		//ref.Queue = h

		f := ref.Frame
		ref.Frame = func(frame *op.Ops) {
			f(frame)

			h.processFrameEvent((*unsafeOps)(unsafe.Pointer(&frame.Internal)))

			for _, index := range h.redirectEvent[reflect.TypeOf(EndFrameEvent{})] {
				h.plugins[index].ListenEvents(evt)
			}

			for i := range h.visited {
				delete(h.visited, i)
			}
		}
	}
}

type EndFrameEvent struct{}

func (e EndFrameEvent) ImplementsEvent() {
}

func (e EndFrameEvent) ImplementsFilter() {
}
