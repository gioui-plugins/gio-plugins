package main

import (
	"fmt"
	"os"
	"sync"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
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

	window := app.NewWindow()

	shaper := text.NewCache(gofont.Collection())

	mutex := new(sync.Mutex)
	textx := ""

	ops := new(op.Ops)
	ready := make(chan int, 1)
	go func() {
		for evt := range window.Events() {
			switch evt := evt.(type) {
			case app.ViewEvent:
				config = giosafedata.NewConfigFromViewEvent(window, evt, config.App)
				sh.Configure(config)
				ready <- 1
			case system.FrameEvent:
				mutex.Lock()
				gtx := layout.NewContext(ops, evt)

				widget.Label{}.Layout(gtx, shaper, text.Font{}, 12, textx)

				op.InvalidateOp{}.Add(gtx.Ops)
				evt.Frame(ops)
				mutex.Unlock()
			case system.DestroyEvent:
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
				Description: "aaaa",
				Data:        []byte("my-secret-data"),
			})

			if err != nil {
				mutex.Lock()
				textx = string("ERR ON ADD->" + err.Error())
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
				textx += "\n" + string("ERR ON GET->"+err.Error())
			} else {
				textx = string(x.Data)
			}
			mutex.Unlock()

			m <- struct{}{}
		}()

		<-m
	}()
	app.Main()
}
