//go:build windows

package webview

import (
	"unsafe"

	"github.com/jchv/go-winloader"
	"golang.org/x/sys/windows"
)

var (
	_CreateCoreWebView2EnvironmentWithOptions winloader.Proc
	_Ole32                                    = windows.NewLazySystemDLL("ole32.dll")
	_Ole32CoTaskMemAlloc                      = _Ole32.NewProc("CoTaskMemAlloc")
)

func init() {
	dll, _ := winloader.LoadFromMemory(dllFile)
	_CreateCoreWebView2EnvironmentWithOptions = dll.Proc("CreateCoreWebView2EnvironmentWithOptions")
}

var (
	// referenceHolder prevents GC from releasing the COM object.
	referenceHolder = make(map[unsafe.Pointer]struct{}, 64)
)

type (
	// _IUnknownVTBL implements IUnknown
	_IUnknownVTBL struct {
		Query   uintptr
		Add     uintptr
		Release uintptr
	}
)

type (
	// _ICoreWebView2 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.622.22
	_ICoreWebView2 struct {
		VTBL *_ICoreWebView2VTBL
	}

	// _ICoreWebView2VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2?view=webview2-1.0.622.22
	_ICoreWebView2VTBL struct {
		_IUnknownVTBL
		GetSettings                            uintptr
		GetSource                              uintptr
		Navigate                               uintptr
		NavigateToString                       uintptr
		AddNavigationStarting                  uintptr
		RemoveNavigationStarting               uintptr
		AddContentLoading                      uintptr
		RemoveContentLoading                   uintptr
		AddSourceChanged                       uintptr
		RemoveSourceChanged                    uintptr
		AddHistoryChanged                      uintptr
		RemoveHistoryChanged                   uintptr
		AddNavigationCompleted                 uintptr
		RemoveNavigationCompleted              uintptr
		AddFrameNavigationStarting             uintptr
		RemoveFrameNavigationStarting          uintptr
		AddFrameNavigationCompleted            uintptr
		RemoveFrameNavigationCompleted         uintptr
		AddScriptDialogOpening                 uintptr
		RemoveScriptDialogOpening              uintptr
		AddPermissionRequested                 uintptr
		RemovePermissionRequested              uintptr
		AddProcessFailed                       uintptr
		RemoveProcessFailed                    uintptr
		AddScriptToExecuteOnDocumentCreated    uintptr
		RemoveScriptToExecuteOnDocumentCreated uintptr
		ExecuteScript                          uintptr
		CapturePreview                         uintptr
		Reload                                 uintptr
		PostWebMessageAsJSON                   uintptr
		PostWebMessageAsString                 uintptr
		AddWebMessageReceived                  uintptr
		RemoveWebMessageReceived               uintptr
		CallDevToolsProtocolMethod             uintptr
		GetBrowserProcessID                    uintptr
		GetCanGoBack                           uintptr
		GetCanGoForward                        uintptr
		GoBack                                 uintptr
		GoForward                              uintptr
		GetDevToolsProtocolEventReceiver       uintptr
		Stop                                   uintptr
		AddNewWindowRequested                  uintptr
		RemoveNewWindowRequested               uintptr
		AddDocumentTitleChanged                uintptr
		RemoveDocumentTitleChanged             uintptr
		GetDocumentTitle                       uintptr
		AddHostObjectToScript                  uintptr
		RemoveHostObjectFromScript             uintptr
		OpenDevToolsWindow                     uintptr
		AddContainsFullScreenElementChanged    uintptr
		RemoveContainsFullScreenElementChanged uintptr
		GetContainsFullScreenElement           uintptr
		AddWebResourceRequested                uintptr
		RemoveWebResourceRequested             uintptr
		AddWebResourceRequestedFilter          uintptr
		RemoveWebResourceRequestedFilter       uintptr
		AddWindowCloseRequested                uintptr
		RemoveWindowCloseRequested             uintptr
	}
)

var _GUIDCoreWebView22, _ = windows.GUIDFromString("{9E8F0CF8-E670-4B5E-B2BC-73E061E3184C}")

type (
	// ICoreWebView22 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView22 struct {
		VTBL *ICoreWebView22VTBL
	}

	// ICoreWebView22VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	ICoreWebView22VTBL struct {
		_ICoreWebView2VTBL
		AddWebResourceResponseReceived    uintptr
		RemoveWebResourceResponseReceived uintptr
		NavigateWithWebResourceRequest    uintptr
		AddDOMContentLoaded               uintptr
		RemoveDOMContentLoaded            uintptr
		CookieManager                     uintptr
		Environment                       uintptr
	}
)

type (
	// _ICoreWebView2Environment implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment
	_ICoreWebView2Environment struct {
		VTBL *_ICoreWebView2EnvironmentVTBL
	}

	// _ICoreWebView2EnvironmentVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2environment
	_ICoreWebView2EnvironmentVTBL struct {
		_IUnknownVTBL
		CreateCoreWebView2Controller     uintptr
		CreateWebResourceResponse        uintptr
		GetBrowserVersionString          uintptr
		AddNewBrowserVersionAvailable    uintptr
		RemoveNewBrowserVersionAvailable uintptr
	}
)

type (
	// _ICoreWebView2Controller implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2controller?view=webview2-1.0.622.22
	_ICoreWebView2Controller struct {
		VTBL *_ICoreWebView2ControllerVTBL
	}

	// _ICoreWebView2ControllerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2controller?view=webview2-1.0.622.22
	_ICoreWebView2ControllerVTBL struct {
		_IUnknownVTBL
		GetIsVisible                      uintptr
		PutIsVisible                      uintptr
		GetBounds                         uintptr
		PutBounds                         uintptr
		GetZoomFactor                     uintptr
		PutZoomFactor                     uintptr
		AddZoomFactorChanged              uintptr
		RemoveZoomFactorChanged           uintptr
		SetBoundsAndZoomFactor            uintptr
		MoveFocus                         uintptr
		AddMoveFocusRequested             uintptr
		RemoveMoveFocusRequested          uintptr
		AddGotFocus                       uintptr
		RemoveGotFocus                    uintptr
		AddLostFocus                      uintptr
		RemoveLostFocus                   uintptr
		AddAcceleratorKeyPressed          uintptr
		RemoveAcceleratorKeyPressed       uintptr
		GetParentWindow                   uintptr
		PutParentWindow                   uintptr
		NotifyParentWindowPositionChanged uintptr
		Close                             uintptr
		GetCoreWebView2                   uintptr
	}
)

type (
	// _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2environmentcompletedhandler.
	_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler struct {
		VTBL *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, err uintptr, val *_ICoreWebView2Environment) uintptr
	}

	// _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2environmentcompletedhandler.
	_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL = &_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2EnvironmentCompletedHandler, err uintptr, val *_ICoreWebView2Environment) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err, val)
	}),
}

type (
	// _ICoreWebView2CreateCoreWebView2ControllerCompletedHandler implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2controllercompletedhandler
	_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler struct {
		VTBL *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, err uintptr, val *_ICoreWebView2Controller) uintptr
	}

	// _ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2createcorewebview2controllercompletedhandler
	_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL = &_ICoreWebView2CreateCoreWebView2ControllerCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2CreateCoreWebView2ControllerCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2CreateCoreWebView2ControllerCompletedHandler, err uintptr, val *_ICoreWebView2Controller) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err, val)
	}),
}

type (
	// _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2addscripttoexecuteondocumentcreatedcompletedhandler?view=webview2-1.0.1264.42
	_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler struct {
		VTBL *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, err uintptr, id uintptr) uintptr
	}

	// _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2addscripttoexecuteondocumentcreatedcompletedhandler?view=webview2-1.0.1264.42
	_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTBL = &_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2AddScriptToExecuteOnDocumentCreatedCompletedHandler, err uintptr, id uintptr) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err, id)
	}),
}

type (
	// _ICoreWebView2CookieManager implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2cookiemanager
	_ICoreWebView2CookieManager struct {
		VTBL *_ICoreWebView2CookieManagerVTBL
	}

	// _ICoreWebView2CookieManagerVTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2cookiemanager
	_ICoreWebView2CookieManagerVTBL struct {
		_IUnknownVTBL
		CreateCookie                   uintptr
		CopyCookie                     uintptr
		GetCookies                     uintptr
		AddOrUpdateCookie              uintptr
		DeleteCookie                   uintptr
		DeleteCookies                  uintptr
		DeleteCookiesWithDomainAndPath uintptr
		DeleteAllCookies               uintptr
	}
)

type (
	_ICoreWebView2GetCookiesCompletedHandler struct {
		VTBL *_ICoreWebView2GetCookiesCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2GetCookiesCompletedHandler, err uintptr, cookies *_ICoreWebView2CookieList) uintptr
	}

	_ICoreWebView2GetCookiesCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2GetCookiesCompletedHandlerVTBL = &_ICoreWebView2GetCookiesCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2GetCookiesCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2GetCookiesCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2GetCookiesCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2GetCookiesCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2GetCookiesCompletedHandler, err uintptr, cookies *_ICoreWebView2CookieList) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err, cookies)
	}),
}

type (
	_ICoreWebView2CookieList struct {
		VTBL *_ICoreWebView2CookieListVTBL
	}

	_ICoreWebView2CookieListVTBL struct {
		_IUnknownVTBL
		Count           uintptr
		GetValueAtIndex uintptr
	}
)

type (
	_ICoreWebView2Cookie struct {
		VTBL *_ICoreWebView2CookieVTBL
	}

	_ICoreWebView2CookieVTBL struct {
		_IUnknownVTBL
		GetName     uintptr
		GetValue    uintptr
		PutValue    uintptr
		GetDomain   uintptr
		GetPath     uintptr
		GetExpires  uintptr
		PutExpires  uintptr
		IsHttpOnly  uintptr
		PutHttpOnly uintptr
		GetSameSite uintptr
		PutSameSite uintptr
		IsSecure    uintptr
		PutSecure   uintptr
		IsSession   uintptr
	}
)

type (
	_ICoreWebView2FrameWebMessageReceivedEventHandler struct {
		VTBL *_ICoreWebView2FrameWebMessageReceivedEventHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler, frame uintptr, args *_ICoreWebView2WebMessageReceivedEventArgs) uintptr
	}

	_ICoreWebView2FrameWebMessageReceivedEventHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var (
	_CoreWebView2FrameWebMessageReceivedEventHandlerVTBL = &_ICoreWebView2FrameWebMessageReceivedEventHandlerVTBL{
		_IUnknownVTBL: _IUnknownVTBL{
			Query: windows.NewCallback(func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler, _, o uintptr) uintptr {
				return 0
			}),
			Add: windows.NewCallback(func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler) uintptr {
				this.Counter += 1
				return this.Counter
			}),
			Release: windows.NewCallback(func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler) uintptr {
				this.Counter -= 1
				if this.Counter == 0 {
					*this = _ICoreWebView2FrameWebMessageReceivedEventHandler{}
				}
				return this.Counter + 1
			}),
		},
		Invoke: windows.NewCallback(func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler, frame uintptr, args *_ICoreWebView2WebMessageReceivedEventArgs) uintptr {
			if this == nil {
				return 0
			}
			return this.Invoke(this, frame, args)
		}),
	}
)

type (
	_ICoreWebView2WebMessageReceivedEventArgs struct {
		VTBL *_ICoreWebView2WebMessageReceivedEventArgsVTBL
	}

	_ICoreWebView2WebMessageReceivedEventArgsVTBL struct {
		_IUnknownVTBL
		Source                   uintptr
		WebMessageAsJson         uintptr
		TryGetWebMessageAsString uintptr
	}
)

type (
	_ICoreWebView2ExecuteScriptCompletedHandler struct {
		VTBL *_ICoreWebView2ExecuteScriptCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2ExecuteScriptCompletedHandler, err uintptr, resultObjectAsJson uintptr) uintptr
	}

	_ICoreWebView2ExecuteScriptCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2ExecuteScriptCompletedHandlerVTBL = &_ICoreWebView2ExecuteScriptCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2ExecuteScriptCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2ExecuteScriptCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2ExecuteScriptCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2ExecuteScriptCompletedHandler, err uintptr, resultObjectAsJson uintptr) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err, resultObjectAsJson)
	}),
}

type (
	_ICoreWebView2DocumentTitleChangedEventHandler struct {
		VTBL *_ICoreWebView2DocumentTitleChangedEventHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2DocumentTitleChangedEventHandler, w *_ICoreWebView2, v uintptr) uintptr
	}

	_ICoreWebView2DocumentTitleChangedEventHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2DocumentTitleChangedEventHandlerVTBL = &_ICoreWebView2DocumentTitleChangedEventHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2DocumentTitleChangedEventHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2DocumentTitleChangedEventHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2DocumentTitleChangedEventHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2DocumentTitleChangedEventHandler, w *_ICoreWebView2, v uintptr) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, w, v)
	}),
}

type (
	_ICoreWebView2SourceChangedEventHandler struct {
		VTBL *_ICoreWebView2SourceChangedEventHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2SourceChangedEventHandler, w *_ICoreWebView2, v uintptr) uintptr
	}

	_ICoreWebView2SourceChangedEventHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2SourceChangedEventHandlerVTBL = &_ICoreWebView2SourceChangedEventHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2SourceChangedEventHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2SourceChangedEventHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2SourceChangedEventHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2SourceChangedEventHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2SourceChangedEventHandler, w *_ICoreWebView2, v uintptr) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, w, v)
	}),
}
