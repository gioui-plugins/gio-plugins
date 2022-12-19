package explorer

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"io"
	"path/filepath"
	"runtime"
	"runtime/cgo"
	"strings"
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/explorer_explorer_android/classes explorer_android.java
//go:generate jar cf explorer_android.jar -C $TEMP/explorer_explorer_android/classes .

type explorer struct {
	window *app.Window
	view   uintptr

	libObject jni.Object
	libClass  jni.Class

	openMethodFile jni.MethodID
	saveMethodFile jni.MethodID
}

func (e *explorerPlugin) listenEvents(evt event.Event) {
	if evt, ok := evt.(app.ViewEvent); ok {
		e.view = evt.View
	}
}

// init will get all necessary MethodID (to future JNI calls) and get our Java library/class (which
// is defined on explorer_android.java file). The Java class doesn't retain information about the view,
// the view (GioView/GioActivity) is passed as argument for each openFile/saveFile function, so it
// can safely change between each call.
func (e *explorerPlugin) init(env jni.Env) error {
	if e.libObject != 0 && e.libClass != 0 {
		return nil // Already initialized
	}

	class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "org/gioui/x/explorer/explorer_android")
	if err != nil {
		return err
	}

	obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", `()V`))
	if err != nil {
		return err
	}

	e.libObject = jni.NewGlobalRef(env, obj)
	e.libClass = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))
	e.openMethodFile = jni.GetMethodID(env, e.libClass, "openFile", "(Landroid/view/View;Ljava/lang/String;I)V")
	e.saveMethodFile = jni.GetMethodID(env, e.libClass, "saveFile", "(Landroid/view/View;Ljava/lang/String;I)V")

	return nil
}

func (e *explorerPlugin) destroy() {
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		if e.libObject != 0 {
			jni.DeleteGlobalRef(env, e.libObject)
			e.libObject = 0
		}
		if e.libClass != 0 {
			jni.DeleteGlobalRef(env, jni.Object(e.libClass))
			e.libClass = 0
		}
		return nil
	})
}

func (e *explorerPlugin) saveFile(filename string, mime mimetype.MimeType) (io.WriteCloser, error) {
	res := make(chan result)
	callback := func(r result) {
		res <- r
	}

	go e.window.Run(func() {
		err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			e.mutex.Lock()
			defer e.mutex.Unlock()

			if err := e.init(env); err != nil {
				return err
			}

			return jni.CallVoidMethod(env, e.libObject, e.explorer.saveMethodFile,
				jni.Value(e.view),
				jni.Value(jni.JavaString(env, strings.TrimPrefix(strings.ToLower(filepath.Ext(filename)), "."))),
				jni.Value(cgo.NewHandle(callback)),
			)
		})

		if err != nil {
			res <- result{error: err}
		}
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	runtime.KeepAlive(callback)
	return r.file.(io.WriteCloser), nil
}

func (e *explorerPlugin) openFile(mime []mimetype.MimeType) (io.ReadCloser, error) {
	s := stringBuilderPool.Get().(*strings.Builder)

	res := make(chan result)
	callback := func(r result) {
		defer stringBuilderPool.Put(s)
		res <- r
		s.Reset()
	}

	go e.window.Run(func() {
		err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			e.mutex.Lock()
			defer e.mutex.Unlock()

			if err := e.init(env); err != nil {
				return err
			}

			for i, v := range mime {
				if i > 0 {
					s.WriteRune(',')
				}
				v.WriteTo(s)
			}

			return jni.CallVoidMethod(env, e.libObject, e.explorer.openMethodFile,
				jni.Value(e.view),
				jni.Value(jni.JavaString(env, s.String())),
				jni.Value(cgo.NewHandle(callback)),
			)
		})
		if err != nil {
			res <- result{error: err}
		}
	})

	r := <-res
	if r.error != nil {
		return nil, r.error
	}
	runtime.KeepAlive(callback)
	return r.file.(io.ReadCloser), nil
}

//export Java_org_gioui_x_explorer_explorer_1android_ImportCallback
func Java_org_gioui_x_explorer_explorer_1android_ImportCallback(env *C.JNIEnv, _ C.jclass, stream C.jobject, handler C.jlong, err C.jstring) {
	fileCallback(env, stream, handler, err)
}

//export Java_org_gioui_x_explorer_explorer_1android_ExportCallback
func Java_org_gioui_x_explorer_explorer_1android_ExportCallback(env *C.JNIEnv, _ C.jclass, stream C.jobject, handler C.jlong, err C.jstring) {
	fileCallback(env, stream, handler, err)
}

func fileCallback(env *C.JNIEnv, stream C.jobject, handler C.jlong, err C.jstring) {
	if callback, ok := cgo.Handle(handler).Value().(func(result)); ok {
		var res result
		env := jni.EnvFor(uintptr(unsafe.Pointer(env)))
		if stream == 0 {
			res.error = ErrUserDecline
			if err != 0 {
				if err := jni.GoString(env, jni.String(uintptr(err))); len(err) > 0 {
					res.error = errors.New(err)
				}
			}
		} else {
			res.file, res.error = newFile(env, jni.NewGlobalRef(env, jni.Object(uintptr(stream))))
		}

		callback(res)
	}
}

var (
	_ io.ReadCloser  = (*File)(nil)
	_ io.WriteCloser = (*File)(nil)
)
