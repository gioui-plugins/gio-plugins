// SPDX-License-Identifier: Unlicense OR MIT

package share

//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/explorer/classes share_android.java
//go:generate jar cf share_android.jar -C $TEMP/explorer/classes .

import "C"
import (
	"gioui.org/app"
	"git.wow.st/gmp/jni"
)

type driver struct {
	config Config

	shareClass         jni.Class
	shareMethodText    jni.MethodID
	shareMethodWebsite jni.MethodID
}

func attachDriver(house *Share, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	old := driver.config.View
	driver.config = config
	if old != driver.config.View {
		driver.destroy()
		driver.init()
	}
}

func (e *driver) init() {
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

func (e *driver) destroy() {
	if e.shareClass == 0 {
		return
	}
	jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
		jni.DeleteGlobalRef(env, jni.Object(e.shareClass))
		return nil
	})
}

func (e *driver) shareText(title, text string) error {
	e.config.RunOnMain(func() {
		jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			title, text := jni.JavaString(env, title), jni.JavaString(env, text)
			return jni.CallStaticVoidMethod(env, e.shareClass, e.shareMethodText, jni.Value(e.config.View), jni.Value(title), jni.Value(text))
		})
	})

	return nil
}

func (e *driver) shareWebsite(title, description, url string) error {
	e.config.RunOnMain(func() {
		jni.Do(jni.JVMFor(app.JavaVM()), func(env jni.Env) error {
			title, text, link := jni.JavaString(env, title), jni.JavaString(env, description), jni.JavaString(env, url)
			return jni.CallStaticVoidMethod(env, e.shareClass, e.shareMethodWebsite, jni.Value(e.config.View), jni.Value(title), jni.Value(text), jni.Value(link))
		})
	})
	return nil
}
