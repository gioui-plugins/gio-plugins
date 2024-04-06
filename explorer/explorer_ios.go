//go:build ios
// +build ios

package explorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fno-objc-arc

#include <UIKit/UIKit.h>
#include <stdint.h>

// Defined on explorer_ios.m file (implements UIDocumentPickerDelegate).
@interface explorer_picker:NSObject<UIDocumentPickerDelegate>
@property (strong) UIDocumentPickerViewController * picker;
@property (strong) UIViewController * controller;
@property uint64_t mode;
@property uintptr_t callback;
@end

extern CFTypeRef saveFile(CFTypeRef view, char * name, uintptr_t callback, CFTypeRef pooled);
extern CFTypeRef openFile(CFTypeRef view, char * ext, uintptr_t callback, CFTypeRef pooled);
*/
import "C"

import (
	"io"
	"os"
	"path/filepath"
	"runtime/cgo"
	"strings"
	"sync"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	picker C.CFTypeRef

	savePool *sync.Pool
	openPool *sync.Pool
}

func attachDriver(house *Explorer, config Config) {
	house.driver = driver{}
	house.driver.savePool = &sync.Pool{New: func() any { return C.CFTypeRef(0) }}
	house.driver.openPool = &sync.Pool{New: func() any { return C.CFTypeRef(0) }}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	driver.config = config
}

func (e *driver) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	if e.picker == 0 {
		return nil, ErrNotAvailable
	}

	name = filepath.Join(os.TempDir(), name)

	f, err := os.Create(name)
	if err != nil {
		return nil, nil
	}
	if err := f.Close(); err != nil {
		return nil, err
	}

	name = "file://" + name
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	res := make(chan result[io.ReadWriteCloser], 1)
	hcallback := cgo.NewHandle(func(r result[io.ReadWriteCloser]) { res <- r })
	defer hcallback.Delete()

	r := e.savePool.Get().(C.CFTypeRef)
	defer func() {
		e.savePool.Put(r)
	}()

	go e.config.RunOnMain(func() {
		if r = C.saveFile(e.picker, cname, C.uintptr_t(hcallback), r); r == 0 {
			res <- result[io.ReadWriteCloser]{error: ErrNotAvailable}
		}
	})

	file := <-res
	if file.error != nil {
		return nil, file.error
	}
	return file.file, nil
}

func (e *driver) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	if e.picker == 0 {
		return nil, ErrNotAvailable
	}

	s := stringBuilderPool.Get().(*strings.Builder)
	for i, ext := range mimes {
		if i > 0 {
			s.WriteRune(',')
		}
		s.WriteString(strings.TrimPrefix(ext.Extension, "."))
	}

	cextensions := C.CString(s.String())
	defer C.free(unsafe.Pointer(cextensions))

	res := make(chan result[io.ReadWriteCloser], 1)
	hcallback := cgo.NewHandle(func(r result[io.ReadWriteCloser]) { res <- r })
	defer hcallback.Delete()

	r := e.openPool.Get().(C.CFTypeRef)
	defer func() {
		e.openPool.Put(r)
	}()

	go e.config.RunOnMain(func() {
		if r = C.openFile(e.picker, cextensions, C.uintptr_t(hcallback), r); r == 0 {
			res <- result[io.ReadWriteCloser]{error: ErrNotAvailable}
		}
	})

	file := <-res
	if file.error != nil {
		return nil, file.error
	}
	return file.file, nil
}

//export pickerCallback
func pickerCallback(u C.CFTypeRef, id C.uintptr_t) {
	if fn, ok := cgo.Handle(id).Value().(func(result[io.ReadWriteCloser])); ok {
		res := result[io.ReadWriteCloser]{error: ErrUserDecline}
		if u != 0 {
			res.file, res.error = newFile(u)
		}
		fn(res)
	}
}
