package internal

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/inkeliz/giosvg"
)

var _, _, _, _, _, _, _, _ = (*f32.Point)(nil), (*op.Ops)(nil), (*clip.Op)(nil), (*paint.PaintOp)(nil), (*giosvg.Vector)(nil), (*color.NRGBA)(nil), (*layout.Dimensions)(nil), (*image.Image)(nil)
var VectorGoogleLogo giosvg.Vector = func(ops *op.Ops, constraints giosvg.Constraints) layout.Dimensions {
	var w, h float32
	if constraints.Max != constraints.Min {

		d := float32(1.000000)
		if constraints.Max.Y*d > constraints.Max.X {
			w, h = constraints.Max.X, constraints.Max.X/d
		} else {
			w, h = constraints.Max.Y*d, constraints.Max.Y
		}
	}

	if constraints.Min.X > w {
		w = constraints.Min.X
	}
	if constraints.Min.Y > h {
		h = constraints.Min.Y
	}

	var (
		size = f32.Point{X: w / 75.000000, Y: h / 75.000000}
		avg  = (size.X + size.Y) / 2
		aff  = f32.Affine2D{}.Scale(f32.Point{X: float32(0 - 0.000000), Y: float32(0 - 0.000000)}, size)

		end             clip.PathSpec
		path            clip.Path
		stroke, outline clip.Stack
	)
	_, _, _, _, _, _ = avg, aff, end, path, stroke, outline

	path = clip.Path{}
	path.Begin(ops)
	path.MoveTo(aff.Transform(f32.Point{X: 73.500000, Y: 38.352001}))
	path.CubeTo(aff.Transform(f32.Point{X: 73.500000, Y: 35.693001}), aff.Transform(f32.Point{X: 73.261002, Y: 33.136002}), aff.Transform(f32.Point{X: 72.818001, Y: 30.681999}))
	path.LineTo(aff.Transform(f32.Point{X: 37.500000, Y: 30.681999}))
	path.LineTo(aff.Transform(f32.Point{X: 37.500000, Y: 45.188000}))
	path.LineTo(aff.Transform(f32.Point{X: 57.681999, Y: 45.188000}))
	path.CubeTo(aff.Transform(f32.Point{X: 56.812000, Y: 49.875000}), aff.Transform(f32.Point{X: 54.169998, Y: 53.847000}), aff.Transform(f32.Point{X: 50.199001, Y: 56.506001}))
	path.LineTo(aff.Transform(f32.Point{X: 50.199001, Y: 65.915001}))
	path.LineTo(aff.Transform(f32.Point{X: 62.318001, Y: 65.915001}))
	path.CubeTo(aff.Transform(f32.Point{X: 69.408997, Y: 59.386002}), aff.Transform(f32.Point{X: 73.500000, Y: 49.772999}), aff.Transform(f32.Point{X: 73.500000, Y: 38.352001}))
	end = path.End()
	outline = clip.Outline{Path: end}.Op().Push(ops)
	paint.ColorOp{Color: color.NRGBA{R: 66, G: 133, B: 244, A: 255}}.Add(ops)
	paint.PaintOp{}.Add(ops)
	outline.Pop()

	path = clip.Path{}
	path.Begin(ops)
	path.MoveTo(aff.Transform(f32.Point{X: 37.500000, Y: 75.000000}))
	path.CubeTo(aff.Transform(f32.Point{X: 47.625000, Y: 75.000000}), aff.Transform(f32.Point{X: 56.113998, Y: 71.641998}), aff.Transform(f32.Point{X: 62.318001, Y: 65.915001}))
	path.LineTo(aff.Transform(f32.Point{X: 50.199001, Y: 56.506001}))
	path.CubeTo(aff.Transform(f32.Point{X: 46.841000, Y: 58.756001}), aff.Transform(f32.Point{X: 42.544998, Y: 60.084999}), aff.Transform(f32.Point{X: 37.500000, Y: 60.084999}))
	path.CubeTo(aff.Transform(f32.Point{X: 27.733000, Y: 60.084999}), aff.Transform(f32.Point{X: 19.466000, Y: 53.488998}), aff.Transform(f32.Point{X: 16.517000, Y: 44.625000}))
	path.LineTo(aff.Transform(f32.Point{X: 3.989000, Y: 44.625000}))
	path.LineTo(aff.Transform(f32.Point{X: 3.989000, Y: 54.341000}))
	path.CubeTo(aff.Transform(f32.Point{X: 10.159000, Y: 66.597000}), aff.Transform(f32.Point{X: 22.841000, Y: 75.000000}), aff.Transform(f32.Point{X: 37.500000, Y: 75.000000}))
	end = path.End()
	outline = clip.Outline{Path: end}.Op().Push(ops)
	paint.ColorOp{Color: color.NRGBA{R: 52, G: 168, B: 83, A: 255}}.Add(ops)
	paint.PaintOp{}.Add(ops)
	outline.Pop()

	path = clip.Path{}
	path.Begin(ops)
	path.MoveTo(aff.Transform(f32.Point{X: 16.517000, Y: 44.625000}))
	path.CubeTo(aff.Transform(f32.Point{X: 15.767000, Y: 42.375000}), aff.Transform(f32.Point{X: 15.341000, Y: 39.972000}), aff.Transform(f32.Point{X: 15.341000, Y: 37.500000}))
	path.CubeTo(aff.Transform(f32.Point{X: 15.341000, Y: 35.028000}), aff.Transform(f32.Point{X: 15.767000, Y: 32.625000}), aff.Transform(f32.Point{X: 16.517000, Y: 30.375000}))
	path.LineTo(aff.Transform(f32.Point{X: 16.517000, Y: 20.659000}))
	path.LineTo(aff.Transform(f32.Point{X: 3.989000, Y: 20.659000}))
	path.CubeTo(aff.Transform(f32.Point{X: 1.449000, Y: 25.722000}), aff.Transform(f32.Point{X: -0.000000, Y: 31.448999}), aff.Transform(f32.Point{X: -0.000000, Y: 37.500000}))
	path.CubeTo(aff.Transform(f32.Point{X: -0.000000, Y: 43.550999}), aff.Transform(f32.Point{X: 1.449000, Y: 49.278000}), aff.Transform(f32.Point{X: 3.989000, Y: 54.341000}))
	path.LineTo(aff.Transform(f32.Point{X: 16.517000, Y: 44.625000}))
	end = path.End()
	outline = clip.Outline{Path: end}.Op().Push(ops)
	paint.ColorOp{Color: color.NRGBA{R: 251, G: 188, B: 5, A: 255}}.Add(ops)
	paint.PaintOp{}.Add(ops)
	outline.Pop()

	path = clip.Path{}
	path.Begin(ops)
	path.MoveTo(aff.Transform(f32.Point{X: 37.500000, Y: 14.915000}))
	path.CubeTo(aff.Transform(f32.Point{X: 43.006001, Y: 14.915000}), aff.Transform(f32.Point{X: 47.949001, Y: 16.806999}), aff.Transform(f32.Point{X: 51.834999, Y: 20.523001}))
	path.LineTo(aff.Transform(f32.Point{X: 62.591000, Y: 9.767000}))
	path.CubeTo(aff.Transform(f32.Point{X: 56.097000, Y: 3.716000}), aff.Transform(f32.Point{X: 47.608002, Y: 0.000000}), aff.Transform(f32.Point{X: 37.500000, Y: 0.000000}))
	path.CubeTo(aff.Transform(f32.Point{X: 22.841000, Y: 0.000000}), aff.Transform(f32.Point{X: 10.159000, Y: 8.403000}), aff.Transform(f32.Point{X: 3.989000, Y: 20.659000}))
	path.LineTo(aff.Transform(f32.Point{X: 16.517000, Y: 30.375000}))
	path.CubeTo(aff.Transform(f32.Point{X: 19.466000, Y: 21.511000}), aff.Transform(f32.Point{X: 27.733000, Y: 14.915000}), aff.Transform(f32.Point{X: 37.500000, Y: 14.915000}))
	end = path.End()
	outline = clip.Outline{Path: end}.Op().Push(ops)
	paint.ColorOp{Color: color.NRGBA{R: 234, G: 67, B: 53, A: 255}}.Add(ops)
	paint.PaintOp{}.Add(ops)
	outline.Pop()
	return layout.Dimensions{Size: image.Point{X: int(w), Y: int(h)}}
}
