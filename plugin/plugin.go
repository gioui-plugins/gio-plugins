package plugin

import (
	"reflect"
	"sync"
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/op"
)

var handlers = new(sync.Map) // map[app.Window]handler

type Plugin struct {
	window *app.Window
	queue  event.Queue

	eventsCustom map[event.Tag][]event.Event
	eventsPool   []event.Event

	redirectEvent map[reflect.Type][]int
	redirectOp    map[reflect.Type][]int

	plugins     []Handler
	invalidated bool
}

func newHandler(w *app.Window) *Plugin {
	h := &Plugin{
		window:        w,
		eventsCustom:  make(map[event.Tag][]event.Event, 128),
		plugins:       make([]Handler, len(registeredPlugins)),
		redirectOp:    map[reflect.Type][]int{},
		redirectEvent: map[reflect.Type][]int{},
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
	if l.eventsCustom == nil {
		l.eventsCustom = make(map[event.Tag][]event.Event, 128)
	}
	if l.eventsCustom[tag] == nil {
		l.eventsCustom[tag] = make([]event.Event, 0, 128)
	}
	l.eventsCustom[tag] = append(l.eventsCustom[tag], data)

	if !l.invalidated {
		l.window.Invalidate()
	}
}

func (l *Plugin) Events(t event.Tag) []event.Event {
	evtsGio := l.queue.Events(t)
	evtsCustom, _ := l.eventsCustom[t]

	switch {
	case len(evtsGio) > 0 && len(evtsCustom) > 0:
		l.eventsPool = l.eventsPool[:0]

		l.eventsPool = append(l.eventsPool, evtsGio...)
		l.eventsPool = append(l.eventsPool, evtsCustom...)

		l.eventsCustom[t] = l.eventsCustom[t][:0]

		return l.eventsPool
	case len(evtsGio) > 0:
		return evtsGio
	case len(evtsCustom) > 0:
		l.eventsCustom[t] = l.eventsCustom[t][:0]
		return evtsCustom
	default:
		return nil
	}
}

type unsafeOps struct {
	version     int
	data        []byte
	refs        []interface{}
	nextStateID int
	multipOp    bool
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
	case system.FrameEvent:
		ref := *(**system.FrameEvent)(unsafe.Add(unsafe.Pointer(&evt), unsafe.Sizeof(uintptr(0))))
		h.invalidated = false

		q := ref.Queue
		h.queue = q
		ref.Queue = h

		f := ref.Frame
		ref.Frame = func(frame *op.Ops) {
			ops := (*unsafeOps)(unsafe.Pointer(&frame.Internal))
			f(frame)

			for _, r := range ops.refs {
				for _, index := range h.redirectOp[reflect.TypeOf(r)] {
					h.plugins[index].ListenOps(r)
				}
			}

			for _, index := range h.redirectEvent[reflect.TypeOf(EndFrameEvent{})] {
				h.plugins[index].ListenEvents(EndFrameEvent{})
			}
		}
	}
}

type EndFrameEvent struct{}

func (EndFrameEvent) ImplementsEvent() {}
