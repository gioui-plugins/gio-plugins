//go:build ios || darwin

package webview

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

extern CFTypeRef config();
extern CFTypeRef create(CFTypeRef config, uintptr_t handler);
extern void resize(CFTypeRef web, CFTypeRef windowRef, float x, float y, float w, float h);
extern void run(CFTypeRef web, CFTypeRef windowRef);
extern void seturl(CFTypeRef web, char *u);
extern void hide(CFTypeRef web);
extern void show(CFTypeRef web);

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

func (r *driver) attach(w *webview) (err error) {
	r.config.RunOnMain(func() {
		w.mutex.Lock()
		defer w.mutex.Unlock()

		r.webviewConfig = C.config()
		r.setDebug()
		r.webviewObject = C.create(r.webviewConfig, C.uintptr_t(uintptr(w.handle)))

		w.javascriptManager = newJavascriptManager(w)
		w.dataManager = newDataManager(w)

		C.run(r.webviewObject, C.CFTypeRef(r.config.View))
	})

	w.mutex.Lock()
	config := r.config
	w.mutex.Unlock()
	w.scheduler.SetRunner(config.RunOnMain)

	return nil
}

func (r *driver) resize(w *webview, pos [4]float32) {
	if pos[2] == 0 && pos[3] == 0 {
		if w.visible {
			C.hide(r.webviewObject)
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

		C.resize(
			r.webviewObject,
			C.CFTypeRef(r.config.View),
			C.float(posFloat[0]),
			C.float(posFloat[1]),
			C.float(posFloat[2]),
			C.float(posFloat[3]),
		)
		if !w.visible {
			C.show(r.webviewObject)
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

	C.seturl(r.webviewObject, u)
}

func (r *driver) close(w *webview) {
	C.webview_cf_release(r.webviewObject)
	C.webview_cf_release(r.webviewConfig)
	return
}
