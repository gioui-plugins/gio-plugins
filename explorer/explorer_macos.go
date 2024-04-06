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

	res := make(chan result[io.WriteCloser])
	callback := func(r result[io.WriteCloser]) {
		res <- r
	}
	hcallback := cgo.NewHandle(callback)
	defer hcallback.Delete()

	go e.config.RunOnMain(func() {
		C.saveFile(e.config.View, cname, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	runtime.KeepAlive(callback)
	return r.file.(io.WriteCloser), r.error
}

func (e *driver) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	res := make(chan result[io.ReadCloser])
	callback := func(r result[io.ReadCloser]) {
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
		C.openFile(e.config.View, cextensions, C.uintptr_t(hcallback))
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	return r.file.(io.ReadCloser), r.error
}

//export importCallback
func importCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result[io.ReadCloser])); ok {
		v(newOSFile[io.ReadCloser](u, func(s string) (io.ReadCloser, error) {
			return os.Open(s)
		}))
	}
}

//export exportCallback
func exportCallback(u *C.char, id uintptr) {
	if v, ok := cgo.Handle(id).Value().(func(result[io.WriteCloser])); ok {
		v(newOSFile[io.WriteCloser](u, func(s string) (io.WriteCloser, error) {
			return os.Create(s)
		}))
	}
}

func newOSFile[T io.ReadCloser | io.ReadCloser](u *C.char, action func(s string) (T, error)) result[T] {
	name := C.GoString(u)
	if name == "" {
		return result[T]{error: ErrUserDecline}
	}

	uri, err := url.Parse(name)
	if err != nil {
		return result[T]{error: err}
	}

	path, err := url.PathUnescape(uri.Path)
	if err != nil {
		return result[T]{error: err}
	}

	f, err := action(path)
	return result[T]{error: err, file: f}
}
