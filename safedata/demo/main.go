package main

import (
	"fmt"
	"image/color"
	"os"
	"sync"

	"gioui.org/font"
	"gioui.org/op/paint"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/safedata"
	"github.com/gioui-plugins/gio-plugins/safedata/giosafedata"
)

func main() {
	config := safedata.Config{
		App: "MyApp",
	}
	sh := safedata.NewSafeData(config)

	window := new(app.Window)

	shaper := text.NewShaper() // todo set fount

	mutex := new(sync.Mutex)
	txt := ""

	ops := new(op.Ops)
	ready := make(chan int, 1)
	go func() {
		for {
			switch evt := window.Event().(type) {
			case app.ViewEvent:
				config = giosafedata.NewConfigFromViewEvent(window, evt, config.App)
				sh.Configure(config)
				ready <- 1
			case app.FrameEvent:
				mutex.Lock()
				gtx := app.NewContext(ops, evt)

				paint.ColorOp{Color: color.NRGBA{0, 0, 0, 255}}.Add(gtx.Ops)
				widget.Label{}.Layout(gtx, shaper, font.Font{}, 12, txt, op.CallOp{})

				// op.InvalidateOp{}.Add(gtx.Ops)//todo
				op.InvalidateCmd{}.ImplementsCommand()
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
			fmt.Print("my-secret4")
			// Add a secret
			err := sh.Set(safedata.Secret{
				Identifier:  "my-secret4",
				Description: "some secret",
				Data:        []byte("my-secret-data"),
			})
			if err != nil {
				mutex.Lock()
				txt = string("ERR ON ADD->" + err.Error())
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
