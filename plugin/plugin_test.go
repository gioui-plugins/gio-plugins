package plugin

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"image"
	"reflect"
	"testing"
	"time"
)

type testingOp struct {
	data string
}

type testingEvent struct {
	data string
}

func (t testingEvent) ImplementsEvent() {}

type testingPlugin struct {
	w *app.Window
	p *Plugin

	ackListenEvent bool
	ackListenOps   bool
}

func (t *testingPlugin) TypeOp() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(&testingOp{})}
}

func (t *testingPlugin) TypeEvent() []reflect.Type {
	return []reflect.Type{reflect.TypeOf(system.FrameEvent{})}
}

func (t *testingPlugin) ListenOps(op interface{}) {
	if op, ok := op.(*testingOp); ok {
		t.ackListenOps = true
	}
}

func (t *testingPlugin) ListenEvents(evt event.Event) {
	if _, ok := evt.(app.FrameEvent); ok {
		t.ackListenEvent = true
		t.p.SendEvent(event.Tag(0), testingEvent{data: "test"})
	}
}

func TestInstall(t *testing.T) {
	window := app.NewWindow()

	frameCalled := false

	evt := event.Event(system.FrameEvent{
		Now:    time.Now(),
		Metric: unit.Metric{},
		Size:   image.Point{},
		Insets: system.Insets{},
		Frame: func(frame *op.Ops) {
			frameCalled = true
		},
		Queue: &router.Router{},
	})

	h := &testingPlugin{}
	Register(func(w *app.Window, handler *Plugin) Handler {
		h.w = w
		h.p = handler
		return h
	})
	Install(window, evt)

	ops := new(op.Ops)

	gtx := layout.NewContext(ops, evt.(system.FrameEvent))
	testingOp{data: "test"}.Add(ops)
	evt.(system.FrameEvent).Frame(gtx.Ops)

	if !frameCalled {
		t.Error("frame not called")
	}

	if !h.ackListenEvent {
		t.Error("ListenEvent not called")
	}
	if !h.ackListenOps {
		t.Error("ListenOps not called")
	}
}
