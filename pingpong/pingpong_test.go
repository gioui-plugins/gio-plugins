package pingpong

import (
	"os"
	"testing"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

func TestMain(m *testing.M) {
	w := app.NewWindow()
	done := make(chan struct{})

	go func() {
		ops := new(op.Ops)
		first := true
		found := false

		for evt := range w.Events() {
			plugin.Install(w, evt)

			switch evt := evt.(type) {
			case system.FrameEvent:
				gtx := layout.NewContext(ops, evt)
				PingOp{Tag: &w, Text: "Ping"}.Add(gtx.Ops)

				if !first {
					first = true
				}

				for _, e := range gtx.Events(&w) {
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
