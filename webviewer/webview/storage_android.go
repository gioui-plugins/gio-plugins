package webview

/*
#cgo CFLAGS: -Werror
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"strconv"
	"time"
	"unsafe"

	"git.wow.st/gmp/jni"
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
	done := make(chan error, 1)
	dr, fr := internal.NewHandle(done), internal.NewHandle(fn)
	defer dr.Delete()
	defer fr.Delete()

	s.scheduler.MustRun(func() {
		s.driver.callArgs("webview_getCookies", "(JJ)V", func(env jni.Env) []jni.Value {
			return []jni.Value{
				jni.Value(int64(fr)),
				jni.Value(int64(dr)),
			}
		})
	})

	return <-done
}

//export Java_com_inkeliz_webview_sys_1android_getCookiesCallback
func Java_com_inkeliz_webview_sys_1android_getCookiesCallback(env *C.JNIEnv, class C.jclass, ptr C.jlong, msg C.jstring) {
	raw := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(msg))
	c := CookieData{}

	last := len(raw) - 1
	start := 0
	ignore := false

	for i := 0; i < len(raw); i++ {
		if !ignore && raw[i] == '=' {
			c.Name = raw[start:i]
			start = i + 1
			ignore = true
		}
		if raw[i] == ';' || last == i {
			if last == i {
				i++
			}
			c.Value = raw[start:i]
			start = i + 2
			ignore = false

			if !internal.Handle(ptr).Value().(DataLooper[CookieData])(&c) {
				break
			}

			c.Value, c.Name = "", ""
		}
	}
}

// AddCookie implements the CookieManager interface.
func (s *cookieManager) AddCookie(c CookieData) error {
	done := make(chan error, 1)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	cookie := c.Name + "=" + c.Value + ";"
	if c.Domain != "" {
		cookie += " Domain=" + c.Domain + ";"
	}
	if c.Path != "" {
		cookie += " Path=" + c.Path + ";"
	}
	if c.Expires.After(time.Now()) {
		cookie += " Max-Age=" + strconv.FormatInt(int64(c.Expires.Sub(time.Now())/time.Second), 10) + ";"
	}
	if c.Features.IsSecure() {
		cookie += " Secure;"
	}
	if c.Features.IsHTTPOnly() {
		cookie += " HttpOnly;"
	}

	s.scheduler.MustRun(func() {
		s.driver.callArgs("webview_addCookie", "(Ljava/lang/String;Ljava/lang/String;J)V", func(env jni.Env) []jni.Value {
			return []jni.Value{
				jni.Value(jni.JavaString(env, c.Domain)),
				jni.Value(jni.JavaString(env, cookie)),
				jni.Value(int64(dr)),
			}
		})
	})

	return nil
}

// RemoveCookie implements the CookieManager interface.
func (s *cookieManager) RemoveCookie(c CookieData) error {
	return ErrNotSupported
}
