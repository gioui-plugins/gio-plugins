package plugin

import (
	"image"
	"reflect"
	"testing"
	"time"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/op"
	"gioui.org/unit"
)

type testingOp struct {
	data string
}

var testingOpPool = NewOpPool[testingOp]()

func (t testingOp) Add(ops *op.Ops) {
	testingOpPool.WriteOp(ops, t)
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
	return []reflect.Type{reflect.TypeOf(app.FrameEvent{})}
}

func (t *testingPlugin) ListenOps(op interface{}) {
	if op, ok := op.(*testingOp); ok {
		testingOpPool.Release(op)
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
	window := new(app.Window)

	frameCalled := false

	evt := event.Event(app.FrameEvent{
		Now:    time.Now(),
		Metric: unit.Metric{},
		Size:   image.Point{},
		Insets: app.Insets{},
		Frame: func(frame *op.Ops) {
			frameCalled = true
		},
		Source: input.Source{},
	})

	h := &testingPlugin{}
	Register(func(w *app.Window, handler *Plugin) Handler {
		h.w = w
		h.p = handler
		return h
	})
	Install(window, evt)

	ops := new(op.Ops)

	gtx := app.NewContext(ops, evt.(app.FrameEvent))
	testingOp{data: "test"}.Add(ops)
	evt.(app.FrameEvent).Frame(gtx.Ops)

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
