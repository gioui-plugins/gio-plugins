//go:build windows

package webview

import (
	"os"
	"path/filepath"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	_CreateCoreWebView2EnvironmentWithOptions *windows.Proc
	_Ole32                                    = windows.NewLazySystemDLL("ole32.dll")
	_Ole32CoTaskMemAlloc                      = _Ole32.NewProc("CoTaskMemAlloc")
)

func init() {
	// current app name
	appName := "gio-webview2-"
	if len(os.Args) > 0 {
		appName = filepath.Base(os.Args[0])
	}

	dllPath := filepath.Join(os.TempDir(), appName+"webview2loader.dll")
	dst, err := os.Create(dllPath)
	if err != nil {
		panic(err)
	}

	if _, err := dst.Write(dllFile); err != nil {
		panic(err)
	}

	if err := dst.Close(); err != nil {
		panic(err)
	}

	dll, err := windows.LoadDLL(dllPath)
	if err != nil {
		panic(err)
	}

	_CreateCoreWebView2EnvironmentWithOptions, err = dll.FindProc("CreateCoreWebView2EnvironmentWithOptions")
	if err != nil {
		panic(err)
	}
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
	// _ICoreWebView22 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView22 struct {
		VTBL *_ICoreWebView22VTBL
	}

	// _ICoreWebView22VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView22VTBL struct {
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
	// _ICoreWebView23 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView23 struct {
		VTBL *_ICoreWebView23VTBL
	}

	// _ICoreWebView23VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView23VTBL struct {
		_ICoreWebView22VTBL
		TrySuspend                          uintptr
		Resume                              uintptr
		IsSuspended                         uintptr
		SetVirtualHostNameToFolderMapping   uintptr
		ClearVirtualHostNameToFolderMapping uintptr
	}
)

type (
	// _ICoreWebView24 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView24 struct {
		VTBL *_ICoreWebView24VTBL
	}

	// _ICoreWebView24VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView24VTBL struct {
		_ICoreWebView23VTBL
		AddFrameCreated        uintptr
		RemoveFrameCreated     uintptr
		AddDownloadStarting    uintptr
		RemoveDownloadStarting uintptr
	}
)

type (
	// _ICoreWebView25 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView25 struct {
		VTBL *_ICoreWebView25VTBL
	}

	// _ICoreWebView25VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView25VTBL struct {
		_ICoreWebView24VTBL
		AddClientCertificateRequested    uintptr
		RemoveClientCertificateRequested uintptr
	}
)

type (
	// _ICoreWebView26 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView26 struct {
		VTBL *_ICoreWebView26VTBL
	}

	// _ICoreWebView26VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView26VTBL struct {
		_ICoreWebView25VTBL
		OpenTaskManagerWindow uintptr
	}
)

type (
	// _ICoreWebView27 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView27 struct {
		VTBL *_ICoreWebView27VTBL
	}

	// _ICoreWebView27VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView27VTBL struct {
		_ICoreWebView26VTBL
		PrintToPDF uintptr
	}
)

type (
	// _ICoreWebView28 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView28 struct {
		VTBL *_ICoreWebView28VTBL
	}

	// _ICoreWebView28VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView28VTBL struct {
		_ICoreWebView27VTBL
		AddIsMutedChanged                   uintptr
		RemoveIsMutedChanged                uintptr
		GetIsMuted                          uintptr
		PutIsMuted                          uintptr
		AddIsDocumentPlayingAudioChanged    uintptr
		RemoveIsDocumentPlayingAudioChanged uintptr
		IsDocumentPlayingAudio              uintptr
	}
)

type (
	// _ICoreWebView29 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView29 struct {
		VTBL *_ICoreWebView29VTBL
	}

	// _ICoreWebView29VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView29VTBL struct {
		_ICoreWebView28VTBL
		AddIsDefaultDownloadDialogOpenChanged    uintptr
		RemoveIsDefaultDownloadDialogOpenChanged uintptr
		GetIsDefaultDownloadDialogOpen           uintptr
		OpenDefaultDownloadDialog                uintptr
		CloseDefaultDownloadDialog               uintptr
		GetDefaultDownloadDialogCornerAlignment  uintptr
		PutDefaultDownloadDialogCornerAlignment  uintptr
		GetDefaultDownloadDialogMargin           uintptr
		PutDefaultDownloadDialogMargin           uintptr
	}
)

type (
	// _ICoreWebView210 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView210 struct {
		VTBL *_ICoreWebView210VTBL
	}

	// _ICoreWebView210VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView210VTBL struct {
		_ICoreWebView29VTBL
		AddBasicAuthenticationRequested    uintptr
		RemoveBasicAuthenticationRequested uintptr
	}
)

type (
	// _ICoreWebView211 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView211 struct {
		VTBL *_ICoreWebView211VTBL
	}

	// _ICoreWebView211VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView211VTBL struct {
		_ICoreWebView210VTBL
		CallDevToolsProtocolMethodForSession uintptr
		AddContextMenuRequested              uintptr
		RemoveContextMenuRequested           uintptr
	}
)

type (
	// _ICoreWebView212 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView212 struct {
		VTBL *_ICoreWebView212VTBL
	}

	// _ICoreWebView212VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView212VTBL struct {
		_ICoreWebView211VTBL
		AddStatusBarTextChanged    uintptr
		RemoveStatusBarTextChanged uintptr
		GetStatusBarText           uintptr
	}
)

var _GUIDCoreWebView213, _ = windows.GUIDFromString("{F75F09A8-667E-4983-88D6-C8773F315E84}")

type (
	// _ICoreWebView213 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView213 struct {
		VTBL *_ICoreWebView213VTBL
	}

	// _ICoreWebView213VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView213VTBL struct {
		_ICoreWebView212VTBL
		GetProfile uintptr
	}
)

type (
	// _ICoreWebView214 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView214 struct {
		VTBL *_ICoreWebView214VTBL
	}

	// _ICoreWebView214VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView214VTBL struct {
		_ICoreWebView213VTBL
		AddServerCertificateErrorDetected    uintptr
		RemoveServerCertificateErrorDetected uintptr
		ClearServerCertificateErrorAction    uintptr
	}
)

type (
	// _ICoreWebView215 implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView215 struct {
		VTBL *_ICoreWebView215VTBL
	}

	// _ICoreWebView215VTBL implements https://docs.microsoft.com/en-us/microsoft-edge/webview2/reference/win32/icorewebview2_2?view=webview2-1.0.622.22
	_ICoreWebView215VTBL struct {
		_ICoreWebView214VTBL
		AddFaviconChanged    uintptr
		RemoveFaviconChanged uintptr
		GetFaviconUri        uintptr
		GetFavicon           uintptr
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

type (
	_ICoreWebView2Profile struct {
		VTBL *_ICoreWebView2ProfileVTBL
	}

	_ICoreWebView2ProfileVTBL struct {
		_IUnknownVTBL
		GetProfileName               uintptr
		IsInPrivateModeEnabled       uintptr
		ProfilePath                  uintptr
		GetDefaultDownloadFolderPath uintptr
		PutDefaultDownloadFolderPath uintptr
		GetPreferredColorScheme      uintptr
		PutPreferredColorScheme      uintptr
	}
)

var _GUIDCoreWebView2Profile2, _ = windows.GUIDFromString("{fa740d4b-5eae-4344-a8ad-74be31925397}")

type (
	_ICoreWebView2Profile2 struct {
		VTBL *_ICoreWebView2Profile2VTBL
	}

	_ICoreWebView2Profile2VTBL struct {
		_ICoreWebView2ProfileVTBL
		ClearBrowsingData            uintptr
		ClearBrowsingDataInTimeRange uintptr
		ClearBrowsingDataAll         uintptr
	}
)

type (
	_ICoreWebView2ClearBrowsingDataCompletedHandler struct {
		VTBL *_ICoreWebView2ClearBrowsingDataCompletedHandlerVTBL

		Counter uintptr
		Invoke  func(this *_ICoreWebView2ClearBrowsingDataCompletedHandler, err uintptr) uintptr
	}

	_ICoreWebView2ClearBrowsingDataCompletedHandlerVTBL struct {
		_IUnknownVTBL
		Invoke uintptr
	}
)

var _CoreWebView2ClearBrowsingDataCompletedHandler = &_ICoreWebView2ClearBrowsingDataCompletedHandlerVTBL{
	_IUnknownVTBL: _IUnknownVTBL{
		Query: windows.NewCallback(func(this *_ICoreWebView2ClearBrowsingDataCompletedHandler, _, o uintptr) uintptr {
			return 0
		}),
		Add: windows.NewCallback(func(this *_ICoreWebView2ClearBrowsingDataCompletedHandler) uintptr {
			referenceHolder[unsafe.Pointer(this)] = struct{}{}
			this.Counter += 1
			return this.Counter
		}),
		Release: windows.NewCallback(func(this *_ICoreWebView2ClearBrowsingDataCompletedHandler) uintptr {
			this.Counter -= 1
			if this.Counter == 0 {
				delete(referenceHolder, unsafe.Pointer(this))
				*this = _ICoreWebView2ClearBrowsingDataCompletedHandler{}
			}
			return this.Counter + 1
		}),
	},
	Invoke: windows.NewCallback(func(this *_ICoreWebView2ClearBrowsingDataCompletedHandler, err uintptr) uintptr {
		if this == nil {
			return 0
		}
		return this.Invoke(this, err)
	}),
}
