// SPDX-License-Identifier: Unlicense OR MIT

package share

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/explorer/classes share_android.java
//go:generate jar cf share_android.jar -C $TEMP/explorer/classes .

import "C"
import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"git.wow.st/gmp/jni"
)

type share struct {
	window *app.Window
	view   uintptr

	shareClass         jni.Class
	shareMethodText    jni.MethodID
	shareMethodWebsite jni.MethodID
}

func newShare(w *app.Window) share {
	return share{window: w}
}

func (e *share) init() {
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		share, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(app.AppContext())), "com/inkeliz/share_android/share_android")
		if err != nil {
			return err
		}

		e.shareClass = jni.Class(jni.NewGlobalRef(env, jni.Object(share)))
		e.shareMethodText = jni.GetStaticMethodID(env, share, "shareText", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;)V")
		e.shareMethodWebsite = jni.GetStaticMethodID(env, share, "shareWebsite", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V")

		return nil
	})
}

func (e *share) destroy() {
	if e.shareClass == 0 {
		return
	}
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		jni.DeleteGlobalRef(env, jni.Object(e.shareClass))
		return nil
	})
}

func (e *sharePlugin) listenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		e.view = evt.View
		e.init()
	case system.DestroyEvent:

	}
}

func (e *sharePlugin) shareText(op TextOp) error {
	e.window.Run(func() {
		jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			title, text := jni.JavaString(env, op.Title), jni.JavaString(env, op.Text)
			return jni.CallStaticVoidMethod(env, e.shareClass, e.shareMethodText, jni.Value(e.view), jni.Value(title), jni.Value(text))
		})
	})

	return nil
}

func (e *sharePlugin) shareWebsite(op WebsiteOp) error {
	e.window.Run(func() {
		jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			title, text, link := jni.JavaString(env, op.Title), jni.JavaString(env, op.Text), jni.JavaString(env, op.Link)
			return jni.CallStaticVoidMethod(env, e.shareClass, e.shareMethodWebsite, jni.Value(e.view), jni.Value(title), jni.Value(text), jni.Value(link))
		})
	})
	return nil
}
