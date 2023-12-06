package authlayout

import (
	"gioui.org/font"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/inkeliz/giosvg"
	"image"
	"image/color"
	"math"
	"time"
)

type Format int

const (
	FormatRounded Format = iota
	FormatRectangle
	FormatPill
)

// ButtonStyle is the style for Google buttons.
//
// Notice that you may violate some branding guidelines, it's
// safer to use DefaultGoogleButtonStyle or DefaultAppleButtonStyle.
type ButtonStyle struct {
	Text                string
	TextSize            unit.Dp
	TextFont            font.Font
	TextShaper          *text.Shaper
	TextColor           color.NRGBA
	TextAlignment       layout.Alignment
	IconAlignment       layout.Alignment
	IconColor           color.NRGBA
	BackgroundColor     color.NRGBA
	BackgroundIconColor color.NRGBA
	BorderColor         color.NRGBA
	BorderThickness     unit.Dp
	Format              Format
}

func (b ButtonStyle) label(gtx layout.Context, text string) (op.CallOp, layout.Dimensions) {
	gtx.Constraints.Min.X = 0
	gtx.Constraints.Min.Y = 0

	r := op.Record(gtx.Ops)
	dims := widget.Label{}.Layout(gtx, b.TextShaper, b.TextFont, gtx.Metric.DpToSp(b.TextSize), text, toLabelColor(gtx, b.TextColor))
	return r.Stop(), dims
}

func toLabelColor(gtx layout.Context, c color.NRGBA) op.CallOp {
	r := op.Record(gtx.Ops)
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	return r.Stop()
}

func (b ButtonStyle) layoutText(gtx layout.Context, icon *giosvg.Icon, pointer *Pointer, text string, logoSize int, logoPadding int) layout.Dimensions {
	label, labelDims := b.label(gtx, text)

	minHeight := int(math.Round(float64(labelDims.Size.Y) * 233 / 100))

	inset := layout.Inset{
		Top:    gtx.Metric.PxToDp((minHeight - labelDims.Size.Y) / 2),
		Bottom: gtx.Metric.PxToDp((minHeight - labelDims.Size.Y) / 2),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}

	if logoSize == 0 {
		logoSize = labelDims.Size.Y
	}

	avalSize := gtx.Constraints.Max.X - (labelDims.Size.X + gtx.Dp(inset.Left) + gtx.Dp(inset.Right) + logoSize)
	if avalSize < 0 {
		// The label is too long, we need to render the icon-only button.
	}

	if avalSize > (logoPadding*2) && b.TextAlignment != layout.Start && b.IconAlignment != layout.Middle {
		logoPadding = 0
	}

	main := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		d := layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, labelDims.Size.Y)}

		{
			// Logo
			var off op.TransformStack
			switch b.IconAlignment {
			case layout.Start:
				off = op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
			case layout.Middle:
				off = op.Offset(image.Pt((gtx.Constraints.Max.X-logoSize-logoPadding-labelDims.Size.X)/2, 0)).Push(gtx.Ops)
			case layout.End:
				off = op.Offset(image.Pt(gtx.Constraints.Max.X-logoSize, 0)).Push(gtx.Ops)
			}

			// Logo Background
			padding := gtx.Dp(6)
			offBackground := op.Offset(image.Pt(-padding/2, -padding/2)).Push(gtx.Ops)
			background := clip.UniformRRect(image.Rectangle{Max: image.Pt(logoSize+padding, logoSize+padding)}, (logoSize+padding)/2).Push(gtx.Ops)
			paint.Fill(gtx.Ops, b.BackgroundIconColor)
			background.Pop()
			offBackground.Pop()

			gtx := gtx
			gtx.Constraints.Min = image.Point{}
			gtx.Constraints.Max.X, gtx.Constraints.Max.Y = logoSize, logoSize
			if b.IconColor.A != 0 {
				paint.ColorOp{Color: b.IconColor}.Add(gtx.Ops)
			}

			iconR := op.Record(gtx.Ops)
			dimsIcon := icon.Layout(gtx)
			iconOp := iconR.Stop()

			iconOff := op.Offset(image.Pt((logoSize-dimsIcon.Size.X)/2, (logoSize-dimsIcon.Size.Y)/2)).Push(gtx.Ops)
			iconOp.Add(gtx.Ops)
			iconOff.Pop()

			off.Pop()
		}

		{
			// Text
			gtx := gtx
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - logoSize - logoPadding

			if b.TextAlignment != layout.End {
				defer op.Offset(image.Pt(logoSize+logoPadding, 0)).Push(gtx.Ops).Pop()
			}

			switch b.TextAlignment {
			case layout.Start:
				label.Add(gtx.Ops)
			case layout.Middle:
				defer op.Offset(image.Pt((gtx.Constraints.Max.X-labelDims.Size.X)/2, 0)).Push(gtx.Ops).Pop()
				label.Add(gtx.Ops)
			case layout.End:
				defer op.Offset(image.Pt(gtx.Constraints.Max.X-labelDims.Size.X, 0)).Push(gtx.Ops).Pop()
				label.Add(gtx.Ops)
			}
		}

		return d
	})
	call := main.Stop()

	borderSize := gtx.Dp(b.BorderThickness)
	switch b.Format {
	case FormatRounded:
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, gtx.Dp(4)).Push(gtx.Ops).Pop()
	case FormatRectangle:
		defer clip.Rect{Max: dims.Size}.Push(gtx.Ops).Pop()
	case FormatPill:
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size}, dims.Size.Y/2).Push(gtx.Ops).Pop()
	}
	paint.ColorOp{Color: b.BorderColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	defer op.Offset(image.Pt(borderSize, borderSize)).Push(gtx.Ops).Pop()
	switch b.Format {
	case FormatRounded:
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size.Sub(image.Pt(borderSize*2, borderSize*2))}, gtx.Dp(4)).Push(gtx.Ops).Pop()
	case FormatRectangle:
		defer clip.Rect{Max: dims.Size.Sub(image.Pt(borderSize*2, borderSize*2))}.Push(gtx.Ops).Pop()
	case FormatPill:
		defer clip.UniformRRect(image.Rectangle{Max: dims.Size.Sub(image.Pt(borderSize*2, borderSize*2))}, dims.Size.Y/2).Push(gtx.Ops).Pop()
	}

	paint.ColorOp{Color: b.BackgroundColor}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	pointer.add(gtx.Ops)
	call.Add(gtx.Ops)

	return dims
}

// Pointer is a pointer handler.
type Pointer struct {
	clickFrame time.Time
	pid        pointer.ID
	clicked    bool
	pressed    bool
	entered    bool
	handler    bool
}

func (e *Pointer) add(ops *op.Ops) {
	pointer.InputOp{Tag: e, Types: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave | pointer.Cancel}.Add(ops)
	pointer.CursorPointer.Add(ops)
}

// Clicked reports whether the button was clicked in the last frame.
// It is safe to call Clicked multiple times in the same frame.
func (e *Pointer) Clicked(gtx layout.Context) bool {
	if gtx.Now == e.clickFrame {
		return e.clicked
	}
	e.clickFrame = gtx.Now
	e.clicked = false

	for _, ev := range gtx.Events(e) {
		switch evt := ev.(type) {
		case pointer.Event:
			switch evt.Type {
			case pointer.Release:
				if !e.pressed || e.pid != evt.PointerID {
					break
				}
				e.pressed = false
				if !e.entered {
					break
				}
				e.clicked = true
			case pointer.Cancel:
				e.pressed = false
				e.entered = false
			case pointer.Press:
				if e.pressed {
					break
				}
				if evt.Source == pointer.Mouse && !evt.Buttons.Contain(pointer.ButtonPrimary) {
					break
				}
				if evt.Source == pointer.Touch {
					e.entered = true
				}
				if !e.entered {
					e.pid = evt.PointerID
				}
				if e.pid != evt.PointerID {
					break
				}
				e.pressed = true
			case pointer.Leave:
				if !e.pressed {
					e.pid = evt.PointerID
				}
				if e.pid == evt.PointerID {
					e.entered = false
				}
			case pointer.Enter:
				if !e.pressed {
					e.pid = evt.PointerID
				}
				if e.pid == evt.PointerID {
					e.entered = true
				}
			}
		}
	}

	return e.clicked
}
