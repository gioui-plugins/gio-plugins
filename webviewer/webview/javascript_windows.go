package webview

import (
	"fmt"
	"strings"
	"sync"
	"syscall"
	"unsafe"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
	"golang.org/x/sys/windows"
)

type javascriptManager struct {
	webview   *webview
	jsHandler internal.Handle
	callbacks sync.Map // map[string]func(message string)
	callback  *_ICoreWebView2FrameWebMessageReceivedEventHandler
}

func newJavascriptManager(w *webview) *javascriptManager {
	r := &javascriptManager{webview: w}
	r.jsHandler = internal.NewHandle(r)
	w.scheduler.MustRun(func() {
		r.installCallback()
		r.installJavascript(fmt.Sprintf(scriptCallback, `window.chrome.webview.postMessage`))
	})
	return r
}

func (j *javascriptManager) installCallback() {
	j.callback = &_ICoreWebView2FrameWebMessageReceivedEventHandler{
		VTBL: _CoreWebView2FrameWebMessageReceivedEventHandlerVTBL,
		Invoke: func(this *_ICoreWebView2FrameWebMessageReceivedEventHandler, frame uintptr, args *_ICoreWebView2WebMessageReceivedEventArgs) uintptr {
			var message *uint16
			syscall.SyscallN(
				args.VTBL.TryGetWebMessageAsString,
				uintptr(unsafe.Pointer(args)),
				uintptr(unsafe.Pointer(&message)),
			)
			if message != nil {
				receiveCallback(uintptr(j.jsHandler), windows.UTF16PtrToString(message))
			}
			return 0
		},
	}

	j.webview.scheduler.MustRun(func() {
		var r uint64
		syscall.SyscallN(
			j.webview.driver.webview2.VTBL.AddWebMessageReceived,
			uintptr(unsafe.Pointer(j.webview.driver.webview2)),
			uintptr(unsafe.Pointer(j.callback)),
			uintptr(unsafe.Pointer(&r)),
		)
	})
}

// RunJavaScript implements the JavascriptManager interface.
func (j *javascriptManager) RunJavaScript(js string) error {
	done := make(chan error)
	fr := internal.NewHandle(done)
	defer fr.Delete()

	for _, c := range js {
		if c == 0x00 {
			return ErrInvalidJavascript
		}
	}

	handler := &_ICoreWebView2ExecuteScriptCompletedHandler{
		VTBL: _CoreWebView2ExecuteScriptCompletedHandlerVTBL,
		Invoke: func(this *_ICoreWebView2ExecuteScriptCompletedHandler, err uintptr, resulAsJson uintptr) uintptr {
			done <- nil
			return 0
		},
	}

	text := syscall.StringToUTF16Ptr(js)

	j.webview.scheduler.MustRun(func() {
		syscall.SyscallN(
			j.webview.driver.webview2.VTBL.ExecuteScript,
			uintptr(unsafe.Pointer(j.webview.driver.webview2)),
			uintptr(unsafe.Pointer(text)),
			uintptr(unsafe.Pointer(handler)),
		)
	})

	return <-done
}

// InstallJavascript implements the JavascriptManager interface.
func (j *javascriptManager) InstallJavascript(js string, when JavascriptInstallationTime) error {
	if when == JavascriptOnLoadFinish {
		js = fmt.Sprintf(`document.addEventListener('DOMContentLoaded', function() { %s };`, js)
	}

	j.webview.scheduler.MustRun(func() {
		j.installJavascript(js)
	})

	return nil
}

func (j *javascriptManager) installJavascript(js string) {
	syscall.SyscallN(
		j.webview.driver.webview2.VTBL.AddScriptToExecuteOnDocumentCreated,
		uintptr(unsafe.Pointer(j.webview.driver.webview2)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(js))),
	)
}

// AddCallback implements the JavascriptManager interface.
func (j *javascriptManager) AddCallback(name string, fn func(message string)) error {
	if len(name) > 255 {
		return ErrJavascriptCallbackInvalidName
	}
	if strings.Contains(name, ".") || strings.Contains(name, " ") {
		return ErrJavascriptCallbackInvalidName
	}
	if _, ok := j.callbacks.Load(name); ok {
		return ErrJavascriptCallbackDuplicate
	}

	j.callbacks.Store(name, fn)
	return nil
}
