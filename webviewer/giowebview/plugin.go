package giowebview

import (
	"container/list"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
	"reflect"
	"sync"

	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/io/event"
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

	frame app.FrameEvent

	config    webview.Config
	viewEvent app.ViewEvent
}

// TypeOp implements plugin.Handler.
func (p *webViewPlugin) TypeOp() []reflect.Type { return wantOps }

// TypeCommand implements plugin.Handler.
func (p *webViewPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements plugin.Handler.
func (p *webViewPlugin) TypeEvent() []reflect.Type { return wantEvent }

// Op implements plugin.Handler.
func (p *webViewPlugin) Op(op interface{}) {
	switch v := op.(type) {
	case *WebViewOp:
		defer _WebViewOpPool.Release(v)

		v.execute(p.window, p, p.frame)
	case *OffsetOp:
		defer _OffsetOpPool.Release(v)

		v.execute(p.window, p, p.frame)
	case *RectOp:
		defer _RectOpPool.Release(v)

		v.execute(p.window, p, p.frame)
	}
}

// Execute implements plugin.Handler.
func (p *webViewPlugin) Execute(cmd interface{}) {
	switch v := cmd.(type) {
	case interface {
		execute(w *app.Window, p *webViewPlugin, _ app.FrameEvent)
	}:
		v.execute(p.window, p, p.frame)
	default:
	}
}

// Event implements plugin.Handler.
func (p *webViewPlugin) Event(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		p.viewEvent = evt

	case app.FrameEvent:
		p.frame = evt

		// Reset the seen map.
		for s := range p.seem {
			p.seem[s] = false
		}

	case app.DestroyEvent:
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

func (p *webViewPlugin) getWebView(tag event.Tag) (webview.WebView, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	i, ok := p.tags[tag]
	if !ok {
		return nil, false
	}
	return p.views[i], true
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
