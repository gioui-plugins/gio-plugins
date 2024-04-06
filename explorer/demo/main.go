package main

import (
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"

	"gioui.org/io/event"

	"gioui.org/font"

	"github.com/gioui-plugins/gio-plugins/explorer/gioexplorer"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	_ "golang.org/x/image/webp"

	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
	"github.com/gioui-plugins/gio-plugins/plugin"
)

var (
	_Sharper     = text.NewShaper()
	_ImageResult = make(chan img)
)

type img struct {
	widget widget.Image
	image  image.Image
}

func main() {
	w := new(app.Window)
	w.Option(app.Size(500, 500))
	ops := new(op.Ops)
	p := new(Page)
	p.tag = new(int)

	go func() {
		for {
			select {
			case img := <-_ImageResult:
				p.image = img.widget
				p.raw = img.image
				p.loading = false
				w.Invalidate()
			}
		}
	}()

	go func() {
		for {
			switch e := w.Event().(type) {
			case app.FrameEvent:
				plugin.Install(w, e)
				gtx := app.NewContext(ops, e)
				p.Layout(gtx)
				e.Frame(ops)
			}
		}
	}()

	app.Main()
}

type Page struct {
	uploadClickable widget.Clickable
	saveClickable   widget.Clickable

	tag     *int
	image   widget.Image
	raw     image.Image
	loading bool
	error   string
	cancel  bool
}

var _FileTypes = []mimetype.MimeType{
	{Extension: "png", Type: "image", Subtype: "png"},
	{Extension: "jpg", Type: "image", Subtype: "jpeg"},
	{Extension: "jpeg", Type: "image", Subtype: "jpeg"},
	{Extension: "gif", Type: "image", Subtype: "gif"},
	{Extension: "webp", Type: "image", Subtype: "webp"},
}

func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if p.uploadClickable.Clicked(gtx) {
		gioexplorer.OpenFileOp{Tag: p.tag, Mimetype: _FileTypes}.Add(gtx.Ops)
		p.error = ""
		p.cancel = false
	}

	if p.saveClickable.Clicked(gtx) {
		gioexplorer.SaveFileOp{Tag: p.tag, Mimetype: _FileTypes[0], Filename: "image.png"}.Add(gtx.Ops)
		p.error = ""
		p.cancel = false
	}
	events := []event.Filter{
		gioexplorer.SaveFileEvent{},
		gioexplorer.OpenFileEvent{},
		gioexplorer.ErrorEvent{},
		gioexplorer.CancelEvent{},
	}

	for {
		if p.loading { //?
			continue
		}
		evt, b := gtx.Event(events...)
		if !b {
			break
		}
		switch evt := evt.(type) {
		case gioexplorer.SaveFileEvent:
			go func() {
				defer evt.File.Close()

				err := png.Encode(evt.File, p.raw)
				if err != nil {
					return
				}
			}()
		case gioexplorer.OpenFileEvent:
			p.loading = true
			go func() {
				defer evt.File.Close()

				i, _, err := image.Decode(evt.File)
				if err != nil {
					p.loading = false
					return
				}
				_ImageResult <- img{widget: widget.Image{Fit: widget.Contain, Position: layout.Center, Src: paint.NewImageOp(i)}, image: i}
			}()
		case gioexplorer.ErrorEvent:
			p.error = evt.Error()
		case gioexplorer.CancelEvent:
			p.cancel = true
		}

	}

	return layout.Flex{Axis: layout.Vertical}.Layout(gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			if p.error != "" {
				return widget.Label{}.Layout(gtx, _Sharper, font.Font{}, 16, p.error, op.CallOp{})
			}
			if p.cancel {
				return widget.Label{}.Layout(gtx, _Sharper, font.Font{}, 16, "Canceled", op.CallOp{})
			}
			return layout.Dimensions{}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return Button{Clickable: &p.uploadClickable, text: "Open", color: color.NRGBA{R: 0, G: 0, B: 255, A: 255}}.Layout(gtx)
			})
		}),
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return p.image.Layout(gtx)
			})
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(20).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
				return Button{Clickable: &p.saveClickable, text: "Save", color: color.NRGBA{R: 38, G: 38, B: 38, A: 255}}.Layout(gtx)
			})
		}),
	)
}

type Button struct {
	*widget.Clickable
	text     string
	color    color.NRGBA
	disabled bool
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	macro := op.Record(gtx.Ops)
	var labelDims layout.Dimensions
	{
		gtx := gtx
		gtx.Constraints.Min = image.Point{}
		paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
		labelDims = widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, _Sharper, font.Font{}, 14, b.text, op.CallOp{})
	}
	call := macro.Stop()

	gtx.Constraints.Max.Y = gtx.Dp(20) + labelDims.Size.Y

	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		defer clip.Rect{Max: gtx.Constraints.Max}.Push(gtx.Ops).Pop()
		paint.ColorOp{Color: b.color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		defer op.Offset(image.Pt((gtx.Constraints.Max.X-labelDims.Size.X)/2, gtx.Dp(10))).Push(gtx.Ops).Pop()
		call.Add(gtx.Ops)

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}
