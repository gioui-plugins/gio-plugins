package webviewer

import (
	"container/list"
	"image"
	"net/url"
	"reflect"
	"sync"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"

	"gioui.org/unit"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/op"
)

var (
	wantOp = []reflect.Type{
		reflect.TypeOf(&webViewOp{}),
		reflect.TypeOf(&OffsetOp{}),
		reflect.TypeOf(&RectOp{}),
		reflect.TypeOf(&NavigateOp{}),
		reflect.TypeOf(&DestroyOp{}),
		reflect.TypeOf(&SetCookieOp{}),
		reflect.TypeOf(&ListCookieOp{}),
		reflect.TypeOf(&RemoveCookieOp{}),
		reflect.TypeOf(&SetStorageOp{}),
		reflect.TypeOf(&ListStorageOp{}),
		reflect.TypeOf(&RemoveStorageOp{}),
		reflect.TypeOf(&ClearCacheOp{}),
		reflect.TypeOf(&ExecuteJavascriptOp{}),
		reflect.TypeOf(&InstallJavascriptOp{}),
		reflect.TypeOf(&MessageReceiverOp{}),
	}
	wantEvent = []reflect.Type{
		reflect.TypeOf(app.ViewEvent{}),
		reflect.TypeOf(system.FrameEvent{}),
		reflect.TypeOf(system.DestroyEvent{}),
		reflect.TypeOf(plugin.EndFrameEvent{}),
	}
)

var webViewPluginInstances = struct {
	sync.Mutex
	list []*webViewPlugin
}{}

func init() {
	plugin.Register(func(w *app.Window, plugin *plugin.Plugin) plugin.Handler {
		p := &webViewPlugin{
			funcs:     new(list.List).Init(),
			funcsChan: make(chan struct{}, 1),

			window: w,
			plugin: plugin,

			tags:   make(map[event.Tag]int),
			views:  make([]webview.WebView, 0, 8),
			seem:   make([]bool, 0, 8),
			bounds: make([][2]f32.Point, 0, 8),
		}
		go p.runFuncs()
		webViewPluginInstances.Lock()
		webViewPluginInstances.list = append(webViewPluginInstances.list, p)
		webViewPluginInstances.Unlock()
		return p
	})
}

type webViewPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	funcs     *list.List
	funcsChan chan struct{}

	mutex sync.Mutex

	tags   map[event.Tag]int
	views  []webview.WebView
	seem   []bool
	bounds [][2]f32.Point

	activeIndex int
	activeTag   event.Tag
	active      webview.WebView

	frame system.FrameEvent

	config    webview.Config
	viewEvent app.ViewEvent
}

// TypeOp implements plugin.Handler.
func (p *webViewPlugin) TypeOp() []reflect.Type {
	return wantOp
}

// TypeEvent implements plugin.Handler.
func (p *webViewPlugin) TypeEvent() []reflect.Type {
	return wantEvent
}

// ListenOps implements plugin.Handler.
func (p *webViewPlugin) ListenOps(op interface{}) {
	switch v := op.(type) {
	case interface {
		execute(w *app.Window, p *webViewPlugin, _ system.FrameEvent)
	}:
		v.execute(p.window, p, p.frame)
	default:
		_ = v
	}
}

// ListenEvents implements plugin.Handler.
func (p *webViewPlugin) ListenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		p.viewEvent = evt

	case system.FrameEvent:
		p.frame = evt

		// Reset the seen map.
		for s := range p.seem {
			p.seem[s] = false
		}

	case system.DestroyEvent:
		p.mutex.Lock()
		defer p.mutex.Unlock()

		for _, v := range p.views {
			v.Close()
		}

	case plugin.EndFrameEvent:
		// If remain unseen, makes it invisible (0x0)
		for i, v := range p.seem {
			if !v {
				p.views[i].Resize(webview.Point{}, webview.Point{})
			}
		}

		for i := range p.bounds {
			p.bounds[i][0] = f32.Point{}
			p.bounds[i][1] = f32.Point{}
		}

		p.activeIndex = 0
		p.active = nil
		p.activeTag = nil
	}
}

func (p *webViewPlugin) run(f func()) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.funcs.PushBack(f)

	select {
	case p.funcsChan <- struct{}{}:
	default:
	}
}

func (p *webViewPlugin) runFuncs() {
	for range p.funcsChan {
		for {
			p.mutex.Lock()
			f := p.funcs.Front()
			var fn func()
			if f != nil {
				fn = f.Value.(func())
				p.funcs.Remove(f)
			}
			p.mutex.Unlock()
			if fn != nil {
				fn()
			}
		}
	}
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

type webViewOp struct {
	tag   event.Tag
	isPop bool
}

var poolWebViewOp = plugin.NewOpPool[webViewOp]()

// Push adds a new WebViewOp to the queue, any subsequent Ops (sucha as RectOp)
// will affect this WebViewOp.
// In order to stop using this WebViewOp, call Pop.
func (o WebViewOp) Push(op *op.Ops) WebViewOp {
	o.isPop = false
	poolWebViewOp.WriteOp(op, *(*webViewOp)(unsafe.Pointer(&o)))
	return o
}

// Pop stops using the WebViewOp.
func (o WebViewOp) Pop(op *op.Ops) {
	o.isPop = true
	poolWebViewOp.WriteOp(op, *(*webViewOp)(unsafe.Pointer(&o)))
}

func (o *webViewOp) execute(w *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	defer poolWebViewOp.Release(o)

	p.mutex.Lock()
	defer p.mutex.Unlock()
	runnerIndex, ok := p.tags[o.tag]
	if !ok {
		p.config = NewConfigFromViewEvent(w, p.viewEvent)
		wv, err := webview.NewWebView(p.config)
		if err != nil {
			panic(err)
		}
		go eventsListener(wv, w, p, o.tag)
		runnerIndex = len(p.views)
		p.views = append(p.views, wv)
		p.tags[o.tag] = runnerIndex
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
		p.activeTag = o.tag
	}
}

func eventsListener(wv webview.WebView, w *app.Window, p *webViewPlugin, tag event.Tag) {
	for evt := range wv.Events() {
		switch evt := evt.(type) {
		case webview.NavigationEvent:
			p.plugin.SendEvent(tag, NavigationEvent(evt))
		case webview.TitleEvent:
			p.plugin.SendEvent(tag, TitleEvent(evt))
		}
		w.Invalidate()
	}
}

// OffsetOp moves the webview by the specified offset.
type OffsetOp struct {
	Point f32.Point
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

var poolOffsetOp = plugin.NewOpPool[OffsetOp]()

// Add adds a new OffsetOp to the queue.
func (o OffsetOp) Add(op *op.Ops) {
	poolOffsetOp.WriteOp(op, o)
}

func (o *OffsetOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	defer poolOffsetOp.Release(o)
	if p.active == nil {
		return
	}
	p.bounds[p.activeIndex][0].Y += o.Point.Y
	p.bounds[p.activeIndex][0].X += o.Point.X
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

var poolRectOp = plugin.NewOpPool[RectOp]()

// Add adds a new RectOp to the queue.
func (o RectOp) Add(op *op.Ops) {
	poolRectOp.WriteOp(op, o)
}

func (o *RectOp) execute(_ *app.Window, p *webViewPlugin, e system.FrameEvent) {
	defer poolRectOp.Release(o)

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

// NavigateOp redirects the last Display to the
// given URL. If the URL have unknown protocols,
// or malformed URL may lead to unknown behaviors.
type NavigateOp struct {
	// URL is the URL to redirect to.
	URL string
}

var poolNavigateOp = plugin.NewOpPool[NavigateOp]()

// Add adds a new NavigateOp to the queue.
func (o NavigateOp) Add(op *op.Ops) {
	poolNavigateOp.WriteOp(op, o)
}

func (o *NavigateOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	defer poolNavigateOp.Release(o)

	if p.active == nil {
		return
	}

	u, err := url.Parse(o.URL)
	if err != nil {
		return
	}
	p.active.Navigate(u)
}

// DestroyOp destroys the webview.
type DestroyOp struct{}

var poolDestroyOp = plugin.NewOpPool[DestroyOp]()

// Add adds a new DestroyOp to the queue.
func (o DestroyOp) Add(op *op.Ops) {
	poolDestroyOp.WriteOp(op, o)
}

// SetCookieOp sets given cookie in the webview.
type SetCookieOp struct {
	Cookie webview.CookieData
}

var poolSetCookieOp = plugin.NewOpPool[SetCookieOp]()

// Add adds a new SetCookieOp to the queue.
func (o SetCookieOp) Add(op *op.Ops) {
	poolSetCookieOp.WriteOp(op, o)
}

func (o *SetCookieOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	if p.active == nil {
		return
	}

	p.run(func() {
		defer poolSetCookieOp.Release(o)
		err := manager.AddCookie(o.Cookie)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// RemoveCookieOp sets given cookie in the webview.
type RemoveCookieOp struct {
	Cookie webview.CookieData
}

var poolRemoveCookieOp = plugin.NewOpPool[RemoveCookieOp]()

// Add adds a new SetCookieOp to the queue.
func (o RemoveCookieOp) Add(op *op.Ops) {
	poolRemoveCookieOp.WriteOp(op, o)
}

func (o *RemoveCookieOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	if p.active == nil {
		return
	}

	p.run(func() {
		defer poolRemoveCookieOp.Release(o)
		err := manager.RemoveCookie(o.Cookie)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// ListCookieOp lists all cookies in the webview.
// The response in sent via CookiesEvent using the
// provided Tag.
type ListCookieOp struct {
	Tag event.Tag
	// Buffer is the buffer to use for the response,
	// that may prevent allocations.
	Buffer []webview.CookieData
}

// CookiesEvent is the event sent when ListCookieOp is executed.
type CookiesEvent struct {
	Cookies []webview.CookieData
}

// ImplementsEvent the event.Event interface.
func (c CookiesEvent) ImplementsEvent() {}

var poolListCookieOp = plugin.NewOpPool[ListCookieOp]()

// Add adds a new ListCookieOp to the queue.
func (o ListCookieOp) Add(op *op.Ops) {
	poolListCookieOp.WriteOp(op, o)
}

func (o *ListCookieOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolListCookieOp.Release(o)
		evt := CookiesEvent{
			Cookies: o.Buffer,
		}
		err := manager.Cookies(func(c *webview.CookieData) bool {
			evt.Cookies = append(evt.Cookies, *c)
			return true
		})
		p.plugin.SendEvent(o.Tag, evt)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// StorageType is the type of storage.
type StorageType int

const (
	// StorageTypeLocal is the local storage.
	StorageTypeLocal StorageType = iota
	// StorageTypeSession is the session storage.
	StorageTypeSession
)

// SetStorageOp sets given Storage in the webview.
type SetStorageOp struct {
	Local   StorageType
	Content webview.StorageData
}

var poolSetStorageOp = plugin.NewOpPool[SetStorageOp]()

// Add adds a new SetStorageOp to the queue.
func (o SetStorageOp) Add(op *op.Ops) {
	poolSetStorageOp.WriteOp(op, o)
}

func (o *SetStorageOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolSetStorageOp.Release(o)
		var err error
		switch o.Local {
		case StorageTypeLocal:
			err = manager.AddLocalStorage(o.Content)
		case StorageTypeSession:
			err = manager.AddSessionStorage(o.Content)
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// RemoveStorageOp sets given Storage in the webview.
type RemoveStorageOp struct {
	Local   StorageType
	Content webview.StorageData
}

var poolRemoveStorageOp = plugin.NewOpPool[RemoveStorageOp]()

// Add adds a new RemoveStorageOp to the queue.
func (o RemoveStorageOp) Add(op *op.Ops) {
	poolRemoveStorageOp.WriteOp(op, o)
}

func (o *RemoveStorageOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolRemoveStorageOp.Release(o)
		var err error
		switch o.Local {
		case StorageTypeLocal:
			err = manager.RemoveLocalStorage(o.Content)
		case StorageTypeSession:
			err = manager.AddSessionStorage(o.Content)
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// ListStorageOp lists all Storage in the webview.
//
// The response in sent via StorageEvent using the
// provided Tag.
type ListStorageOp struct {
	Tag   event.Tag
	Local StorageType
	// Buffer is the buffer to use for the response,
	// that may prevent allocations.
	Buffer []webview.StorageData
}

// StorageEvent is the event sent when ListStorageOp is executed.
type StorageEvent struct {
	Storage []webview.StorageData
}

// ImplementsEvent the event.Event interface.
func (c StorageEvent) ImplementsEvent() {}

var poolListStorageOp = plugin.NewOpPool[ListStorageOp]()

// Add adds a new ListStorageOp to the queue.
func (o ListStorageOp) Add(op *op.Ops) {
	poolListStorageOp.WriteOp(op, o)
}

func (o *ListStorageOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolListStorageOp.Release(o)
		evt := StorageEvent{
			Storage: o.Buffer,
		}

		fn := manager.LocalStorage
		if o.Local == StorageTypeSession {
			fn = manager.SessionStorage
		}

		err := fn(func(c *webview.StorageData) bool {
			evt.Storage = append(evt.Storage, *c)
			return true
		})

		p.plugin.SendEvent(o.Tag, evt)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

type ClearCacheOp struct{}

var poolClearCacheOp = plugin.NewOpPool[ClearCacheOp]()

func (o ClearCacheOp) Add(op *op.Ops) {
	poolClearCacheOp.WriteOp(op, o)
}

func (o *ClearCacheOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.DataManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolClearCacheOp.Release(o)
		err := manager.ClearAll()
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// ExecuteJavascriptOp executes given JavaScript in the webview.
type ExecuteJavascriptOp struct {
	Script string
}

var poolExecuteJavascriptOp = plugin.NewOpPool[ExecuteJavascriptOp]()

// Add adds a new ExecuteJavascriptOp to the queue.
func (o ExecuteJavascriptOp) Add(op *op.Ops) {
	poolExecuteJavascriptOp.WriteOp(op, o)
}

func (o *ExecuteJavascriptOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.JavascriptManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolExecuteJavascriptOp.Release(o)
		err := manager.RunJavaScript(o.Script)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// InstallJavascriptOp installs given JavaScript in the webview, executing
// it every time the webview loads a new page. The script is executed before
// the page is fully loaded.
type InstallJavascriptOp struct {
	Script string
}

var poolInstallJavascriptOp = plugin.NewOpPool[InstallJavascriptOp]()

// Add adds a new ExecuteJavascriptOp to the queue.
func (o InstallJavascriptOp) Add(op *op.Ops) {
	poolInstallJavascriptOp.WriteOp(op, o)
}

func (o *InstallJavascriptOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.JavascriptManager()
	wvTag := p.activeTag

	p.run(func() {
		defer poolInstallJavascriptOp.Release(o)
		err := manager.InstallJavascript(o.Script, webview.JavascriptOnLoadStart)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// MessageReceiverOp receives a message from the webview,
// and sends it to the provided Tag. The message is sent
// as a string.
//
// You can use this to communicate with the webview, by using:
//
//	window.callback.<name>(<message>);
//
// Consider that <name> is the provided Name of the callback,
// and <message> is the message to send to Tag. The Tag will
// receive the message as a string, with MessageEvent.
//
// For further information, see webview.JavascriptManager.
type MessageReceiverOp struct {
	Tag  event.Tag
	Name string
}

var poolMessageReceiverOp = plugin.NewOpPool[MessageReceiverOp]()

// Add adds a new ExecuteJavascriptOp to the queue.
func (o MessageReceiverOp) Add(op *op.Ops) {
	poolMessageReceiverOp.WriteOp(op, o)
}

func (o *MessageReceiverOp) execute(_ *app.Window, p *webViewPlugin, _ system.FrameEvent) {
	manager := p.active.JavascriptManager()
	wvTag := p.activeTag
	tag := o.Tag

	p.run(func() {
		defer poolMessageReceiverOp.Release(o)
		err := manager.AddCallback(o.Name, func(msg string) {
			p.plugin.SendEvent(tag, MessageEvent{Message: msg})
		})

		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

// MessageEvent is the event sent when receiving a message,
// from previously defined MessageReceiverOp.
type MessageEvent struct {
	Message string
}

// ImplementsEvent the event.Event interface.
func (c MessageEvent) ImplementsEvent() {}

// NavigationEvent is issued when the webview change the URL.
type NavigationEvent webview.NavigationEvent

// ImplementsEvent the event.Event interface.
func (NavigationEvent) ImplementsEvent() {}

// TitleEvent is issued when the webview change the title.
type TitleEvent webview.TitleEvent

// ImplementsEvent the event.Event interface.
func (TitleEvent) ImplementsEvent() {}

// ErrorEvent is issued when the webview encounters an error.
type ErrorEvent struct {
	error
}

// ImplementsEvent implements the event.Event interface.
func (ErrorEvent) ImplementsEvent() {}
