package webview

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

extern void getCookies(CFTypeRef config, uintptr_t handler, uintptr_t done);
extern void addCookie(CFTypeRef config, uintptr_t done, char *name, char *value, char *domain, char *path, int64_t expires, uint64_t features);
extern void removeCookie(CFTypeRef config, uintptr_t done, char *name, char *domain, char *path);

*/
import "C"
import (
	"time"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

type cookieManager struct {
	*webview
}

func newCookieManager(w *webview) *cookieManager {
	r := &cookieManager{webview: w}
	return r
}

// Cookies implements the CookieManager interface.
func (s *cookieManager) Cookies(fn DataLooper[CookieData]) (err error) {
	done := make(chan error)

	fr, dr := internal.NewHandle(fn), internal.NewHandle(done)
	defer fr.Delete()
	defer dr.Delete()

	s.scheduler.MustRun(func() {
		C.getCookies(s.driver.webviewConfig, C.uintptr_t(fr), C.uintptr_t(dr))
	})

	return <-done
}

//export getCookiesCallback
func getCookiesCallback(handler uintptr, features uint64, name, value, domain, path *C.char, expires int64) bool {
	return internal.Handle(handler).Value().(DataLooper[CookieData])(&CookieData{
		Name:     C.GoString(name),
		Value:    C.GoString(value),
		Domain:   C.GoString(domain),
		Path:     C.GoString(path),
		Expires:  time.Unix(expires, 0),
		Features: CookieFeatures(features),
	})
}

// AddCookie implements the CookieManager interface.
func (s *cookieManager) AddCookie(c CookieData) error {
	done := make(chan error)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	name, value, domain, path := C.CString(c.Name), C.CString(c.Value), C.CString(c.Domain), C.CString(c.Path)
	var expires C.int64_t
	if c.Expires.After(time.Now()) {
		expires = C.int64_t(c.Expires.Unix())
	}

	defer C.free(unsafe.Pointer(name))
	defer C.free(unsafe.Pointer(value))
	defer C.free(unsafe.Pointer(domain))
	defer C.free(unsafe.Pointer(path))

	s.webview.scheduler.MustRun(func() {
		C.addCookie(s.driver.webviewConfig, C.uintptr_t(dr), name, value, domain, path, expires, C.uint64_t(c.Features))
	})
	return <-done
}

// RemoveCookie implements the CookieManager interface.
func (s *cookieManager) RemoveCookie(c CookieData) error {
	done := make(chan error)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	name := C.CString(c.Name)
	defer C.free(unsafe.Pointer(name))
	domain := C.CString(c.Domain)
	defer C.free(unsafe.Pointer(domain))
	path := C.CString(c.Path)
	defer C.free(unsafe.Pointer(path))

	s.webview.scheduler.MustRun(func() {
		C.removeCookie(s.driver.webviewConfig, C.uintptr_t(dr), name, domain, path)
	})

	return <-done
}
