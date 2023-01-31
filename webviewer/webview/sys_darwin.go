package webview

import "C"
import (
	"errors"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

//export reportDone
func reportDone(done uintptr, e *C.char) {
	if err := C.GoString(e); err != "" {
		internal.Handle(done).Value().(chan error) <- errors.New(err)
	} else {
		internal.Handle(done).Value().(chan error) <- nil
	}
}

//export reportLoadStatus
func reportLoadStatus(handler uintptr, url *C.char) {
	internal.Handle(handler).Value().(*webview).fan.Send(NavigationEvent{
		URL: C.GoString(url),
	})
}

//export reportTitleStatus
func reportTitleStatus(handler uintptr, title *C.char) {
	internal.Handle(handler).Value().(*webview).fan.Send(TitleEvent{
		Title: C.GoString(title),
	})
}

/*
func _reportDone(msg string) {
	if len(msg) < 16 {
		return
	}
	ptr, err := strconv.ParseUint(msg, 16, 64)
	if err != nil {
		return
	}
	reportDone(uintptr(ptr), nil)
}
*/
