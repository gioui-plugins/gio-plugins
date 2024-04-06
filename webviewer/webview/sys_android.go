package webview

/*
#cgo CFLAGS: -Werror
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"

import (
	"errors"
	"unsafe"

	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

//export Java_com_inkeliz_webview_sys_1android_reportDone
func Java_com_inkeliz_webview_sys_1android_reportDone(env *C.JNIEnv, class C.jclass, done C.jlong, v C.jstring) {
	if done == 0 {
		return
	}

	var r error
	if err := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(v)); err != "" {
		r = errors.New(err)
	}

	go func() {
		internal.Handle(uintptr(done)).Value().(chan error) <- r
	}()
}

//export Java_com_inkeliz_webview_sys_1android_reportLoadStatus
func Java_com_inkeliz_webview_sys_1android_reportLoadStatus(env *C.JNIEnv, class C.jclass, handler C.jlong, v C.jstring) {
	url := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(v))
	internal.Handle(uintptr(handler)).Value().(*webview).fan.Send(NavigationEvent{
		URL: url,
	})
}

//export Java_com_inkeliz_webview_sys_1android_reportTitleStatus
func Java_com_inkeliz_webview_sys_1android_reportTitleStatus(env *C.JNIEnv, class C.jclass, handler C.jlong, v C.jstring) {
	title := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(v))
	internal.Handle(uintptr(handler)).Value().(*webview).fan.Send(TitleEvent{
		Title: title,
	})
}
