package gioplugins

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"sync"
	"unsafe"
)

func NewWindow() *app.Window {
	return new(app.Window)
}

// Event is the main event handler for the plugin, you MUST use it as a wrapper,
// for the app.Window.Event method.
//
// It will COMBINE the events from the app.Window.Event method with the plugin events.
//
// Instead of:
//
//	 w := app.Window{}
//
//		for {
//		    e := w.Event()
//		    // ...
//		}
//
// You should use:
//
//	 w := app.Window{}
//
//		for {
//		   e := gioplugins.Event(w)
//		   // ...
//		}
func Event(w *app.Window) event.Event {
	instance := getInstanceByWindow(w)
	if instance == nil {
		instance = createInstance(w)
	}

	evt := instance.Plugin.ProcessEventFromGio(w.Event())
	if e, ok := evt.(app.FrameEvent); ok {
		updateInstance(instance, e.Source)
	}

	return evt
}

// gioInputSource must match the input.Source.
type gioInputSource struct {
	r *input.Router
}

func init() {
	if unsafe.Sizeof(gioInputSource{}) != unsafe.Sizeof(input.Source{}) {
		panic("Gio version not supported")
	}

	input.SourceEventProcessor = func(r *input.Router, filters ...event.Filter) (evt event.Event, ok bool) {
		if li := getInstanceByRouter(r); li != nil {
			if evt, ok := li.Plugin.Event(filters...); ok {
				return evt, true
			}
		}
		return r.Event(filters...)
	}

	input.SourceExecuteProcessor = func(r *input.Router, c input.Command) {
		if li := getInstanceByRouter(r); li != nil {
			if li.Plugin.Execute(c) {
				return
			}
		}
		r.Execute(c)
	}
}

type instance struct {
	Window *app.Window
	Router *input.Router
	Plugin *plugin.Plugin
}

var handlers = new(sync.Map)         // map[app.Window]*instance
var handlersByRouter = new(sync.Map) // map[*input.Router]*instance

func createInstance(w *app.Window) *instance {
	instance := &instance{
		Window: w,
		Plugin: plugin.NewPlugin(w),
	}

	handlers.Store(w, instance)
	return instance
}

func updateInstance(instance *instance, r input.Source) {
	rr := (*gioInputSource)(unsafe.Pointer(&r))
	if instance.Router != rr.r {
		handlersByRouter.Delete(unsafe.Pointer(instance.Router))
		instance.Router = rr.r
		if rr.r != nil {
			handlersByRouter.Store(unsafe.Pointer(rr.r), instance)
		}
	}
}

func getInstanceByRouter(r *input.Router) *instance {
	li, ok := handlersByRouter.Load(unsafe.Pointer(r))
	if !ok {
		return nil
	}

	return li.(*instance)
}

func getInstanceByWindow(w *app.Window) *instance {
	li, ok := handlers.Load(w)
	if !ok {
		return nil
	}

	return li.(*instance)
}
