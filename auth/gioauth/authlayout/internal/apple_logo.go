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
var VectorAppleLogo giosvg.Vector = func(ops *op.Ops, constraints giosvg.Constraints) layout.Dimensions {
	var w, h float32
	if constraints.Max != constraints.Min {

		d := float32(1.219512)
		if constraints.Max.X*d > constraints.Max.Y {
			w, h = constraints.Max.Y/d, constraints.Max.Y
		} else {
			w, h = constraints.Max.X, constraints.Max.X*d
		}
	}

	if constraints.Min.X > w {
		w = constraints.Min.X
	}
	if constraints.Min.Y > h {
		h = constraints.Min.Y
	}

	var (
		size = f32.Point{X: w / 82.000000, Y: h / 100.000000}
		avg  = (size.X + size.Y) / 2
		aff  = f32.Affine2D{}.Scale(f32.Point{X: float32(0 - 0.000000), Y: float32(0 - 0.000000)}, size)

		end             clip.PathSpec
		path            clip.Path
		stroke, outline clip.Stack
	)
	_, _, _, _, _, _ = avg, aff, end, path, stroke, outline

	path = clip.Path{}
	path.Begin(ops)
	path.MoveTo(aff.Transform(f32.Point{X: 41.894001, Y: 23.077000}))
	path.CubeTo(aff.Transform(f32.Point{X: 46.408001, Y: 23.077000}), aff.Transform(f32.Point{X: 52.066002, Y: 20.025000}), aff.Transform(f32.Point{X: 55.435001, Y: 15.957000}))
	path.CubeTo(aff.Transform(f32.Point{X: 58.487000, Y: 12.270000}), aff.Transform(f32.Point{X: 60.712002, Y: 7.120000}), aff.Transform(f32.Point{X: 60.712002, Y: 1.971000}))
	path.CubeTo(aff.Transform(f32.Point{X: 60.712002, Y: 1.271000}), aff.Transform(f32.Point{X: 60.647999, Y: 0.572000}), aff.Transform(f32.Point{X: 60.521000, Y: 0.000000}))
	path.CubeTo(aff.Transform(f32.Point{X: 55.499001, Y: 0.191000}), aff.Transform(f32.Point{X: 49.459999, Y: 3.369000}), aff.Transform(f32.Point{X: 45.835999, Y: 7.629000}))
	path.CubeTo(aff.Transform(f32.Point{X: 42.974998, Y: 10.871000}), aff.Transform(f32.Point{X: 40.368999, Y: 15.957000}), aff.Transform(f32.Point{X: 40.368999, Y: 21.170000}))
	path.CubeTo(aff.Transform(f32.Point{X: 40.368999, Y: 21.933001}), aff.Transform(f32.Point{X: 40.495998, Y: 22.695000}), aff.Transform(f32.Point{X: 40.558998, Y: 22.950001}))
	path.CubeTo(aff.Transform(f32.Point{X: 40.876999, Y: 23.013000}), aff.Transform(f32.Point{X: 41.386002, Y: 23.077000}), aff.Transform(f32.Point{X: 41.894001, Y: 23.077000}))
	path.MoveTo(aff.Transform(f32.Point{X: 26.000999, Y: 100.000000}))
	path.CubeTo(aff.Transform(f32.Point{X: 32.167999, Y: 100.000000}), aff.Transform(f32.Point{X: 34.901001, Y: 95.867996}), aff.Transform(f32.Point{X: 42.594002, Y: 95.867996}))
	path.CubeTo(aff.Transform(f32.Point{X: 50.412998, Y: 95.867996}), aff.Transform(f32.Point{X: 52.130001, Y: 99.873001}), aff.Transform(f32.Point{X: 58.995998, Y: 99.873001}))
	path.CubeTo(aff.Transform(f32.Point{X: 65.734001, Y: 99.873001}), aff.Transform(f32.Point{X: 70.248001, Y: 93.642998}), aff.Transform(f32.Point{X: 74.507004, Y: 87.540001}))
	path.CubeTo(aff.Transform(f32.Point{X: 79.275002, Y: 80.546997}), aff.Transform(f32.Point{X: 81.246002, Y: 73.681000}), aff.Transform(f32.Point{X: 81.373001, Y: 73.362999}))
	path.CubeTo(aff.Transform(f32.Point{X: 80.928001, Y: 73.236000}), aff.Transform(f32.Point{X: 68.023003, Y: 67.959000}), aff.Transform(f32.Point{X: 68.023003, Y: 53.146999}))
	path.CubeTo(aff.Transform(f32.Point{X: 68.023003, Y: 40.305000}), aff.Transform(f32.Point{X: 78.195000, Y: 34.520000}), aff.Transform(f32.Point{X: 78.766998, Y: 34.075001}))
	path.CubeTo(aff.Transform(f32.Point{X: 72.028000, Y: 24.412001}), aff.Transform(f32.Point{X: 61.792999, Y: 24.158001}), aff.Transform(f32.Point{X: 58.995998, Y: 24.158001}))
	path.CubeTo(aff.Transform(f32.Point{X: 51.430000, Y: 24.158001}), aff.Transform(f32.Point{X: 45.264000, Y: 28.735001}), aff.Transform(f32.Point{X: 41.386002, Y: 28.735001}))
	path.CubeTo(aff.Transform(f32.Point{X: 37.189999, Y: 28.735001}), aff.Transform(f32.Point{X: 31.659000, Y: 24.412001}), aff.Transform(f32.Point{X: 25.111000, Y: 24.412001}))
	path.CubeTo(aff.Transform(f32.Point{X: 12.651000, Y: 24.412001}), aff.Transform(f32.Point{X: -0.000000, Y: 34.710999}), aff.Transform(f32.Point{X: -0.000000, Y: 54.164001}))
	path.CubeTo(aff.Transform(f32.Point{X: -0.000000, Y: 66.242996}), aff.Transform(f32.Point{X: 4.704000, Y: 79.021004}), aff.Transform(f32.Point{X: 10.490000, Y: 87.285004}))
	path.CubeTo(aff.Transform(f32.Point{X: 15.448000, Y: 94.278000}), aff.Transform(f32.Point{X: 19.771000, Y: 100.000000}), aff.Transform(f32.Point{X: 26.000999, Y: 100.000000}))
	end = path.End()
	outline = clip.Outline{Path: end}.Op().Push(ops)
	paint.PaintOp{}.Add(ops)
	outline.Pop()
	return layout.Dimensions{Size: image.Point{X: int(w), Y: int(h)}}
}
