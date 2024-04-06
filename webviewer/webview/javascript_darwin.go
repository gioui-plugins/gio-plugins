package webview

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

void addCallbackJavascript(CFTypeRef config, char *name, uintptr_t handler);
void runJavascript(CFTypeRef web, char *js, uintptr_t done);
void installJavascript(CFTypeRef config, char *js, uint64_t when);

*/
import "C"

import (
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

type javascriptManager struct {
	webview   *webview
	jsHandler internal.Handle
	callbacks sync.Map // map[string]func(message string)
}

func newJavascriptManager(w *webview) *javascriptManager {
	r := &javascriptManager{webview: w}
	r.installJavascript(fmt.Sprintf(scriptCallback, `window.webkit.messageHandlers.callback.postMessage`), JavascriptOnLoadStart)

	r.jsHandler = internal.NewHandle(r)

	name := C.CString("callback")
	defer C.free(unsafe.Pointer(name))

	C.addCallbackJavascript(r.webview.driver.webviewConfig, name, C.uintptr_t(r.jsHandler))

	return r
}

//export javascriptManagerCallback
func javascriptManagerCallback(handler uintptr, input *C.char) {
	receiveCallback(handler, C.GoString(input))
}

// RunJavaScript implements the JavascriptManager interface.
func (j *javascriptManager) RunJavaScript(js string) error {
	done := make(chan error, 1)

	dr := internal.NewHandle(done)
	defer dr.Delete()

	j.webview.scheduler.MustRun(func() {
		code := C.CString(js)
		defer C.free(unsafe.Pointer(code))

		C.runJavascript(j.webview.driver.webviewObject, code, C.uintptr_t(dr))
	})

	return <-done
}

// InstallJavascript implements the JavascriptManager interface.
func (j *javascriptManager) InstallJavascript(js string, when JavascriptInstallationTime) error {
	j.webview.scheduler.MustRun(func() {
		j.installJavascript(js, when)
	})
	return nil
}

func (j *javascriptManager) installJavascript(js string, when JavascriptInstallationTime) {
	code := C.CString(js)
	defer C.free(unsafe.Pointer(code))

	C.installJavascript(j.webview.driver.webviewConfig, code, C.uint64_t(when))
}

// AddCallback implements the JavascriptManager interface.
func (j *javascriptManager) AddCallback(name string, fn func(message string)) error {
	if len(name) > 255 {
		return ErrJavascriptCallbackInvalidName
	}
	if strings.Contains(name, ".") || strings.Contains(name, " ") {
		return ErrJavascriptCallbackInvalidName
	}
	if _, ok := j.callbacks.Load(name); ok {
		return ErrJavascriptCallbackDuplicate
	}

	j.callbacks.Store(name, fn)
	return nil
}
