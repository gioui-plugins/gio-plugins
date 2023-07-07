package main

import (
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/deeplink"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"image"
	"image/color"
	"os"
)

func main() {
	w := app.NewWindow()
	lt := text.NewShaper(gofont.Collection())

	var texts = []string{"Schemes:"}

	ops := new(op.Ops)

	go func() {
		d := deeplink.NewDeepLink()
		for e := range d.Events() {
			switch e := e.(type) {
			case deeplink.Linked:
				texts = append(texts, e.URL.String())
			}
			w.Invalidate()
		}
	}()

	go func() {
		for e := range w.Events() {
			plugin.Install(w, e)

			switch e := e.(type) {
			case system.DestroyEvent:
				os.Exit(0)
				return
			case system.FrameEvent:
				gtx := layout.NewContext(ops, e)
				_ = gtx

				s := clip.RRect{Rect: image.Rectangle{Max: gtx.Constraints.Max}}.Push(gtx.Ops)
				paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, A: 255}}.Add(gtx.Ops)
				paint.PaintOp{}.Add(gtx.Ops)
				s.Pop()

				offset := image.Pt(0, 0)

				p := op.Record(gtx.Ops)
				paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
				painter := p.Stop()

				for _, txt := range texts {
					gtx.Constraints.Min.Y = 0
					gtx.Constraints.Min.X = 0

					o := op.Offset(offset).Push(gtx.Ops)
					dims := widget.Label{}.Layout(gtx, lt, font.Font{}, 16, txt, painter)
					o.Pop()
					offset.Y += dims.Size.Y
				}

				e.Frame(ops)
			}
		}
	}()

	app.Main()
}
