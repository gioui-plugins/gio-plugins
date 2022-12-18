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

	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type explorer struct {
	window *app.Window
	picker C.CFTypeRef

	savePool *sync.Pool
	openPool *sync.Pool
}

func (e *explorerPlugin) listenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		e.picker = C.CFTypeRef(evt.ViewController)
		if e.savePool == nil {
			e.savePool = &sync.Pool{New: func() any { return C.CFTypeRef(0) }}
		}
		if e.openPool == nil {
			e.openPool = &sync.Pool{New: func() any { return C.CFTypeRef(0) }}
		}
	}
}

func (e *explorerPlugin) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	if e.picker == 0 {
		return nil, ErrNotAvailable
	}

	name = filepath.Join(os.TempDir(), name)

	f, err := os.Create(name)
	if err != nil {
		return nil, nil
	}
	f.Close()

	name = "file://" + name
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	res := make(chan result, 1)
	hcallback := cgo.NewHandle(func(r result) { res <- r })
	defer hcallback.Delete()

	r := e.savePool.Get().(C.CFTypeRef)
	defer func() {
		e.savePool.Put(r)
	}()

	go e.window.Run(func() {
		if r = C.saveFile(e.explorer.picker, cname, C.uintptr_t(hcallback), r); r == 0 {
			res <- result{error: ErrNotAvailable}
		}
	})

	file := <-res
	if file.error != nil {
		return nil, file.error
	}
	return file.file.(io.WriteCloser), nil
}

func (e *explorerPlugin) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
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

	res := make(chan result, 1)
	hcallback := cgo.NewHandle(func(r result) { res <- r })
	defer hcallback.Delete()

	r := e.openPool.Get().(C.CFTypeRef)
	defer func() {
		e.openPool.Put(r)
	}()

	go e.window.Run(func() {
		if r = C.openFile(e.explorer.picker, cextensions, C.uintptr_t(hcallback), r); r == 0 {
			res <- result{error: ErrNotAvailable}
		}
	})

	file := <-res
	if file.error != nil {
		return nil, file.error
	}
	return file.file.(io.ReadCloser), nil
}

//export pickerCallback
func pickerCallback(u C.CFTypeRef, id C.uintptr_t) {
	if fn, ok := cgo.Handle(id).Value().(func(result)); ok {
		res := result{error: ErrUserDecline}
		if u != 0 {
			res.file, res.error = newFile(u)
		}
		fn(res)
	}
}
