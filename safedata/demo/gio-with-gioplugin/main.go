package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"
	"github.com/gioui-plugins/gio-plugins/safedata"
	"github.com/gioui-plugins/gio-plugins/safedata/giosafedata"
	"image/color"
	"os"
)

// Actions
var (
	Save = &widget.Clickable{}
	Load = &widget.Clickable{}

	TextKey  = &widget.Editor{SingleLine: true, Submit: true}
	TextData = &widget.Editor{SingleLine: true, Submit: true}
)

func main() {
	window := &app.Window{}
	theme := material.NewTheme()

	ops := new(op.Ops)
	go func() {
		for {
			evt := gioplugins.Event(window)

			switch evt := evt.(type) {
			case app.FrameEvent:
				gtx := app.NewContext(ops, evt)

				paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)

				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Label(theme, unit.Sp(14), "Key").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Editor(theme, TextKey, "Key").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if Load.Clicked(gtx) {
							gtx.Execute(giosafedata.ReadSecretCmd{Tag: Load, Identifier: TextKey.Text()})
							TextKey.SetText("")
							TextData.SetText("")
						}
						for {
							evt, ok := gtx.Event(giosafedata.Filter{Source: Load})
							if !ok {
								break
							}

							if evt, ok := evt.(giosafedata.SecretsEvent); ok {
								if len(evt.Secrets) == 0 {
									TextData.SetText("Not found")
								} else {
									TextData.SetText(string(evt.Secrets[0].Data))
								}
							}

							if evt, ok := evt.(giosafedata.ErrorEvent); ok {
								fmt.Println("Error:", evt.Error)
							}
						}
						return material.Button(theme, Load, "Load").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Label(theme, unit.Sp(14), "Data").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return material.Editor(theme, TextData, "Data").Layout(gtx)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						if Save.Clicked(gtx) {
							gtx.Execute(giosafedata.WriteSecretCmd{
								Secret: safedata.Secret{
									Identifier: TextKey.Text(),
									Data:       []byte(TextData.Text()),
								},
							})

							TextKey.SetText("")
							TextData.SetText("")
						}
						return material.Button(theme, Save, "Save").Layout(gtx)
					}),
				)

				gtx.Execute(op.InvalidateCmd{})
				evt.Frame(ops)
			case app.DestroyEvent:
				os.Exit(0)
			}
		}
	}()

	app.Main()
}
