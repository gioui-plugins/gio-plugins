package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/safehouse"
	"sync"
)

func main() {
	sh := safehouse.NewSafeHouse()

	window := app.NewWindow()

	shaper := text.NewCache(gofont.Collection())

	mutex := new(sync.Mutex)
	textx := ""

	ops := new(op.Ops)
	go func() {
		for evt := range window.Events() {
			switch evt := evt.(type) {
			case system.FrameEvent:
				mutex.Lock()
				gtx := layout.NewContext(ops, evt)

				widget.Label{}.Layout(gtx, shaper, text.Font{}, 12, textx)

				op.InvalidateOp{}.Add(gtx.Ops)
				evt.Frame(ops)
				mutex.Unlock()
			}
		}
	}()

	go func() {
		m := make(chan struct{})
		go func() {
			// Add a secret
			err := sh.Set(safehouse.Secret{
				Identifier: "my-secret3",
				Data:       []byte("my-secret-data"),
			})

			fmt.Println(err)
			m <- struct{}{}
		}()
		<-m

		go func() {
			// Add a secret
			x, _ := sh.Get("my-secret3")

			mutex.Lock()
			textx = string(x)
			mutex.Unlock()
			m <- struct{}{}
		}()

		<-m
	}()
	app.Main()

}
