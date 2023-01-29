package webview

/*
#cgo CFLAGS: -Werror
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"strings"
	"sync"
	"unsafe"

	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

type javascriptManager struct {
	webview   *webview
	jsHandler internal.Handle
	callbacks sync.Map // map[string]func(message string)
}

func newJavascriptManager(w *webview) *javascriptManager {
	r := &javascriptManager{webview: w}
	r.jsHandler = internal.NewHandle(r)
	r.installCallback()
	r.installJavascript(fmt.Sprintf(scriptCallback, `_callback.callback`), JavascriptOnLoadStart)
	return r
}

func (j *javascriptManager) installCallback() {
	j.webview.driver.callArgs("webview_set_callback", "(J)V", func(env jni.Env) []jni.Value {
		return []jni.Value{jni.Value(j.jsHandler)}
	})
}

func (j *javascriptManager) RunJavaScript(js string) error {
	done := make(chan error, 1)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	j.webview.scheduler.MustRun(func() {
		j.webview.driver.callArgs("webview_run_javascript", "(Ljava/lang/String;J)V", func(env jni.Env) []jni.Value {
			return []jni.Value{
				jni.Value(jni.JavaString(env, js)),
				jni.Value(int64(dr)),
			}
		})
	})

	return <-done
}

func (j *javascriptManager) InstallJavascript(js string, when JavascriptInstallationTime) (err error) {
	j.webview.scheduler.MustRun(func() {
		err = j.installJavascript(js, when)
	})
	return err
}

func (j *javascriptManager) installJavascript(js string, when JavascriptInstallationTime) error {
	done := make(chan error, 1)
	dr := internal.NewHandle(done)
	defer dr.Delete()

	j.webview.driver.callArgs("webview_install_javascript", "(Ljava/lang/String;JJ)V", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, js)),
			jni.Value(when),
			jni.Value(int64(dr)),
		}
	})

	return <-done
}

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

//export Java_com_inkeliz_webview_sys_1android_sendCallback
func Java_com_inkeliz_webview_sys_1android_sendCallback(env *C.JNIEnv, class C.jclass, ptr C.jlong, msg C.jstring) {
	receiveCallback(uintptr(C.jlong(ptr)), jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(msg)))
}
