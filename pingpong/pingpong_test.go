package pingpong

import (
	"fmt"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"
	"os"
	"testing"

	"gioui.org/app"
	"gioui.org/op"
)

func TestMain(m *testing.M) {
	w := &app.Window{}
	done := make(chan struct{})

	go func() {
		ops := new(op.Ops)
		first := true

		for {
			evt := gioplugins.Hijack(w)

			switch evt := evt.(type) {
			case app.FrameEvent:
				gtx := app.NewContext(ops, evt)

				// By Op
				PingOp{Tag: w, Text: "Ping"}.Add(gtx.Ops)

				// By Command
				gioplugins.Execute(gtx, PingCmd{Tag: w, Text: "Ping"})

				founds := 0
				for {
					evt, ok := gioplugins.Event(gtx, Filter{Target: w})
					if !ok {
						break
					}

					if _, ok := evt.(PongEvent); ok {
						founds++
					}
				}

				if !first {
					if founds != 2 {
						panic("failed to receive pong")
					}

					if founds == 2 {
						done <- struct{}{}
						return
					}
				}

				if first {
					first = false
				}

				evt.Frame(gtx.Ops)
			}
		}
	}()

	go func() {
		<-done
		os.Exit(0)
	}()

	fmt.Println("TestMain")
	app.Main()
}
