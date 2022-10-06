package webview

import (
	"math"
	"syscall"
	"time"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
	"golang.org/x/sys/windows"
)

type cookieManager struct {
	*webview
	*_ICoreWebView2CookieManager
}

func newCookieManager(w *webview) *cookieManager {
	r := &cookieManager{webview: w}
	w.scheduler.MustRun(func() {
		syscall.SyscallN(
			r.webview.driver.webview22.VTBL.CookieManager,
			uintptr(unsafe.Pointer(r.webview.driver.webview22)),
			uintptr(unsafe.Pointer(&r._ICoreWebView2CookieManager)),
		)
	})
	return r
}

func (s *cookieManager) Cookies(fn DataLooper[CookieData]) (err error) {
	done := make(chan error)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	handler := &_ICoreWebView2GetCookiesCompletedHandler{
		VTBL: _CoreWebView2GetCookiesCompletedHandlerVTBL,
		Invoke: func(this *_ICoreWebView2GetCookiesCompletedHandler, err uintptr, cookies *_ICoreWebView2CookieList) uintptr {
			var (
				cookie *_ICoreWebView2Cookie
				data   CookieData
				length uint32
			)

			syscall.SyscallN(cookies.VTBL.Count, uintptr(unsafe.Pointer(cookies)), uintptr(unsafe.Pointer(&length)))

			for index := uint32(0); index < length; index++ {
				syscall.SyscallN(
					cookies.VTBL.GetValueAtIndex,
					uintptr(unsafe.Pointer(cookies)),
					uintptr(index),
					uintptr(unsafe.Pointer(&cookie)),
				)

				var name, value, domain, path *uint16
				syscall.SyscallN(cookie.VTBL.GetName, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&name)))
				syscall.SyscallN(cookie.VTBL.GetValue, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&value)))
				syscall.SyscallN(cookie.VTBL.GetDomain, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&domain)))
				syscall.SyscallN(cookie.VTBL.GetPath, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&path)))

				var expires float64
				syscall.SyscallN(cookie.VTBL.GetExpires, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&expires)))

				var secure, httponly uintptr
				syscall.SyscallN(cookie.VTBL.IsSecure, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&secure)))
				syscall.SyscallN(cookie.VTBL.IsHttpOnly, uintptr(unsafe.Pointer(cookie)), uintptr(unsafe.Pointer(&httponly)))

				data.Name = windows.UTF16PtrToString(name)
				data.Value = windows.UTF16PtrToString(value)
				data.Domain = windows.UTF16PtrToString(domain)
				data.Path = windows.UTF16PtrToString(path)
				data.Expires = time.Unix(int64(expires), 0)
				data.Features = 0
				if secure != 0 {
					data.Features |= CookieSecure
				}
				if httponly != 0 {
					data.Features |= CookieHTTPOnly
				}

				next := fn(&data)
				for _, v := range []*uint16{name, value, domain, path} {
					windows.CoTaskMemFree(unsafe.Pointer(v))
				}

				if !next {
					break
				}
			}

			done <- nil
			return 0
		},
	}

	var url uintptr
	s.scheduler.MustRun(func() {
		syscall.SyscallN(
			s.webview.driver.webview2.VTBL.GetSource,
			uintptr(unsafe.Pointer(s.webview.driver.webview2)),
			uintptr(unsafe.Pointer(&url)),
		)

		syscall.SyscallN(
			s._ICoreWebView2CookieManager.VTBL.GetCookies,
			uintptr(unsafe.Pointer(s._ICoreWebView2CookieManager)),
			url,
			uintptr(unsafe.Pointer(handler)),
		)
	})

	return <-done
}

func (s *cookieManager) AddCookie(c CookieData) error {
	s.scheduler.MustRun(func() {
		var cookie *_ICoreWebView2Cookie
		syscall.SyscallN(
			s._ICoreWebView2CookieManager.VTBL.CreateCookie,
			uintptr(unsafe.Pointer(s._ICoreWebView2CookieManager)),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Name))),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Value))),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Domain))),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Path))),
			uintptr(unsafe.Pointer(&cookie)),
		)

		expiresFloat := float64(c.Expires.Unix())
		syscall.SyscallN(cookie.VTBL.PutExpires, uintptr(unsafe.Pointer(cookie)), *(*uintptr)(unsafe.Pointer(&expiresFloat)))
		if c.Features.IsSecure() {
			syscall.SyscallN(cookie.VTBL.PutSecure, uintptr(unsafe.Pointer(cookie)), uintptr(math.MaxInt))
		}
		if c.Features.IsHTTPOnly() {
			syscall.SyscallN(cookie.VTBL.PutHttpOnly, uintptr(unsafe.Pointer(cookie)), uintptr(math.MaxInt))
		}

		syscall.SyscallN(
			s._ICoreWebView2CookieManager.VTBL.AddOrUpdateCookie,
			uintptr(unsafe.Pointer(s._ICoreWebView2CookieManager)),
			uintptr(unsafe.Pointer(cookie)),
		)
	})

	return nil
}

func (s *cookieManager) RemoveCookie(c CookieData) error {
	s.scheduler.MustRun(func() {
		syscall.SyscallN(
			s._ICoreWebView2CookieManager.VTBL.DeleteCookiesWithDomainAndPath,
			uintptr(unsafe.Pointer(s._ICoreWebView2CookieManager)),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Name))),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Domain))),
			uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(c.Path))),
		)
	})

	return nil
}
