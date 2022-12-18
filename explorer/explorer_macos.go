// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

package explorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Appkit/AppKit.h>

// Defined on explorer_macos.m file.
extern void saveFile(CFTypeRef viewRef, char * name, uintptr_t id);
extern void openFile(CFTypeRef viewRef, char * ext, uintptr_t id);
*/
import "C"
import (
	"io"
	"net/url"
	"os"
	"runtime"
	"runtime/cgo"
	"strings"
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type explorer struct {
	view C.CFTypeRef
}

func (e *explorerPlugin) listenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		e.view = C.CFTypeRef(evt.View)
	}
}

func (e *explorerPlugin) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	res := make(chan result)
	callback := func(r result) {
		res <- r
	}
	hcallback := cgo.NewHandle(callback)
	defer hcallback.Delete()

	go e.window.Run(func() {
		C.saveFile(e.view, cname, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	runtime.KeepAlive(callback)
	return r.file.(io.WriteCloser), r.error

}

func (e *explorerPlugin) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	res := make(chan result)
	callback := func(r result) {
		res <- r
	}
	hcallback := cgo.NewHandle(callback)
	defer hcallback.Delete()

	extensions := make([]string, len(mimes))
	for i, ext := range mimes {
		extensions[i] = strings.TrimPrefix(ext.Extension, ".")
	}

	cextensions := C.CString(strings.Join(extensions, ","))
	defer C.free(unsafe.Pointer(cextensions))

	go e.window.Run(func() {
		C.openFile(e.view, cextensions, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	return r.file.(io.ReadCloser), r.error
}

//export importCallback
func importCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result)); ok {
		v(newOSFile(u, os.Open))
	}
}

//export exportCallback
func exportCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result)); ok {
		v(newOSFile(u, os.Create))
	}
}

func newOSFile(u *C.char, action func(s string) (*os.File, error)) result {
	name := C.GoString(u)
	if name == "" {
		return result{error: ErrUserDecline, file: nil}
	}

	uri, err := url.Parse(name)
	if err != nil {
		return result{error: err, file: nil}
	}

	path, err := url.PathUnescape(uri.Path)
	if err != nil {
		return result{error: err, file: nil}
	}

	f, err := action(path)
	return result{error: err, file: f}
}
