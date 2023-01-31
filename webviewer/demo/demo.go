package main

import (
	"flag"
	"image"
	"image/color"
	"math"
	"net/url"
	"os"
	"time"

	"golang.org/x/exp/shiny/materialdesign/icons"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

var (
	GlobalShaper = text.NewCache(gofont.Collection())
	DefaultURL   = "https://google.com"

	IconAdd, _            = widget.NewIcon(icons.ContentAdd)
	IconClose, _          = widget.NewIcon(icons.NavigationClose)
	IconGo, _             = widget.NewIcon(icons.NavigationArrowForward)
	IconCookie, _         = widget.NewIcon(icons.ContentArchive)
	IconLocalStorage, _   = widget.NewIcon(icons.DeviceStorage)
	IconSessionStorage, _ = widget.NewIcon(icons.ImageTimer)
	IconJavascript, _     = widget.NewIcon(icons.AVPlayArrow)
)

func main() {
	proxy := flag.String("proxy", "", "proxy")
	if proxy != nil && *proxy != "" {
		u, err := url.Parse(*proxy)
		if err != nil {
			panic(err)
		}
		if err := webview.SetProxy(u); err != nil {
			panic(err)
		}
	}
	flag.Parse()

	webview.SetDebug(true)
	window := app.NewWindow()

	browsers := NewBrowser()
	browsers.add()

	go func() {
		ops := new(op.Ops)
		// first := true
		for evt := range window.Events() {
			plugin.Install(window, evt)

			switch evt := evt.(type) {
			case system.DestroyEvent:
				os.Exit(0)
				return
			case system.FrameEvent:
				gtx := layout.NewContext(ops, evt)
				browsers.Layout(gtx)
				evt.Frame(ops)
			}
		}
	}()

	app.Main()
}

const (
	VisibleLocal = 1 << iota
	VisibleSession
	VisibleCookies
)

type Browsers struct {
	Selected int

	Go    widget.Clickable
	Add   widget.Clickable
	Close widget.Clickable

	JavascriptCode widget.Editor
	JavascriptRun  widget.Clickable

	Tabs    []widget.Clickable
	Address []widget.Editor

	Tags   []*int
	Titles []string

	LocalStorage   [][]webview.StorageData
	SessionStorage [][]webview.StorageData
	CookieStorage  [][]webview.CookieData

	StorageVisible uint8

	LocalButton   widget.Clickable
	SessionButton widget.Clickable
	CookieButton  widget.Clickable

	HeaderFlex []layout.FlexChild
	TabsFlex   []layout.FlexChild
}

func NewBrowser() *Browsers {
	b := &Browsers{}
	b.HeaderFlex = []layout.FlexChild{
		layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
			defer clip.Outline{Path: clip.Rect{Max: gtx.Constraints.Max}.Path()}.Op().Push(gtx.Ops).Pop()
			paint.ColorOp{Color: color.NRGBA{R: 24, G: 26, B: 33, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			gtx.Constraints.Min.Y = 0
			gtx.Constraints.Max.X -= gtx.Dp(16)
			macro := op.Record(gtx.Ops)
			dims := b.Address[b.Selected].Layout(gtx, GlobalShaper, text.Font{}, gtx.Metric.DpToSp(16), func(gtx layout.Context) layout.Dimensions {
				paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
				b.Address[b.Selected].PaintText(gtx)

				paint.ColorOp{Color: color.NRGBA{R: 123, G: 123, B: 123, A: 255}}.Add(gtx.Ops)
				b.Address[b.Selected].PaintCaret(gtx)

				paint.ColorOp{Color: color.NRGBA{R: 123, G: 123, B: 123, A: 255}}.Add(gtx.Ops)
				b.Address[b.Selected].PaintSelection(gtx)

				return layout.Dimensions{Size: gtx.Constraints.Max}
			})
			call := macro.Stop()

			defer op.Offset(image.Point{X: gtx.Dp(8), Y: (gtx.Constraints.Max.Y - dims.Size.Y - dims.Baseline) / 2}).Push(gtx.Ops).Pop()
			call.Add(gtx.Ops)

			gtx.Constraints.Max.X += gtx.Dp(16)

			return layout.Dimensions{Size: gtx.Constraints.Max}
		}),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.Go, Icon: IconGo, Text: "Go"}.Layout),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.Close, Icon: IconClose, Text: "Close"}.Layout),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.Add, Icon: IconAdd, Text: "Add"}.Layout),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.CookieButton, Icon: IconCookie}.Layout),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.LocalButton, Icon: IconLocalStorage}.Layout),
		layout.Rigid(layout.Spacer{Width: 4}.Layout),
		layout.Rigid(Button{Clickable: &b.SessionButton, Icon: IconSessionStorage}.Layout),
	}
	return b
}

func (b *Browsers) add() {
	b.Tabs = append(b.Tabs, widget.Clickable{})
	b.Tags = append(b.Tags, new(int))
	b.Titles = append(b.Titles, "")
	b.Address = append(b.Address, widget.Editor{SingleLine: true, Submit: true})
	b.LocalStorage = append(b.LocalStorage, nil)
	b.SessionStorage = append(b.SessionStorage, nil)
	b.CookieStorage = append(b.CookieStorage, nil)

	if cap(b.TabsFlex) < len(b.Tabs) {
		b.TabsFlex = make([]layout.FlexChild, len(b.Tabs))
	} else {
		b.TabsFlex = b.TabsFlex[:len(b.Tabs)]
	}
}

func (b *Browsers) remove(i int) {
	if len(b.Tabs) == 1 {
		return
	}
	if b.Selected >= len(b.Tabs)-1 {
		b.Selected--
	}
	b.Tabs = append(b.Tabs[:i], b.Tabs[i+1:]...)
	b.Tags = append(b.Tags[:i], b.Tags[i+1:]...)
	b.Titles = append(b.Titles[:i], b.Titles[i+1:]...)
	b.TabsFlex = append(b.TabsFlex[:i], b.TabsFlex[i+1:]...)
	b.Address = append(b.Address[:i], b.Address[i+1:]...)
	b.SessionStorage = append(b.SessionStorage[:i], b.SessionStorage[i+1:]...)
	b.LocalStorage = append(b.LocalStorage[:i], b.LocalStorage[i+1:]...)
	b.CookieStorage = append(b.CookieStorage[:i], b.CookieStorage[i+1:]...)
}

func (b *Browsers) Layout(gtx layout.Context) layout.Dimensions {
	if b.Add.Clicked() {
		b.add()
	}
	if b.Close.Clicked() {
		b.remove(b.Selected)
	}

	currentStoragePanel := b.StorageVisible
	if b.LocalButton.Clicked() {
		currentStoragePanel ^= VisibleLocal
	}
	if b.SessionButton.Clicked() {
		currentStoragePanel ^= VisibleSession
	}
	if b.CookieButton.Clicked() {
		currentStoragePanel ^= VisibleCookies
	}
	b.StorageVisible = currentStoragePanel

	submittedIndex := -1
	if b.Go.Clicked() {
		submittedIndex = b.Selected
	}

	for i, t := range b.Address {
		submited := i == submittedIndex
		if !t.Focused() && t.Text() == "" {
			submited = true
			t.SetText(DefaultURL)
		}

		for _, evt := range t.Events() {
			switch evt.(type) {
			case widget.SubmitEvent:
				submited = true
			}
		}

		if submited {
			w := webviewer.WebViewOp{Tag: b.Tags[i]}.Push(gtx.Ops)
			webviewer.NavigateOp{URL: t.Text()}.Add(gtx.Ops)
			w.Pop(gtx.Ops)
		}
	}

	for i, t := range b.Tabs {
		if t.Clicked() {
			b.Selected = i
		}
	}

	for i := range b.Tags {
		for _, evt := range gtx.Events(b.Tags[i]) {
			switch evt := evt.(type) {
			case webviewer.TitleEvent:
				b.Titles[i] = evt.Title
			case webviewer.NavigationEvent:
				b.Address[i].SetText(evt.URL)
			case webviewer.CookiesEvent:
				// fmt.Println(evt.Cookies)
			case webviewer.StorageEvent:
				// fmt.Println(evt.Storage)
			case webviewer.MessageEvent:
				// fmt.Println(evt.Message)
			}
		}
	}

	gtxi := gtx
	return Rows{}.Layout(gtx, 4, func(i int, gtx layout.Context) layout.Dimensions {
		switch i {
		case 0:
			gtx.Constraints.Max.Y = gtx.Dp(48)
			gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
			defer clip.Outline{Path: clip.Rect{Max: gtx.Constraints.Max}.Path()}.Op().Push(gtx.Ops).Pop()
			paint.ColorOp{Color: color.NRGBA{R: 48, G: 52, B: 67, A: 255}}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)

			gtx.Constraints.Max.Y = gtx.Dp(40)
			gtx.Constraints.Max.X = gtx.Constraints.Max.X - gtx.Dp(20)
			defer op.Offset(image.Point{X: gtx.Dp(20) / 2, Y: 4}).Push(gtx.Ops).Pop()

			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, b.HeaderFlex...)
		case 1:
			gtx.Constraints.Max.Y = gtx.Dp(38)

			b.TabsFlex = b.TabsFlex[:len(b.Tags)]
			for i := range b.Tags {
				i := i
				b.TabsFlex[i] = layout.Flexed(1, func(gtx layout.Context) layout.Dimensions {
					gtx.Constraints.Min = gtx.Constraints.Max

					defer clip.Outline{Path: clip.Rect{Max: gtx.Constraints.Max}.Path()}.Op().Push(gtx.Ops).Pop()
					if b.Selected == i {
						paint.ColorOp{Color: color.NRGBA{R: 48, G: 52, B: 67, A: 255}}.Add(gtx.Ops)
					} else {
						paint.ColorOp{Color: color.NRGBA{R: 61, G: 61, B: 69, A: 255}}.Add(gtx.Ops)
					}
					paint.PaintOp{}.Add(gtx.Ops)

					return b.Tabs[i].Layout(gtx, func(gtx layout.Context) layout.Dimensions {
						return layout.UniformInset(4).Layout(gtx, func(gtx layout.Context) layout.Dimensions {
							paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
							macro := op.Record(gtx.Ops)
							gtx.Constraints.Min.Y = 0
							dims := widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, GlobalShaper, text.Font{}, gtx.Metric.DpToSp(16), b.Titles[i])
							call := macro.Stop()

							defer op.Offset(image.Point{X: 0, Y: (gtx.Constraints.Max.Y - dims.Size.Y) / 2}).Push(gtx.Ops).Pop()
							call.Add(gtx.Ops)
							return dims
						})
					})
				})
			}
			return layout.Flex{Axis: layout.Horizontal}.Layout(gtx, b.TabsFlex...)
		case 2:
			defer webviewer.WebViewOp{Tag: b.Tags[b.Selected]}.Push(gtx.Ops).Pop(gtx.Ops)
			webviewer.OffsetOp{Point: f32.Point{Y: float32(gtxi.Constraints.Max.Y - gtx.Constraints.Max.Y)}}.Add(gtx.Ops)
			webviewer.RectOp{Size: f32.Point{X: float32(gtx.Constraints.Max.X), Y: float32(gtx.Constraints.Max.Y)}}.Add(gtx.Ops)
			return layout.Dimensions{Size: gtx.Constraints.Max}
		default:
			return layout.Dimensions{}
		}
	})
}

type Rows struct {
	Size layout.Dimensions
}

func (r Rows) Layout(gtx layout.Context, n int, fn func(i int, gtx layout.Context) layout.Dimensions) layout.Dimensions {
	for i := 0; i < n; i++ {
		offset := op.Offset(image.Point{Y: r.Size.Size.Y}).Push(gtx.Ops)
		dims := fn(i, gtx)
		if dims.Size.X > r.Size.Size.X {
			r.Size.Size.X = dims.Size.X
		}
		r.Size.Size.Y += dims.Size.Y
		gtx.Constraints.Max.Y -= dims.Size.Y
		offset.Pop()
	}
	return r.Size
}

type Columns struct {
	Size layout.Dimensions
}

func (c Columns) Layout(gtx layout.Context, n int, fn func(i int, gtx layout.Context) layout.Dimensions) layout.Dimensions {
	for i := 0; i < n; i++ {
		offset := op.Offset(image.Point{X: c.Size.Size.X}).Push(gtx.Ops)
		dims := fn(i, gtx)
		if dims.Size.Y > c.Size.Size.Y {
			c.Size.Size.Y = dims.Size.Y
		}
		c.Size.Size.X += dims.Size.X
		gtx.Constraints.Max.X -= dims.Size.X
		offset.Pop()
	}
	return c.Size
}

type Button struct {
	Clickable *widget.Clickable
	Icon      *widget.Icon
	Text      string
}

func (b Button) Layout(gtx layout.Context) layout.Dimensions {
	return b.Clickable.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		macro := op.Record(gtx.Ops)
		c := color.NRGBA{R: 32, G: 32, B: 32, A: 255}
		paint.ColorOp{Color: c}.Add(gtx.Ops)
		gtx.Constraints.Min.Y = 0
		var dims layout.Dimensions
		if b.Icon == nil {
			dims = widget.Label{Alignment: text.Start, MaxLines: 1}.Layout(gtx, GlobalShaper, text.Font{}, gtx.Metric.DpToSp(16), b.Text)
		} else {
			gtx := gtx
			gtx.Constraints.Max.Y = gtx.Sp(gtx.Metric.DpToSp(16))
			gtx.Constraints.Max.X = gtx.Constraints.Max.Y
			dims = b.Icon.Layout(gtx, c)
		}
		call := macro.Stop()

		gtx.Constraints.Max.X = dims.Size.X + gtx.Dp(16)

		defer clip.Outline{Path: clip.Rect{Max: gtx.Constraints.Max}.Path()}.Op().Push(gtx.Ops).Pop()
		paint.ColorOp{Color: color.NRGBA{R: 237, G: 237, B: 237, A: 255}}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		pointer.CursorPointer.Add(gtx.Ops)

		defer op.Offset(image.Point{X: gtx.Dp(8), Y: (gtx.Constraints.Max.Y - dims.Size.Y) / 2}).Push(gtx.Ops).Pop()
		call.Add(gtx.Ops)

		return layout.Dimensions{Size: gtx.Constraints.Max}
	})
}

type Loading struct {
	loadingLast time.Time
	loadingDt   float32
}

func (l *Loading) Layout(gtx layout.Context) layout.Dimensions {
	gtx.Constraints.Max.X = gtx.Constraints.Max.Y

	diff := gtx.Now.Sub(l.loadingLast)
	dt := float32(math.Round(float64(diff/(time.Millisecond*32))) * 0.032)
	l.loadingDt += dt
	if l.loadingDt >= 1 {
		l.loadingDt = 0
	}
	if dt > 0 {
		l.loadingLast = gtx.Now
	}

	width := float32(gtx.Dp(4))

	radius := float32(gtx.Constraints.Max.Y / 5)
	defer op.Affine(f32.Affine2D{}.Offset(f32.Point{
		X: float32(gtx.Constraints.Max.X/2) - (radius / 2) + (width),
		Y: float32(gtx.Constraints.Max.Y/2) - (radius / 2) + (width),
	})).Push(gtx.Ops).Pop()

	rot := f32.Affine2D{}.Rotate(f32.Pt(0, 0), l.loadingDt*math.Pi*2)

	path := clip.Path{}
	path.Begin(gtx.Ops)
	path.Move(rot.Transform(f32.Pt(radius, radius)))
	path.Arc(
		rot.Transform(f32.Pt(-radius, -radius)),
		rot.Transform(f32.Pt(-radius, -radius)),
		float32((math.Pi*2)/8)*7,
	)

	defer clip.Stroke{Path: path.End(), Width: width}.Op().Push(gtx.Ops).Pop()
	paint.ColorOp{Color: color.NRGBA{R: 255, G: 255, B: 255, A: 255}}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	op.InvalidateOp{}.Add(gtx.Ops)

	return layout.Dimensions{Size: gtx.Constraints.Max}
}
