package main

import (
	"gioui.org/font"
	"gioui.org/op/paint"
	"image/color"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/safedata"
	"github.com/gioui-plugins/gio-plugins/safedata/giosafedata"
)

func main() {
	config := safedata.Config{App: "MyApp"}
	sh := safedata.NewSafeData(config)

	window := &app.Window{}

	shaper := text.NewShaper(text.WithCollection(gofont.Collection()))

	mutex := new(sync.Mutex)
	txt := ""

	ops := new(op.Ops)
	ready := make(chan int, 1)
	go func() {
		for {
			evt := window.Event()

			switch evt := evt.(type) {
			case app.ViewEvent:
				config = giosafedata.NewConfigFromViewEvent(window, evt, config.App)
				sh.Configure(config)
				ready <- 1
			case app.FrameEvent:
				mutex.Lock()
				gtx := app.NewContext(ops, evt)

				paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
				widget.Label{}.Layout(gtx, shaper, font.Font{}, 12, txt, op.CallOp{})

				gtx.Execute(op.InvalidateCmd{})
				evt.Frame(ops)
				mutex.Unlock()
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	go func() {
		<-ready
		m := make(chan struct{})
		go func() {
			// Add a secret
			err := sh.Set(safedata.Secret{
				Identifier:  "my-secret4",
				Description: "some secret",
				Data:        []byte("my-secret-data"),
			})

			if err != nil {
				mutex.Lock()
				txt = "ERR ON ADD->" + err.Error()
				mutex.Unlock()
			}
			m <- struct{}{}
		}()
		<-m

		go func() {
			// Add a secret
			x, err := sh.Get("my-secret4")
			mutex.Lock()
			if err != nil {
				txt += "\n" + string("ERR ON GET -> "+err.Error())
			} else {
				txt = string(x.Data)
			}
			mutex.Unlock()

			m <- struct{}{}
		}()

		<-m
	}()
	app.Main()
}
