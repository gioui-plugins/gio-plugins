package authlayout

import (
	"gioui.org/font"
	"gioui.org/io/event"
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

// ButtonTexts is a list of texts for the button.
// The first is preferred, the second is fallback,
// when the first doesn't fit the button width.
type ButtonTexts [2]string

// ButtonStyle is the style for Google buttons.
//
// Notice that you may violate some branding guidelines, it's
// safer to use DefaultGoogleButtonStyle or DefaultAppleButtonStyle.
type ButtonStyle struct {
	// TextSize is the size of the text.
	TextSize unit.Dp
	// TextFont is the font to use for the text.
	TextFont font.Font
	// TextShaper is the shaper to use for the text.
	TextShaper *text.Shaper
	// TextColor is the color of the text.
	TextColor color.NRGBA
	// TextAlignment is the alignment of the text.
	TextAlignment layout.Alignment
	// IconAlignment is the alignment of the icon.
	IconAlignment layout.Alignment
	// IconColor is the color of the icon.
	IconColor color.NRGBA
	// IconSize is the size of the icon.
	IconSize unit.Dp
	// IconPadding is the padding of the icon.
	IconPadding unit.Dp
	// IconVector is the vector of the icon.
	IconVector giosvg.Vector
	// BackgroundColor is the color of the background.
	BackgroundColor color.NRGBA
	// BackgroundIconColor is the color of the background of the icon.
	BackgroundIconColor color.NRGBA
	// BorderColor is the color of the border.
	BorderColor color.NRGBA
	// BorderThickness is the thickness of the border.
	BorderThickness unit.Dp
	// Format is the format of the button.
	Format Format

	icon *giosvg.Icon
}

func (b *ButtonStyle) label(gtx layout.Context, text string) (call op.CallOp, dims layout.Dimensions) {
	gtx.Constraints.Min.X = 0
	gtx.Constraints.Min.Y = 0

	r := op.Record(gtx.Ops)
	dims = widget.Label{}.Layout(gtx, b.TextShaper, b.TextFont, gtx.Metric.DpToSp(b.TextSize), text, toLabelColor(gtx, b.TextColor))
	return r.Stop(), dims
}

func toLabelColor(gtx layout.Context, c color.NRGBA) op.CallOp {
	r := op.Record(gtx.Ops)
	paint.ColorOp{Color: c}.Add(gtx.Ops)
	return r.Stop()
}

func (b *ButtonStyle) LayoutText(gtx layout.Context, pointer *Pointer, texts ButtonTexts) layout.Dimensions {
	return b.layoutText(gtx, pointer, texts, 0)
}

func (b *ButtonStyle) layoutText(gtx layout.Context, pointer *Pointer, texts ButtonTexts, textIndex int) layout.Dimensions {
	if b.icon == nil {
		b.icon = giosvg.NewIcon(b.IconVector)
	}

	label, labelDims := b.label(gtx, texts[textIndex])

	minHeight := int(math.Round(float64(labelDims.Size.Y) * 233 / 100))

	inset := layout.Inset{
		Top:    gtx.Metric.PxToDp((minHeight - labelDims.Size.Y) / 2),
		Bottom: gtx.Metric.PxToDp((minHeight - labelDims.Size.Y) / 2),
		Left:   unit.Dp(16),
		Right:  unit.Dp(16),
	}

	iconSize := gtx.Dp(b.IconSize)
	iconPadding := gtx.Dp(b.IconPadding)

	if b.IconSize == 0 {
		iconSize = labelDims.Size.Y
	}

	avalSize := gtx.Constraints.Max.X - (labelDims.Size.X + gtx.Dp(inset.Left) + gtx.Dp(inset.Right) + iconSize + iconPadding)
	isLogoOnly := avalSize < 0
	if avalSize > iconPadding && b.TextAlignment != layout.Start && b.IconAlignment != layout.Middle {
		iconPadding = 0
	}

	if isLogoOnly && textIndex != len(texts)-1 {
		return b.layoutText(gtx, pointer, texts, textIndex+1)
	}

	if isLogoOnly {
		labelDims.Size.X = 0
		iconPadding = 0
	}

	main := op.Record(gtx.Ops)
	dims := inset.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		d := layout.Dimensions{Size: image.Pt(gtx.Constraints.Max.X, labelDims.Size.Y)}

		{
			align := b.IconAlignment
			if isLogoOnly {
				align = layout.Middle
			}

			// Logo
			var off op.TransformStack
			switch align {
			case layout.Start, layout.Baseline:
				off = op.Offset(image.Pt(0, 0)).Push(gtx.Ops)
			case layout.Middle:
				off = op.Offset(image.Pt((gtx.Constraints.Max.X-iconSize-iconPadding-labelDims.Size.X)/2, 0)).Push(gtx.Ops)
			case layout.End:
				off = op.Offset(image.Pt(gtx.Constraints.Max.X-iconSize, 0)).Push(gtx.Ops)
			}

			// Logo Background
			padding := gtx.Dp(6)
			offBackground := op.Offset(image.Pt(-padding/2, -padding/2)).Push(gtx.Ops)
			background := clip.UniformRRect(image.Rectangle{Max: image.Pt(iconSize+padding, iconSize+padding)}, (iconSize+padding)/2).Push(gtx.Ops)
			paint.Fill(gtx.Ops, b.BackgroundIconColor)
			background.Pop()
			offBackground.Pop()

			gtx := gtx
			gtx.Constraints.Min = image.Point{}
			gtx.Constraints.Max.X, gtx.Constraints.Max.Y = iconSize, iconSize
			if b.IconColor.A != 0 {
				paint.ColorOp{Color: b.IconColor}.Add(gtx.Ops)
			}

			iconR := op.Record(gtx.Ops)
			dimsIcon := b.icon.Layout(gtx)
			iconOp := iconR.Stop()

			iconOff := op.Offset(image.Pt((iconSize-dimsIcon.Size.X)/2, (iconSize-dimsIcon.Size.Y)/2)).Push(gtx.Ops)
			iconOp.Add(gtx.Ops)
			iconOff.Pop()

			off.Pop()
		}

		if !isLogoOnly {
			// Text
			gtx := gtx
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - iconSize - iconPadding

			if b.TextAlignment != layout.End {
				defer op.Offset(image.Pt(iconSize+iconPadding, 0)).Push(gtx.Ops).Pop()
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
	event.Op(ops, e)
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

	for {
		evt, ok := gtx.Event(pointer.Filter{Target: e, Kinds: pointer.Press | pointer.Release | pointer.Enter | pointer.Leave | pointer.Cancel})
		if !ok {
			break
		}

		switch evt := evt.(type) {
		case pointer.Event:
			switch evt.Kind {
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
			default:
			}
		}
	}

	return e.clicked
}
