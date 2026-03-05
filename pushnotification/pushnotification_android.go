//go:build android

package pushnotification

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"log"
	"runtime/cgo"
	"sync"
	"unsafe"

	"git.wow.st/gmp/jni"

	_ "github.com/gioui-plugins/gio-plugins/pushnotification/vendors/android/jar"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	clientObject jni.Object
	clientClass  jni.Class

	// Methods
	getTokenMethod   jni.MethodID
	initializeMethod jni.MethodID
}

func attachDriver(push *Push, config Config) {
	d := &driver{}
	push.driver = d
	configureDriver(d, config)
}

func configureDriver(d *driver, config Config) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	old := d.config.View
	d.config = config

	if old != d.config.View {
		if err := initDriver(d); err != nil {
			panic(err)
		}
	}
}

func initDriver(d *driver) error {
	return jni.Do(jni.JVMFor(d.config.VM), func(env jni.Env) error {
		if d.clientClass != 0 && d.clientObject != 0 {
			return nil
		}

		class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(d.config.Context)), "com/inkeliz/pushnotification_android/pushnotification_android")
		if err != nil {
			log.Println("Push: Failed to load class", err)
			return err
		}

		obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", "()V"))
		if err != nil {
			log.Println("Push: Failed to create object", err)
			return err
		}

		d.clientObject = jni.NewGlobalRef(env, obj)
		d.clientClass = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))

		d.getTokenMethod = jni.GetMethodID(env, d.clientClass, "getToken", "(Landroid/view/View;J)V")
		d.initializeMethod = jni.GetMethodID(env, d.clientClass, "initialize", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;Ljava/lang/String;)V")

		// Initialize Firebase with config
		err = jni.CallVoidMethod(env, d.clientObject, d.initializeMethod,
			jni.Value(d.config.View),
			jni.Value(jni.JavaString(env, d.config.AppID)),
			jni.Value(jni.JavaString(env, d.config.ProjectID)),
			jni.Value(jni.JavaString(env, d.config.APIKey)),
			jni.Value(jni.JavaString(env, d.config.SenderID)),
		)
		if err != nil {
			return err
		}

		return nil
	})
}

func (d *driver) requestToken() (Token, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	c := make(chan struct {
		Token Token
		Error error
	}, 1)
	defer close(c)

	fn := func(token Token, err error) {
		c <- struct {
			Token Token
			Error error
		}{
			Token: token,
			Error: err,
		}
	}

	err := jni.Do(jni.JVMFor(d.config.VM), func(env jni.Env) error {
		if d.clientObject == 0 {
			return ErrNotConfigured
		}

		return jni.CallVoidMethod(env, d.clientObject, d.getTokenMethod,
			jni.Value(d.config.View),
			jni.Value(cgo.NewHandle(fn)),
		)
	})
	if err != nil {
		return Token{}, err
	}

	r := <-c
	return r.Token, r.Error
}

//export Java_com_inkeliz_pushnotification_1android_pushnotification_1android_onTokenReceived
func Java_com_inkeliz_pushnotification_1android_pushnotification_1android_onTokenReceived(env *C.JNIEnv, _ C.jclass, handler C.jlong, token C.jstring) {
	t := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(token))

	h := cgo.Handle(handler)
	fn, ok := h.Value().(func(Token, error))
	if !ok {
		return
	}
	defer h.Delete()

	fn(Token{Token: t, Platform: PlatformAndroid}, nil)
}

//export Java_com_inkeliz_pushnotification_1android_pushnotification_1android_onError
func Java_com_inkeliz_pushnotification_1android_pushnotification_1android_onError(env *C.JNIEnv, _ C.jclass, handler C.jlong, msg C.jstring) {
	m := jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(msg))

	h := cgo.Handle(handler)
	fn, ok := h.Value().(func(Token, error))
	if !ok {
		return
	}
	defer h.Delete()

	fn(Token{}, fmt.Errorf("android push error: %s", m))
}
