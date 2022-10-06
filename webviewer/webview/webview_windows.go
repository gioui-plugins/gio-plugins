package webview

import (
	"net/url"
	"os"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

func init() {
	for _, s := range []string{"WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", "WEBVIEW2_USER_DATA_FOLDER", "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "WEBVIEW2_RELEASE_CHANNEL_PREFERENCE"} {
		os.Setenv(s, "")
	}
}

type driver struct {
	config Config
	active uint32

	controllerCompletedHandler  *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler
	environmentCompletedHandler *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controller                  *_ICoreWebView2Controller
	webview2                    *_ICoreWebView2
	webview22                   *_ICoreWebView22

	callbackTitle *_ICoreWebView2DocumentTitleChangedEventHandler
	callbackLoad  *_ICoreWebView2SourceChangedEventHandler
}

func (r *driver) attach(w *webview) error {
	cerr := make(chan error, 1)

	if err := r.checkInstall(); err != nil {
		return err
	}

	// [Windows] Certs and Proxies must be defined before the initialization
	r.setCerts()
	r.setProxy()

	go r.config.RunOnMain(func() {
		windows.CoInitializeEx(0, 0x2)

		r.callbackLoad = &_ICoreWebView2SourceChangedEventHandler{
			VTBL: _CoreWebView2SourceChangedEventHandlerVTBL,
			Invoke: func(this *_ICoreWebView2SourceChangedEventHandler, wv *_ICoreWebView2, v uintptr) uintptr {
				var url *uint16
				syscall.SyscallN(wv.VTBL.GetSource, uintptr(unsafe.Pointer(wv)), uintptr(unsafe.Pointer(&url)))

				w.fan.Send(NavigationEvent{
					URL: windows.UTF16PtrToString(url),
				})
				return 0
			},
		}

		r.callbackTitle = &_ICoreWebView2DocumentTitleChangedEventHandler{
			VTBL: _CoreWebView2DocumentTitleChangedEventHandlerVTBL,
			Invoke: func(this *_ICoreWebView2DocumentTitleChangedEventHandler, wv *_ICoreWebView2, v uintptr) uintptr {
				var title *uint16
				syscall.SyscallN(wv.VTBL.GetDocumentTitle, uintptr(unsafe.Pointer(wv)), uintptr(unsafe.Pointer(&title)))

				w.fan.Send(TitleEvent{
					Title: windows.UTF16PtrToString(title),
				})
				return 0
			},
		}

		r.controllerCompletedHandler = &_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{
			VTBL: _CoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL,
			Invoke: func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, err uintptr, val *_ICoreWebView2Controller) uintptr {
				r.controller = val
				syscall.SyscallN(
					val.VTBL._IUnknownVTBL.Add,
					uintptr(unsafe.Pointer(r.controller)),
				)
				syscall.SyscallN(
					val.VTBL.GetCoreWebView2,
					uintptr(unsafe.Pointer(val)),
					uintptr(unsafe.Pointer(&r.webview2)),
				)
				syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Add, uintptr(unsafe.Pointer(r.webview2)))

				syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Query, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(&_GUIDCoreWebView22)), uintptr(unsafe.Pointer(&r.webview22)))

				// [Windows] Hook the events
				syscall.SyscallN(r.webview2.VTBL.AddSourceChanged, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(r.callbackLoad)))
				syscall.SyscallN(r.webview2.VTBL.AddDocumentTitleChanged, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(r.callbackTitle)))

				atomic.AddUint32(&r.active, 1)

				return 0
			},
		}

		r.environmentCompletedHandler = &_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
			VTBL: _CoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL,
			Invoke: func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, err uintptr, val *_ICoreWebView2Environment) uintptr {
				syscall.SyscallN(val.VTBL._IUnknownVTBL.Add, uintptr(unsafe.Pointer(val)))

				syscall.SyscallN(val.VTBL.CreateCoreWebView2Controller, uintptr(unsafe.Pointer(val)), r.config.HWND, uintptr(unsafe.Pointer(r.controllerCompletedHandler)))

				return 0
			},
		}

		hr, _, _ := _CreateCoreWebView2EnvironmentWithOptions.Call(0, 0, 0, uint64(uintptr(unsafe.Pointer(r.environmentCompletedHandler))))
		if hr != 0 {
			cerr <- ErrInvalidOptionChange
			return
		}

		cerr <- nil
	})

	if err := <-cerr; err != nil {
		return err
	}

	go func() {
		for atomic.LoadUint32(&r.active) == 0 {
		}
		w.scheduler.SetRunner(r.config.RunOnMain)
	}()

	w.javascriptManager = newJavascriptManager(w)
	w.dataManager = newDataManager(w)

	return nil
}

func (r *driver) configure(w *webview, config Config) {
	r.config = config
	if atomic.LoadUint32(&r.active) == 1 {
		w.scheduler.SetRunner(w.driver.config.RunOnMain)
	}
}

func (r *driver) resize(w *webview, pos [4]float32) {
	if pos[2] == 0 && pos[3] == 0 {
		syscall.SyscallN(
			r.controller.VTBL.PutIsVisible,
			uintptr(unsafe.Pointer(r.controller)),
			0,
		)
	} else {
		pos := [4]int32{int32(pos[0] + 0.5), int32(pos[1] + 0.5), int32(pos[0]+0.5) + int32(pos[2]+0.5), int32(pos[1]+0.5) + int32(pos[3]+0.5)}
		syscall.SyscallN(
			r.controller.VTBL.PutIsVisible,
			uintptr(unsafe.Pointer(r.controller)),
			1,
		)
		syscall.SyscallN(
			r.controller.VTBL.PutBounds,
			uintptr(unsafe.Pointer(r.controller)),
			uintptr(unsafe.Pointer(&pos)),
		)
	}
}

func (r *driver) navigate(w *webview, url *url.URL) {
	syscall.SyscallN(
		r.webview2.VTBL.Navigate,
		uintptr(unsafe.Pointer(r.webview2)),
		uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url.String()))),
		0,
	)
}

func (r *driver) close(w *webview) {
	if r.webview2 != nil {
		syscall.SyscallN(
			r.webview2.VTBL._IUnknownVTBL.Release,
			uintptr(unsafe.Pointer(r.webview2)),
		)
	}
	if r.controller != nil {
		syscall.SyscallN(
			r.controller.VTBL._IUnknownVTBL.Release,
			uintptr(unsafe.Pointer(r.controller)),
		)
	}
}

func (r *driver) checkInstall() error {
	haveInstalled := false

	for _, local := range [...]registry.Key{registry.LOCAL_MACHINE, registry.CURRENT_USER} {
		for _, p := range registryPaths {
			key, err := registry.OpenKey(local, p, registry.QUERY_VALUE)
			if err != nil {
				continue
			}

			version, _, err := key.GetStringValue(`pv`)
			if err != nil || version == "" {
				continue
			}

			haveInstalled = true
			break
		}
	}

	if !haveInstalled {
		return ErrNotInstalled
	}
	return nil
}
