// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

package explorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Appkit/AppKit.h>

// Defined on explorer_macos.m file.
extern void gioplugins_explorer_saveFile(CFTypeRef viewRef, char * name, uintptr_t id);
extern void gioplugins_explorer_openFile(CFTypeRef viewRef, char * ext, uintptr_t id);
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

	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type driver struct {
	config Config
}

func attachDriver(house *Explorer, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.config = config
}

func (e *driver) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	res := make(chan result)
	callback := func(r result) {
		res <- r
	}
	hcallback := cgo.NewHandle(callback)
	defer hcallback.Delete()

	go e.config.RunOnMain(func() {
		C.gioplugins_explorer_saveFile(C.CFTypeRef(e.config.View), cname, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	runtime.KeepAlive(callback)
	return r.file.(io.WriteCloser), r.error

}

func (e *driver) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
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

	go e.config.RunOnMain(func() {
		C.gioplugins_explorer_openFile(C.CFTypeRef(e.config.View), cextensions, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	return r.file.(io.ReadCloser), r.error
}

//export gioplugins_explorer_importCallback
func gioplugins_explorer_importCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result)); ok {
		v(newOSFile(u, func(s string) (any, error) {
			return os.Open(s)
		}))
	}
}

//export gioplugins_explorer_exportCallback
func gioplugins_explorer_exportCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result)); ok {
		v(newOSFile(u, func(s string) (any, error) {
			return os.Create(s)
		}))
	}
}

func newOSFile(u *C.char, action func(s string) (any, error)) result {
	name := C.GoString(u)
	if name == "" {
		return result{error: ErrUserDecline}
	}

	uri, err := url.Parse(name)
	if err != nil {
		return result{error: err}
	}

	path, err := url.PathUnescape(uri.Path)
	if err != nil {
		return result{error: err}
	}

	f, err := action(path)
	return result{error: err, file: f}
}
