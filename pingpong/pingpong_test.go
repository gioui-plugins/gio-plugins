package pingpong

import (
	"gioui.org/io/event"
	"os"
	"testing"

	"gioui.org/app"

	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

func TestMain(m *testing.M) {
	w := new(app.Window)
	done := make(chan struct{})

	go func() {
		ops := new(op.Ops)
		first := true
		found := false

		filter := []event.Filter{PongEvent{}}

		for {
			switch evt := w.Event().(type) {
			case app.FrameEvent:
				plugin.Install(w, evt)
				gtx := app.NewContext(ops, evt)
				PingOp{Tag: &w, Text: "Ping"}.Add(gtx.Ops)

				if !first {
					first = true
				}

				for {
					e, b := gtx.Event(filter...)
					if !b {
						break
					}
					if _, ok := e.(PongEvent); ok {
						found = true
					}
				}

				if !first && !found {
					panic("failed to receive pong")
				} else {
					done <- struct{}{}
				}

				evt.Frame(gtx.Ops)
			}
		}
	}()

	go func() {
		<-done
		os.Exit(0)
	}()
	app.Main()
}
