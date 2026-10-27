package webview

import (
	"net/url"
	"os"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	_user32            = windows.NewLazySystemDLL("user32.dll")
	_SetWindowPos      = _user32.NewProc("SetWindowPos")
	_CreateWindowExW   = _user32.NewProc("CreateWindowExW")
	_DefWindowProcW    = _user32.NewProc("DefWindowProcW")
	_RegisterClassExW  = _user32.NewProc("RegisterClassExW")
	_DestroyWindow     = _user32.NewProc("DestroyWindow")
	_ShowWindow        = _user32.NewProc("ShowWindow")
	_GetClientRect     = _user32.NewProc("GetClientRect")
	_SetWindowLongPtrW = _user32.NewProc("SetWindowLongPtrW")
	_CallWindowProcW   = _user32.NewProc("CallWindowProcW")

	_kernel32         = windows.NewLazySystemDLL("kernel32.dll")
	_GetModuleHandleW = _kernel32.NewProc("GetModuleHandleW")

	// Window styles
	_WS_CHILD        = uintptr(0x40000000)
	_WS_CLIPSIBLINGS = uintptr(0x04000000)
	_WS_VISIBLE      = uintptr(0x10000000)

	// SetWindowPos flags
	_SWP_NOACTIVATE = uintptr(0x0010)

	// Window position constants
	_HWND_BOTTOM = uintptr(1)

	// ShowWindow commands
	_SW_HIDE = uintptr(0)
	_SW_SHOW = uintptr(5)

	// Window messages
	_WM_WINDOWPOSCHANGED = uintptr(0x0047)
	_WM_SIZE             = uintptr(0x0005)
	_WM_MOVE             = uintptr(0x0003)

	_GWLP_WNDPROC = ^uintptr(3 - 1)
)

// wndClassExW is the WNDCLASSEXW structure
type wndClassExW struct {
	cbSize        uint32
	style         uint32
	lpfnWndProc   uintptr
	cbClsExtra    int32
	cbWndExtra    int32
	hInstance     windows.Handle
	hIcon         windows.Handle
	hCursor       windows.Handle
	hbrBackground windows.Handle
	lpszMenuName  *uint16
	lpszClassName *uint16
	hIconSm       windows.Handle
}

func init() {
	for _, s := range []string{"WEBVIEW2_BROWSER_EXECUTABLE_FOLDER", "WEBVIEW2_USER_DATA_FOLDER", "WEBVIEW2_ADDITIONAL_BROWSER_ARGUMENTS", "WEBVIEW2_RELEASE_CHANNEL_PREFERENCE"} {
		_ = os.Setenv(s, "")
	}
}

type driver struct {
	config Config
	active uint32
	dir    []uint16

	controllerCompletedHandler  *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler
	environmentCompletedHandler *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler
	controller                  *_ICoreWebView2Controller
	webview2                    *_ICoreWebView2
	webview22                   *_ICoreWebView22
	webview213                  *_ICoreWebView213

	pinner runtime.Pinner

	callbackTitle *_ICoreWebView2DocumentTitleChangedEventHandler
	callbackLoad  *_ICoreWebView2SourceChangedEventHandler

	// hwnd is the child window that hosts the webview
	hwnd uintptr

	view *webview

	lastNativePos [4]float32
	parentWndProc uintptr
}

func (r *driver) attach(w *webview) error {
	cerr := make(chan error, 1)
	r.view = w

	if err := r.checkInstall(); err != nil {
		return err
	}

	// [Windows] Certs and Proxies must be defined before the initialization
	r.setCerts()
	r.setProxy()
	r.setDir()

	w.mutex.Lock()
	config := r.config
	w.mutex.Unlock()

	go config.RunOnMain(func() {
		if err := windows.CoInitializeEx(0, 0x2); err != nil {
			cerr <- err
			return
		}

		hwndParent := config.HWND
		hModule, _, _ := _GetModuleHandleW.Call(0)

		className, _ := windows.UTF16PtrFromString("GioWebViewHost")
		wc := &wndClassExW{
			cbSize:        uint32(unsafe.Sizeof(wndClassExW{})),
			hInstance:     windows.Handle(hModule),
			lpfnWndProc:   hostWindowProc,
			lpszClassName: className,
		}
		_RegisterClassExW.Call(uintptr(unsafe.Pointer(wc)))

		hwnd, _, _ := _CreateWindowExW.Call(
			0,
			uintptr(unsafe.Pointer(className)),
			0,
			_WS_CHILD|_WS_CLIPSIBLINGS|_WS_VISIBLE,
			0, 0, 100, 100,
			hwndParent,
			0,
			hModule,
			0,
		)
		r.hwnd = hwnd
		windowDriverMu.Lock()
		windowDriverRefs[hwnd] = r
		windowDriverMu.Unlock()

		oldProc, _, _ := _SetWindowLongPtrW.Call(hwndParent, _GWLP_WNDPROC, parentWindowProc)
		r.parentWndProc = oldProc
		parentDriverMu.Lock()
		parentDriverRefs[hwndParent] = r
		parentDriverMu.Unlock()

		_SetWindowPos.Call(r.hwnd, _HWND_BOTTOM, 0, 0, 100, 100, _SWP_NOACTIVATE)
		r.lastNativePos = [4]float32{0, 0, 100, 100}

		r.callbackLoad = &_ICoreWebView2SourceChangedEventHandler{
			VTBL: _CoreWebView2SourceChangedEventHandlerVTBL,
			Invoke: func(this *_ICoreWebView2SourceChangedEventHandler, wv *_ICoreWebView2, v uintptr) uintptr {
				var source *uint16
				syscall.SyscallN(wv.VTBL.GetSource, uintptr(unsafe.Pointer(wv)), uintptr(unsafe.Pointer(&source)))

				w.fan.Send(NavigationEvent{URL: windows.UTF16PtrToString(source)})
				return 0
			},
		}
		r.pinner.Pin(r.callbackLoad)
		r.pinner.Pin(r.callbackLoad.VTBL)
		r.pinner.Pin(&r.callbackLoad.Invoke)

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
		r.pinner.Pin(r.callbackTitle)
		r.pinner.Pin(r.callbackTitle.VTBL)
		r.pinner.Pin(&r.callbackTitle.Invoke)

		r.controllerCompletedHandler = &_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{
			VTBL: _CoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL,
			Invoke: func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, err uintptr, val *_ICoreWebView2Controller) uintptr {
				r.controller = val

				syscall.SyscallN(val.VTBL._IUnknownVTBL.Add, uintptr(unsafe.Pointer(r.controller)))
				syscall.SyscallN(val.VTBL.GetCoreWebView2, uintptr(unsafe.Pointer(val)), uintptr(unsafe.Pointer(&r.webview2)))
				syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Add, uintptr(unsafe.Pointer(r.webview2)))

				syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Query, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(&_GUIDCoreWebView22)), uintptr(unsafe.Pointer(&r.webview22)))
				syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Query, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(&_GUIDCoreWebView213)), uintptr(unsafe.Pointer(&r.webview213)))

				// [Windows] Hook the events
				syscall.SyscallN(r.webview2.VTBL.AddSourceChanged, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(r.callbackLoad)), uintptr(unsafe.Pointer(&r.callbackLoad.Token)))
				syscall.SyscallN(r.webview2.VTBL.AddDocumentTitleChanged, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(r.callbackTitle)), uintptr(unsafe.Pointer(&r.callbackTitle.Token)))

				atomic.AddUint32(&r.active, 1)

				return 0
			},
		}
		r.pinner.Pin(r.controllerCompletedHandler)
		r.pinner.Pin(r.controllerCompletedHandler.VTBL)
		r.pinner.Pin(&r.controllerCompletedHandler.Invoke)

		r.environmentCompletedHandler = &_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{
			VTBL: _CoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL,
			Invoke: func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, err uintptr, val *_ICoreWebView2Environment) uintptr {

				syscall.SyscallN(val.VTBL._IUnknownVTBL.Add, uintptr(unsafe.Pointer(val)))

				// Create the controller using our child window as parent
				syscall.SyscallN(val.VTBL.CreateCoreWebView2Controller, uintptr(unsafe.Pointer(val)), r.hwnd, uintptr(unsafe.Pointer(r.controllerCompletedHandler)))

				return 0
			},
		}
		r.pinner.Pin(r.environmentCompletedHandler)
		r.pinner.Pin(r.environmentCompletedHandler.VTBL)
		r.pinner.Pin(&r.environmentCompletedHandler.Invoke)

		pathPtr := uintptr(0)
		if len(r.dir) > 0 {
			pathPtr = uintptr(unsafe.Pointer(&r.dir[0]))
		}

		hr, _, _ := _CreateCoreWebView2EnvironmentWithOptions.Call(0, pathPtr, 0, uintptr(unsafe.Pointer(r.environmentCompletedHandler)))
		if hr != 0 {
			cerr <- ErrInvalidOptionChange
			return
		}

		cerr <- nil
	})

	finalErr := make(chan error, 1)
	go func() {
		if err := <-cerr; err != nil {
			finalErr <- err
			return
		}

		for atomic.LoadUint32(&r.active) == 0 {
		}

		w.scheduler.SetRunner(config.RunOnMain)
		finalErr <- nil
	}()

	err := <-finalErr

	w.javascriptManager = newJavascriptManager(w)
	w.dataManager = newDataManager(w)

	return err
}

func (r *driver) configure(w *webview, config Config) {
	if atomic.LoadUint32(&r.active) == 1 {
		w.scheduler.SetRunner(config.RunOnMain)
	}
	r.config = config
}

func (r *driver) resize(w *webview, pos [4]float32) {
	if r.controller == nil || r.hwnd == 0 {
		return
	}

	if pos[2] == 0 && pos[3] == 0 {
		r.lastNativePos = pos
		syscall.SyscallN(r.controller.VTBL.PutIsVisible, uintptr(unsafe.Pointer(r.controller)), 0)
		_ShowWindow.Call(r.hwnd, _SW_HIDE)
	} else {
		left := int32(pos[0] + 0.5)
		top := int32(pos[1] + 0.5)
		right := int32(pos[0]+0.5) + int32(pos[2]+0.5)
		bottom := int32(pos[1]+0.5) + int32(pos[3]+0.5)
		width := right - left
		height := bottom - top

		// Move and show, keeping at bottom
		_SetWindowPos.Call(r.hwnd, _HWND_BOTTOM, uintptr(left), uintptr(top), uintptr(width), uintptr(height), _SWP_NOACTIVATE)
		_ShowWindow.Call(r.hwnd, _SW_SHOW)

		bounds := [4]int32{0, 0, width, height}
		r.lastNativePos = pos
		syscall.SyscallN(r.controller.VTBL.PutIsVisible, uintptr(unsafe.Pointer(r.controller)), 1)
		syscall.SyscallN(r.controller.VTBL.PutBounds, uintptr(unsafe.Pointer(r.controller)), uintptr(unsafe.Pointer(&bounds)))
	}
}

func (r *driver) navigate(w *webview, url *url.URL) {
	syscall.SyscallN(r.webview2.VTBL.Navigate, uintptr(unsafe.Pointer(r.webview2)), uintptr(unsafe.Pointer(windows.StringToUTF16Ptr(url.String()))), 0)
}

func (r *driver) close(w *webview) {
	if r.webview2 != nil {
		syscall.SyscallN(r.webview2.VTBL._IUnknownVTBL.Release, uintptr(unsafe.Pointer(r.webview2)))
	}
	if r.controller != nil {
		syscall.SyscallN(r.controller.VTBL._IUnknownVTBL.Release, uintptr(unsafe.Pointer(r.controller)))
	}
	if r.config.HWND != 0 && r.parentWndProc != 0 {
		parentDriverMu.Lock()
		delete(parentDriverRefs, r.config.HWND)
		parentDriverMu.Unlock()
		_SetWindowLongPtrW.Call(r.config.HWND, _GWLP_WNDPROC, r.parentWndProc)
		r.parentWndProc = 0
	}
	if r.hwnd != 0 {
		windowDriverMu.Lock()
		delete(windowDriverRefs, r.hwnd)
		windowDriverMu.Unlock()
		_DestroyWindow.Call(r.hwnd)
		r.hwnd = 0
	}
}

func (r *driver) viewPointer() uintptr {
	return r.hwnd
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

type point struct {
	x int32
	y int32
}

type rect struct {
	left   int32
	top    int32
	right  int32
	bottom int32
}

func (r rect) pos() [4]float32 {
	return [4]float32{
		float32(r.left),
		float32(r.top),
		float32(r.right - r.left),
		float32(r.bottom - r.top),
	}
}

var (
	windowDriverMu   sync.RWMutex
	windowDriverRefs = map[uintptr]*driver{}
	parentDriverMu   sync.RWMutex
	parentDriverRefs = map[uintptr]*driver{}
)

var hostWindowProc = windows.NewCallback(func(hwnd, msg, wParam, lParam uintptr) uintptr {
	windowDriverMu.RLock()
	r := windowDriverRefs[hwnd]
	windowDriverMu.RUnlock()
	if r != nil {
		switch msg {
		case _WM_WINDOWPOSCHANGED, _WM_SIZE, _WM_MOVE:
			r.syncResizeFromHWND()
		}
	}

	ret, _, _ := _DefWindowProcW.Call(hwnd, msg, wParam, lParam)
	return ret
})

var parentWindowProc = windows.NewCallback(func(hwnd, msg, wParam, lParam uintptr) uintptr {
	parentDriverMu.RLock()
	r := parentDriverRefs[hwnd]
	parentDriverMu.RUnlock()

	if r == nil || r.parentWndProc == 0 {
		ret, _, _ := _DefWindowProcW.Call(hwnd, msg, wParam, lParam)
		return ret
	}

	ret, _, _ := _CallWindowProcW.Call(r.parentWndProc, hwnd, msg, wParam, lParam)

	switch msg {
	case _WM_WINDOWPOSCHANGED, _WM_SIZE, _WM_MOVE:
		r.syncResizeFromParent()
	}

	return ret
})

func (r *driver) syncResizeFromParent() {
	if r.view == nil || r.config.RunOnMain == nil || r.hwnd == 0 {
		return
	}

	var rc rect
	_GetClientRect.Call(r.config.HWND, uintptr(unsafe.Pointer(&rc)))
	pos := rc.pos()
	if pos == r.lastNativePos {
		return
	}

	r.lastNativePos = pos
	view := r.view
	go r.config.RunOnMain(func() {
		if view == nil {
			return
		}
		view.Resize(
			Point{X: pos[2], Y: pos[3]},
			Point{X: 0, Y: 0},
		)
	})
}

func (r *driver) syncResizeFromHWND() {
	if r.view == nil || r.config.RunOnMain == nil || r.hwnd == 0 {
		return
	}

	var rc rect
	_GetClientRect.Call(r.hwnd, uintptr(unsafe.Pointer(&rc)))
	pos := [4]float32{r.lastNativePos[0], r.lastNativePos[1], float32(rc.right - rc.left), float32(rc.bottom - rc.top)}
	if pos == r.lastNativePos {
		return
	}

	r.lastNativePos = pos
	view := r.view
	go r.config.RunOnMain(func() {
		if view == nil {
			return
		}
		view.Resize(
			Point{X: pos[2], Y: pos[3]},
			Point{X: pos[0], Y: pos[1]},
		)
	})
}
