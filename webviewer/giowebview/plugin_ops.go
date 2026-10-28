package giowebview

import (
	"reflect"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/op"
	"gioui.org/op/paint"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
)

var wantOps = []reflect.Type{
	reflect.TypeOf(WebViewOp{}),
}

var externalOpsMap = make(map[event.Tag]paint.EmbedOp)

// WebViewOp shows the webview into the specified area.
// Needs to be used inside a clip.Rect and use paint.PaintOp.
type WebViewOp struct {
	Tag event.Tag
}

func (o WebViewOp) Add(op *op.Ops) {
	if extOp, ok := externalOpsMap[o.Tag]; ok {
		extOp.Add(op)
	} else {
		plugin.WriteOp(op, o)
	}
}

func (o WebViewOp) execute(w *app.Window, p *webViewPlugin) {
	p.tagsMutex.Lock()
	runnerIndex, ok := p.tags[o.Tag]
	p.tagsMutex.Unlock()

	if !ok {
		p.tagsMutex.Lock()
		p.tags[o.Tag] = -1 // Mark as being created
		p.tagsMutex.Unlock()

		go func() {
			wv, err := webview.NewWebView(p.config)
			if err != nil {
				panic(err)
			}
			go eventsListener(wv, w, p, o.Tag)

			p.tagsMutex.Lock()
			runnerIndex = len(p.views)
			p.views = append(p.views, wv)
			p.tags[o.Tag] = runnerIndex
			p.tagsMutex.Unlock()
			externalOpsMap[o.Tag] = plugin.NewEmbedOp(wv.ViewPointer())
			p.plugin.SendEvent(o.Tag, ReadyEvent{})
		}()
	}
}
