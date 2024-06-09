package giowebview

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
	"image"
	"reflect"
)

var wantOps = []reflect.Type{
	reflect.TypeOf(&WebViewOp{}),
	reflect.TypeOf(&OffsetOp{}),
	reflect.TypeOf(&RectOp{}),
}

// WebViewOp shows the webview into the specified area.
// The RectOp is not context-aware, and will overlay
// any other widget on the screen.
//
// WebViewOp also takes the foreground and clicks events
// and keyboard events will not be routed to Gio.
//
// Performance: changing the size/bounds or radius can
// be expensive. If applicable, change the Offset, instead
// of changing the size.
type WebViewOp struct {
	Tag   event.Tag
	isPop bool
}

// OffsetOp moves the webview by the specified offset.
type OffsetOp struct {
	Point f32.Point
}

// RectOp shows the webview into the specified area.
// The RectOp is not context-aware, and will overlay
// any other widget on the screen.
//
// RectOp also takes the foreground and clicks events
// and keyboard events will not be routed to Gio.
//
// Performance: changing the size/bounds or radius can
// be expensive. If applicable, change the Rect, instead
// of changing the size.
//
// Only one RectOp can be active at each frame for the
// same WebViewOp.
type RectOp struct {
	Size           f32.Point
	SE, SW, NW, NE float32
}

var _WebViewOpPool = plugin.NewOpPool[WebViewOp]()
var _OffsetOpPool = plugin.NewOpPool[OffsetOp]()
var _RectOpPool = plugin.NewOpPool[RectOp]()

// Push adds a new WebViewOp to the queue, any subsequent Ops (sucha as RectOp)
// will affect this WebViewOp.
// In order to stop using this WebViewOp, call Pop.
func (o WebViewOp) Push(op *op.Ops) WebViewOp {
	o.isPop = false
	opc := _WebViewOpPool.Get()
	*opc = o

	plugin.WriteOp(op, opc)
	return o
}

// Pop stops using the WebViewOp.
func (o WebViewOp) Pop(op *op.Ops) {
	o.isPop = true
	opc := _WebViewOpPool.Get()
	*opc = o

	plugin.WriteOp(op, opc)
}

func (o WebViewOp) execute(w *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	runnerIndex, ok := p.tags[o.Tag]
	if !ok {
		p.config = NewConfigFromViewEvent(w, p.viewEvent)
		wv, err := webview.NewWebView(p.config)
		if err != nil {
			panic(err)
		}
		go eventsListener(wv, w, p, o.Tag)
		runnerIndex = len(p.views)
		p.views = append(p.views, wv)
		p.tags[o.Tag] = runnerIndex
		p.seem = append(p.seem, false)
		p.bounds = append(p.bounds, [2]f32.Point{})
	}

	runner := p.views[runnerIndex]

	if o.isPop {
		p.active = nil
		p.activeIndex = 0
		p.activeTag = nil
	} else {
		p.activeIndex = runnerIndex
		p.active = runner
		p.activeTag = o.Tag
	}
}

// NewOffsetOp creates a new OffsetOp.
func NewOffsetOp[POINT image.Point | f32.Point](v POINT) OffsetOp {
	switch v := any(v).(type) {
	case image.Point:
		return OffsetOp{Point: f32.Point{X: float32(v.X), Y: float32(v.Y)}}
	case f32.Point:
		return OffsetOp{Point: v}
	default:
		return OffsetOp{}
	}
}

// Add adds a new OffsetOp to the queue.
func (o OffsetOp) Add(op *op.Ops) {
	opc := _OffsetOpPool.Get()
	*opc = o
	plugin.WriteOp(op, opc)
}

func (o OffsetOp) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	if p.active == nil {
		return
	}
	p.bounds[p.activeIndex][0].Y += o.Point.Y
	p.bounds[p.activeIndex][0].X += o.Point.X
}

// NewRectOp creates a new RectOp.
func NewRectOp[POINT image.Point | f32.Point](v POINT) RectOp {
	switch v := any(v).(type) {
	case image.Point:
		return RectOp{Size: f32.Point{X: float32(v.X), Y: float32(v.Y)}}
	case f32.Point:
		return RectOp{Size: v}
	default:
		return RectOp{}
	}
}

// Add adds a new RectOp to the queue.
func (o RectOp) Add(op *op.Ops) {
	opc := _RectOpPool.Get()
	*opc = o

	plugin.WriteOp(op, opc)
}

func (o RectOp) execute(_ *app.Window, p *webViewPlugin, e app.FrameEvent) {
	if p.active == nil {
		return
	}

	p.mutex.Lock()
	defer p.mutex.Unlock()

	if e.Metric.PxPerDp != p.config.PxPerDp {
		p.config.PxPerDp = e.Metric.PxPerDp
		p.active.Configure(p.config)
	}

	p.seem[p.activeIndex] = true

	p.bounds[p.activeIndex][1].X += o.Size.X
	p.bounds[p.activeIndex][1].Y += o.Size.Y

	p.bounds[p.activeIndex][0].X += float32(unit.Dp(e.Metric.PxPerDp) * e.Insets.Left)
	p.bounds[p.activeIndex][0].Y += float32(unit.Dp(e.Metric.PxPerDp) * e.Insets.Top)

	p.active.Resize(webview.Point{X: p.bounds[p.activeIndex][1].X, Y: p.bounds[p.activeIndex][1].Y}, webview.Point{X: p.bounds[p.activeIndex][0].X, Y: p.bounds[p.activeIndex][0].Y})
}
