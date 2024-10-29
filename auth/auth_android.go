package auth

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"git.wow.st/gmp/jni"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/auth/providers/google"
	_ "github.com/gioui-plugins/gio-plugins/auth/vendors/android/jar"
	"runtime/cgo"
	"sync"
	"unsafe"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	authObject jni.Object
	authClass  jni.Class

	generalAuthMethodOpen jni.MethodID
	googleAuthMethodOpen  jni.MethodID

	send func(event Event)

	cgoHandler cgo.Handle
}

func attachDriver(house *Auth, config Config) {
	house.driver = driver{send: house.sendResponse}
	house.driver.cgoHandler = cgo.NewHandle(&house.driver)
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	old := driver.config.View
	driver.config = config

	if old != driver.config.View {
		destroyDriver(driver)
		if err := initDriver(driver); err != nil {
			panic(fmt.Errorf("auth: failed to initialize: %w", err))
		}
	}
}

func initDriver(e *driver) error {
	return jni.Do(jni.JVMFor(e.config.VM), func(env jni.Env) error {
		if e.authClass != 0 && e.authObject != 0 {
			return nil // Already initialized
		}

		class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(e.config.Context)), "com/inkeliz/auth_android/auth_android")
		if err != nil {
			return err
		}

		obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", `()V`))
		if err != nil {
			return err
		}

		e.authObject = jni.NewGlobalRef(env, obj)
		e.authClass = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))
		e.googleAuthMethodOpen = jni.GetMethodID(env, e.authClass, "openNative", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;J)V")
		e.generalAuthMethodOpen = jni.GetMethodID(env, e.authClass, "openGeneral", "(Landroid/view/View;Ljava/lang/String;J)V")

		return nil
	})
}

func destroyDriver(e *driver) error {
	return jni.Do(jni.JVMFor(e.config.VM), func(env jni.Env) error {
		if e.authObject != 0 {
			jni.DeleteGlobalRef(env, e.authObject)
			e.authObject = 0
		}
		if e.authClass != 0 {
			jni.DeleteGlobalRef(env, jni.Object(e.authClass))
			e.authClass = 0
		}
		return nil
	})
}

func (e *driver) open(provider providers.Provider, nonce string) error {
	switch provider := provider.(type) {
	case *google.Provider:
		return e.openNative(provider, nonce)
	default:
		return e.openAny(provider, nonce)
	}
}

func (e *driver) openAny(provider providers.Provider, nonce string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return jni.Do(jni.JVMFor(e.config.VM), func(env jni.Env) error {
		if e.authObject == 0 {
			return ErrNotConfigured
		}

		return jni.CallVoidMethod(env, e.authObject, e.generalAuthMethodOpen,
			jni.Value(e.config.View),
			jni.Value(jni.JavaString(env, provider.URL(nonce))),
			jni.Value(e.cgoHandler),
		)
	})
}

func (e *driver) openNative(provider *google.Provider, nonce string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	return jni.Do(jni.JVMFor(e.config.VM), func(env jni.Env) error {
		if e.authObject == 0 {
			return ErrNotConfigured
		}

		return jni.CallVoidMethod(env, e.authObject, e.googleAuthMethodOpen,
			jni.Value(e.config.View),
			jni.Value(jni.JavaString(env, provider.WebClientID)),
			jni.Value(jni.JavaString(env, nonce)),
			jni.Value(e.cgoHandler),
		)
	})
}

//export Java_com_inkeliz_auth_1android_auth_1android_NativeAuthCallback
func Java_com_inkeliz_auth_1android_auth_1android_NativeAuthCallback(env *C.JNIEnv, _ C.jclass, handler C.jlong, idToken C.jstring) {
	r := AuthenticatedEvent{
		Provider: google.IdentifierGoogle,
		IDToken:  jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(idToken)),
	}

	driver, ok := cgo.Handle(handler).Value().(*driver)
	if !ok {
		panic("auth: invalid handler")
	}

	if r.IDToken == "" {
		driver.send(ErrorEvent{Error: fmt.Errorf("auth: empty id token")})
	} else {
		driver.send(r)
	}
}
