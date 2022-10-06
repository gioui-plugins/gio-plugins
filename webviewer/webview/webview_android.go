//go:build android
// +build android

package webview

import (
	"net/url"

	"git.wow.st/gmp/jni"
)

type driver struct {
	config Config

	clsWebView jni.Class
	objWebView jni.Object
}

func (r *driver) attach(w *webview) (err error) {
	w.scheduler.SetRunner(w.driver.config.RunOnMain)

	w.driver.config.RunOnMain(func() {
		err = jni.Do(r.config.VM, func(env jni.Env) error {
			cls, err := jni.LoadClass(env, jni.ClassLoaderFor(env, w.driver.config.Context), "com/inkeliz/webview/sys_android")
			if err != nil {
				return err
			}

			// [Android] We need to create an GlobalRef of our class, otherwise we can't manipulate that afterwards.
			r.clsWebView = jni.Class(jni.NewGlobalRef(env, jni.Object(cls)))

			obj, err := jni.NewObject(env, r.clsWebView, jni.GetMethodID(env, r.clsWebView, "<init>", `()V`))
			if err != nil {
				return err
			}

			// [Android] We need to create an GlobalRef of our class, otherwise we can't manipulate that afterwards.
			r.objWebView = jni.Object(jni.NewGlobalRef(env, obj))

			err = r.callArgs("webview_create", "(Landroid/view/View;J)V", func(env jni.Env) []jni.Value {
				return []jni.Value{
					jni.Value(r.config.View),
					jni.Value(uintptr(w.handle)),
				}
			})

			if err != nil {
				return err
			}

			if err := r.setProxy(); err != nil {
				return err
			}

			if err := r.setCerts(); err != nil {
				return err
			}

			w.javascriptManager = newJavascriptManager(w)
			w.dataManager = newDataManager(w)

			return nil
		})
	})

	if err != nil {
		return err
	}

	return nil
}

func (r *driver) configure(w *webview, config Config) {
	r.config = config
	w.scheduler.SetRunner(w.driver.config.RunOnMain)
}

func (r *driver) resize(w *webview, pos [4]float32) {
	if pos[2] == 0 && pos[3] == 0 {
		if w.visible {
			r.call("webview_hide", "()V")
			w.visible = false
		}
	} else {
		r.callArgs("webview_resize", "(IIII)V", func(env jni.Env) []jni.Value {
			return []jni.Value{
				jni.Value(int32(pos[0] + 0.5)),
				jni.Value(int32(pos[1] + 0.5)),
				jni.Value(int32(pos[2] + 0.5)),
				jni.Value(int32(pos[3] + 0.5)),
			}
		})
		if !w.visible {
			r.call("webview_show", "()V")
			w.visible = true
		}
	}
}

func (r *driver) navigate(w *webview, url *url.URL) {
	r.callArgs("webview_navigate", "(Ljava/lang/String;)V", func(env jni.Env) []jni.Value {
		return []jni.Value{
			jni.Value(jni.JavaString(env, url.String())),
		}
	})
}

func (r *driver) close(w *webview) {
	if r.objWebView == 0 || r.clsWebView == 0 {
		return
	}

	r.call("webview_destroy", "()V")

	go jni.Do(r.config.VM, func(env jni.Env) error {
		jni.DeleteGlobalRef(env, jni.Object(r.clsWebView))
		jni.DeleteGlobalRef(env, r.objWebView)

		return nil
	})

	r.objWebView, r.clsWebView = 0, 0
}

func (r *driver) call(name, sig string) (err error) {
	// The arguments may need the `env`
	// In that case there's no input, so it's using func(env jni.Env) []jni.Value { return nil } instead
	return r.callArgs(name, sig, func(env jni.Env) []jni.Value { return nil })
}

func (r *driver) callArgs(name, sig string, args func(env jni.Env) []jni.Value) (err error) {
	err = jni.Do(r.config.VM, func(env jni.Env) error {
		return jni.CallVoidMethod(env, r.objWebView, jni.GetMethodID(env, r.clsWebView, name, sig), args(env)...)
	})
	return err
}

func (r *driver) callBooleanArgs(name, sig string, args func(env jni.Env) []jni.Value) (b bool, err error) {
	err = jni.Do(r.config.VM, func(env jni.Env) error {
		b, err = jni.CallBooleanMethod(env, r.objWebView, jni.GetMethodID(env, r.clsWebView, name, sig), args(env)...)
		return err
	})
	return b, err
}
