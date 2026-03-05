package main

import (
	"encoding/json"
	"image"
	"image/color"
	"log"
	"os"
	"strings"

	_ "unsafe"

	"gioui.org/font"
	"gioui.org/io/key"
	"github.com/gioui-plugins/gio-plugins/inapppay/gioinapppay"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(app.Size(unit.Dp(800), unit.Dp(700)), app.MinSize(unit.Dp(400), unit.Dp(400)))

		if err := loop(w); err != nil {
			panic(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops
	for {
		e := gioplugins.Hijack(w)

		switch e := e.(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			var submitted = false

			for {
				evt, ok := InputAction.Update(gtx)
				if !ok {
					break
				}

				switch evt.(type) {
				case widget.SubmitEvent:
					submitted = true
					break
				}
			}

			if ButtonAction.Clicked(gtx) || submitted {
				gioplugins.Execute(gtx, gioinapppay.ListProductsCmd{ProductIDs: []string{InputAction.Text()}})
				InputDetail.SetText("Querying...")
				submitted = false
			}

			if BuyAction.Clicked(gtx) {
				gioplugins.Execute(gtx, gioinapppay.PurchaseCmd{ProductID: InputAction.Text()})
				InputDetail.SetText("Buying...")
			}

			for {
				evt, ok := gioplugins.Event(gtx, gioinapppay.Filter{})
				if !ok {
					break
				}

				switch evt := evt.(type) {
				case gioinapppay.ProductDetailsEvent:
					v, _ := json.Marshal(evt.Products)
					InputDetail.SetText(string(v))
					log.Printf("Products: %s\n", string(v))
					break
				case gioinapppay.PaymentResultEvent:
					v, _ := json.Marshal(evt)
					InputDetail.SetText(string(v))
					log.Printf("Payment Result: %s\n", string(v))
					break

				case gioinapppay.ErrorEvent:
					InputDetail.SetText("Error: " + evt.Error.Error())
					log.Printf("Error: %s\n", evt.Error)
				}
			}

			render(gtx)
			e.Frame(gtx.Ops)
		}
	}
}

func render(gtx layout.Context) layout.Dimensions {
	return layout.Flex{Axis: layout.Vertical, Spacing: layout.SpaceSides}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return MarginDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return InputBackgroundDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return InputDesign.Layout(gtx, InputAction, "ProductID", "1000")
				})
			})
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return MarginDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return ButtonDesign.Layout(gtx, ButtonAction, "QUERY")
			})
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return MarginDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return InputBackgroundDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					return InputDesign.Layout(gtx, InputDetail, "", "")
				})
			})
		}),

		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return MarginDesign.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return ButtonDesign.Layout(gtx, BuyAction, "BUY")
			})
		}),
	)
}

// Actions
var (
	ButtonAction = &widget.Clickable{}
	BuyAction    = &widget.Clickable{}
	InputAction  = &widget.Editor{SingleLine: true, Submit: true, InputHint: key.HintURL}
	InputDetail  = &widget.Editor{SingleLine: true, Submit: true, InputHint: key.HintURL}
)

// Design
var (
	ButtonDesign          = &Button{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}, TextSize: unit.Sp(16), Background: color.NRGBA{R: 135, G: 156, B: 251, A: 255}, BorderRadius: unit.Dp(4), Modifier: strings.ToUpper, Inset: layout.Inset{Top: unit.Dp(10), Right: unit.Dp(12), Bottom: unit.Dp(10), Left: unit.Dp(12)}}
	InputDesign           = &Input{Font: font.Font{}, TextSize: unit.Sp(14), Color: color.NRGBA{R: 100, G: 130, B: 60, A: 255}, HintColor: color.NRGBA{R: 120, G: 120, B: 120, A: 255}}
	InputBackgroundDesign = &Background{Color: color.NRGBA{R: 234, G: 236, B: 231, A: 255}, Inset: layout.UniformInset(unit.Dp(13)), BorderRadius: unit.Dp(10)}

	MarginDesign = layout.Inset{Right: unit.Dp(30), Bottom: unit.Dp(6), Left: unit.Dp(30), Top: unit.Dp(6)}
)

var defaultMaterial = material.NewTheme()

type Input struct {
	Font      font.Font
	TextSize  unit.Sp
	Color     color.NRGBA
	HintColor color.NRGBA
	notSet    bool
}

func (i *Input) Layout(gtx layout.Context, editor *widget.Editor, hint string, value string) layout.Dimensions {
	e := material.Editor(defaultMaterial, editor, hint)
	e.TextSize = i.TextSize
	e.Color = i.Color
	e.Hint = hint
	e.HintColor = i.HintColor

	if value != "" && !i.notSet {
		editor.SetText(value)
		editor.MoveCaret(editor.Len(), editor.Len())
		i.notSet = true
	}

	return e.Layout(gtx)
}

type Button struct {
	Color        color.NRGBA
	Font         font.Font
	TextSize     unit.Sp
	Background   color.NRGBA
	BorderRadius unit.Dp
	Modifier     func(string) string
	Inset        layout.Inset
}

func (b *Button) Layout(gtx layout.Context, clickable *widget.Clickable, s string) layout.Dimensions {
	style := material.Button(defaultMaterial, clickable, s)
	style.Color = b.Color
	style.Font = b.Font
	style.TextSize = b.TextSize
	style.Background = b.Background
	style.CornerRadius = b.BorderRadius
	style.Inset = b.Inset

	if b.Modifier != nil {
		style.Text = b.Modifier(s)
	}

	return style.Layout(gtx)
}

type Background struct {
	Color        color.NRGBA
	BorderRadius unit.Dp
	Inset        layout.Inset
}

func (b *Background) Layout(gtx layout.Context, w layout.Widget) layout.Dimensions {

	macro := op.Record(gtx.Ops)
	dimensions := b.Inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		return w(gtx)
	})
	saved := macro.Stop()

	return layout.Stack{}.Layout(gtx,
		layout.Stacked(func(gtx layout.Context) layout.Dimensions {
			background := image.Rectangle{Max: image.Point{X: dimensions.Size.X, Y: dimensions.Size.Y}}

			rr := gtx.Dp(b.BorderRadius)
			stack := clip.RRect{Rect: background, NE: rr, NW: rr, SE: rr, SW: rr}.Op(gtx.Ops).Push(gtx.Ops)
			paint.Fill(gtx.Ops, b.Color)
			stack.Pop()

			return dimensions
		}),

		layout.Expanded(func(gtx layout.Context) layout.Dimensions {
			saved.Add(gtx.Ops)
			return dimensions
		}),
	)
}
