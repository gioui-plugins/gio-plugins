package gioplugins

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/layout"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"runtime"
	"sync"
	"unsafe"
)

func NewWindow() *app.Window {
	return new(app.Window)
}

// Hijack is the main event handler for the plugin, you MUST use it as a wrapper,
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
//		   e := gioplugins.Hook(w)
//		   // ...
//		}
func Hijack(w *app.Window) event.Event {
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

// Event returns custom events from the last frame.
func Event(gtx layout.Context, filters ...event.Filter) (evt event.Event, ok bool) {
	ptr := unsafe.Pointer(&filters)
	defer runtime.KeepAlive(ptr)

	return _event(gtx, uintptr(ptr))
}

func _event(gtx layout.Context, fptr uintptr) (evt event.Event, ok bool) {
	filters := *(*[]event.Filter)(unsafe.Pointer(fptr))
	source := (*gioInputSource)(unsafe.Pointer(&gtx.Source))

	if li := getInstanceByRouter(source.r); li != nil {
		if evt, ok := li.Plugin.Event(filters...); ok {
			return evt, true
		}
	}
	return nil, false
}

// Execute executes the command.
func Execute(gtx layout.Context, c input.Command) {
	ptr := unsafe.Pointer(&c)
	defer runtime.KeepAlive(ptr)

	_execute(gtx, uintptr(ptr))
}

func _execute(gtx layout.Context, fptr uintptr) {
	cmd := *(*input.Command)(unsafe.Pointer(fptr))
	source := (*gioInputSource)(unsafe.Pointer(&gtx.Source))

	if li := getInstanceByRouter(source.r); li != nil {
		if li.Plugin.Execute(cmd) {
			return
		}
	}
}

// gioInputSource must match the input.Source.
type gioInputSource struct {
	r *input.Router
}

func init() {
	if unsafe.Sizeof(gioInputSource{}) != unsafe.Sizeof(input.Source{}) {
		panic("Gio version not supported")
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
