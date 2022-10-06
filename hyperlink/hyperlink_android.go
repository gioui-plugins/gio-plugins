//go:build android
// +build android

package hyperlink

import (
	"net/url"
	"sync"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"git.wow.st/gmp/jni"
)

//go:generate javac -source 8 -target 8 -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d ./hyperlink/classes hyperlink_android.java
//go:generate jar cf hyperlink_android.jar -C ./hyperlink/classes .

type hyperlink struct {
	mutex   sync.Mutex
	view    uintptr
	context jni.Object
	mid     jni.MethodID
	cls     jni.Class
}

func (h *hyperlinkPlugin) listenEvents(event event.Event) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	switch e := event.(type) {
	case app.ViewEvent:
		h.view = e.View
		if e.View != 0 {
			h.init()
		}
	case system.DestroyEvent:
		h.destroy()
	}
}

func (h *hyperlinkPlugin) init() {
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "com/inkeliz/hyperlink_android/hyperlink_android")
		if err != nil {
			panic(err)
		}

		h.cls = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))
		h.mid = jni.GetStaticMethodID(env, h.cls, "open", "(Landroid/view/View;Ljava/lang/String;)V")

		return nil
	})
}

func (h *hyperlinkPlugin) destroy() {
	if h.cls == 0 {
		return
	}
	jni.DeleteGlobalRef(jni.EnvFor(app.JavaVM()), jni.Object(h.cls))
	h.cls = 0
	h.mid = nil
}

func (h *hyperlinkPlugin) open(u *url.URL) error {
	if h.view == 0 {
		return ErrNotReady
	}

	return jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		h.mutex.Lock()
		defer h.mutex.Unlock()

		err := jni.CallStaticVoidMethod(env, h.cls, h.mid, jni.Value(h.view), jni.Value(jni.JavaString(env, u.String())))
		if err != nil {
			panic(err)
		}

		return nil
	})
}
