//go:build ios || darwin

package webview

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

extern CFTypeRef webview_config();
extern CFTypeRef webview_create(CFTypeRef config, uintptr_t handler);
extern void webview_resize(CFTypeRef web, CFTypeRef windowRef, float x, float y, float w, float h);
extern void webview_run(CFTypeRef web, CFTypeRef windowRef);
extern void webview_seturl(CFTypeRef web, char *u);
extern void webview_hide(CFTypeRef web);
extern void webview_show(CFTypeRef web);

void webview_cf_release(CFTypeRef obj) {
	CFRelease(obj);
}

*/
import "C"
import (
	"math"
	"net/url"
	"unsafe"
)

type driver struct {
	config Config

	webviewObject C.CFTypeRef
	webviewConfig C.CFTypeRef
}

func (r *driver) attach(w *webview) (_ error) {
	cerr := make(chan error, 1)
	go r.config.RunOnMain(func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()

		r.webviewConfig = C.webview_config()
		r.setDebug()
		r.webviewObject = C.webview_create(r.webviewConfig, C.uintptr_t(uintptr(w.handle)))

		w.javascriptManager = newJavascriptManager(w)
		w.dataManager = newDataManager(w)

		C.webview_run(r.webviewObject, C.CFTypeRef(r.config.View))
		cerr <- nil
	})

	err := <-cerr

	w.mutex.Lock()
	config := r.config
	w.mutex.Unlock()
	w.scheduler.SetRunner(config.RunOnMain)

	return err
}

func (r *driver) resize(w *webview, pos [4]float32) {
	if pos[2] == 0 && pos[3] == 0 {
		if w.visible {
			C.webview_hide(r.webviewObject)
			w.visible = false
		}
	} else {
		pp := r.config.PxPerDp
		if pp == 0 {
			pp = 1
		}

		posFloat := [4]float32{
			pos[0] / pp, pos[1] / pp, pos[2] / pp, pos[3] / pp,
		}

		for i, v := range posFloat {
			if math.IsNaN(float64(v)) || math.IsInf(float64(v), 0) {
				posFloat[i] = 0
			}
		}

		C.webview_resize(
			r.webviewObject,
			C.CFTypeRef(r.config.View),
			C.float(posFloat[0]),
			C.float(posFloat[1]),
			C.float(posFloat[2]),
			C.float(posFloat[3]),
		)
		if !w.visible {
			C.webview_show(r.webviewObject)
			w.visible = true
		}
	}
}

func (r *driver) configure(w *webview, config Config) {
	r.config = config
	w.scheduler.SetRunner(w.driver.config.RunOnMain)
}

func (r *driver) navigate(w *webview, url *url.URL) {
	u := C.CString(url.String())
	defer C.free(unsafe.Pointer(u))

	C.webview_seturl(r.webviewObject, u)
}

func (r *driver) close(w *webview) {
	C.webview_cf_release(r.webviewObject)
	C.webview_cf_release(r.webviewConfig)
	return
}

func (r *driver) viewPointer() uintptr {
	return uintptr(unsafe.Pointer(r.webviewObject))
}
