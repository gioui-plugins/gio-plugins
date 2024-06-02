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
	"sync"
	"unsafe"

	"gioui.org/app"
	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/explorer_explorer_android/classes explorer_android.java
//go:generate jar cf explorer_android.jar -C $TEMP/explorer_explorer_android/classes .

type driver struct {
	config Config
	mutex  sync.Mutex

	explorerObject jni.Object
	explorerClass  jni.Class

	explorerMethodOpen jni.MethodID
	explorerMethodSave jni.MethodID
}

func attachDriver(house *Explorer, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.mutex.Lock()
	old := driver.config.View
	driver.config = config
	driver.mutex.Unlock()

	if old != driver.config.View {
		driver.destroy()
		driver.init()
	}
}

// init will get all necessary MethodID (to future JNI calls) and get our Java library/class (which
// is defined on explorer_android.java file). The Java class doesn't retain information about the view,
// the view (GioView/GioActivity) is passed as argument for each openFile/saveFile function, so it
// can safely change between each call.
func (e *driver) init() error {
	return jni.Do(jni.JVMFor(e.config.VM), func(env jni.Env) error {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		if e.explorerClass != 0 && e.explorerObject != 0 {
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

		e.explorerObject = jni.NewGlobalRef(env, obj)
		e.explorerClass = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))
		e.explorerMethodOpen = jni.GetMethodID(env, e.explorerClass, "openFile", "(Landroid/view/View;Ljava/lang/String;I)V")
		e.explorerMethodSave = jni.GetMethodID(env, e.explorerClass, "saveFile", "(Landroid/view/View;Ljava/lang/String;I)V")

		return nil
	})
}

func (e *driver) destroy() {
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		if e.explorerObject != 0 {
			jni.DeleteGlobalRef(env, e.explorerObject)
			e.explorerObject = 0
		}
		if e.explorerClass != 0 {
			jni.DeleteGlobalRef(env, jni.Object(e.explorerClass))
			e.explorerClass = 0
		}
		return nil
	})
}

func (e *driver) saveFile(filename string, mime mimetype.MimeType) (io.WriteCloser, error) {
	res := make(chan result)
	callback := func(r result) {
		res <- r
	}

	go e.config.RunOnMain(func() {
		err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			e.mutex.Lock()
			defer e.mutex.Unlock()

			return jni.CallVoidMethod(env, e.explorerClass, e.explorerMethodSave,
				jni.Value(e.config.View),
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

func (e *driver) openFile(mime []mimetype.MimeType) (io.ReadCloser, error) {
	s := stringBuilderPool.Get().(*strings.Builder)

	res := make(chan result)
	callback := func(r result) {
		defer stringBuilderPool.Put(s)
		res <- r
		s.Reset()
	}

	go e.config.RunOnMain(func() {
		err := jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			e.mutex.Lock()
			defer e.mutex.Unlock()

			for i, v := range mime {
				if i > 0 {
					s.WriteRune(',')
				}
				v.WriteTo(s)
			}

			return jni.CallVoidMethod(env, e.explorerObject, e.explorerMethodOpen,
				jni.Value(e.config.View),
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
			file, err := newFile(env, jni.NewGlobalRef(env, jni.Object(uintptr(stream))))
			if file != nil {
				res.file = file
			}
			res.error = err
		}

		callback(res)
	}
}

var (
	_ io.ReadCloser  = (*File)(nil)
	_ io.WriteCloser = (*File)(nil)
)
