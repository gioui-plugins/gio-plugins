package main

import (
	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"
	"github.com/gioui-plugins/gio-plugins/share/gioshare"
)

var (
	submit widget.Clickable
	isURL  widget.Bool
	list   widget.List

	title, desc, url widget.Editor
)

type mode int

const (
	modeText mode = iota
	modeLink
)

var (
	currentMode mode
	header      []layout.Widget
	footer      []layout.Widget
	inputs      map[mode][]layout.Widget
)

func main() {
	theme := material.NewTheme()
	list.Axis = layout.Vertical

	if header == nil {
		header = []layout.Widget{
			material.Switch(theme, &isURL, "Enable URL").Layout,
		}
	}

	if footer == nil {
		footer = []layout.Widget{
			material.Button(theme, &submit, "Submit").Layout,
		}
	}

	if inputs == nil {
		inputs = make(map[mode][]layout.Widget)
		inputs[modeText] = []layout.Widget{
			material.Editor(theme, &title, "Title").Layout,
			material.Editor(theme, &desc, "Desc").Layout,
		}
		inputs[modeLink] = []layout.Widget{
			material.Editor(theme, &url, "URL").Layout,
			material.Editor(theme, &title, "Title").Layout,
			material.Editor(theme, &desc, "Desc").Layout,
		}
	}

	url.SetText("https://google.com")
	title.SetText("Example Title")
	desc.SetText("Example Text")

	w := &app.Window{}
	ops := new(op.Ops)

	go func() {
		for {
			e := gioplugins.Event(w)

			switch e := e.(type) {
			case app.FrameEvent:
				gtx := app.NewContext(ops, e)

				if isURL.Update(gtx) {
					if currentMode == modeText {
						currentMode = modeLink
					} else {
						currentMode = modeText
					}
				}

				if submit.Clicked(gtx) {
					switch currentMode {
					case modeText:
						gtx.Execute(gioshare.TextCmd{Title: title.Text(), Text: desc.Text()})
					case modeLink:
						gtx.Execute(gioshare.WebsiteCmd{Title: title.Text(), Text: desc.Text(), Link: url.Text()})
					}
				}

				layout.UniformInset(unit.Dp(30)).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
					headerIndex, inputsIndex := 0, 0
					return material.List(theme, &list).Layout(gtx, len(inputs[currentMode])+len(header)+len(footer), func(gtx layout.Context, index int) layout.Dimensions {
						if len(header) > headerIndex {
							d := header[index](gtx)
							d.Size.Y += gtx.Dp(20)
							headerIndex++
							return d
						}
						index -= headerIndex
						if len(inputs[currentMode]) > inputsIndex {
							d := inputs[currentMode][index](gtx)
							d.Size.Y += gtx.Dp(20)
							inputsIndex++
							return d
						}
						index -= inputsIndex
						if len(footer) > index {
							return footer[index](gtx)
						}
						return layout.Dimensions{}
					})

				})

				e.Frame(gtx.Ops)
			}
		}
	}()

	app.Main()
}
