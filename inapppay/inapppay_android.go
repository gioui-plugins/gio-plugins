//go:build android && !aptoide

package inapppay

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"runtime/cgo"
	"sync"

	_ "github.com/gioui-plugins/gio-plugins/inapppay/vendors/android/jar"

	"git.wow.st/gmp/jni"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	clientObject jni.Object
	clientClass  jni.Class

	// Methods
	listProductsMethod jni.MethodID
	purchaseMethod     jni.MethodID

	send       func(event Event)
	cgoHandler cgo.Handle
}

func attachDriver(iap *InAppPay, config Config) {
	d := &driver{send: iap.sendResponse}
	d.cgoHandler = cgo.NewHandle(d)
	iap.driver = d
	configureDriver(d, config)
}

func configureDriver(d *driver, config Config) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	old := d.config.View
	d.config = config

	if old != d.config.View {
		if err := initDriver(d); err != nil {
			d.send(ErrorEvent{Error: fmt.Errorf("inapppay: failed to initialize: %w", err)})
		}
	}
}

func initDriver(d *driver) error {
	return jni.Do(jni.JVMFor(d.config.VM), func(env jni.Env) error {
		if d.clientClass != 0 && d.clientObject != 0 {
			return nil
		}

		class, err := jni.LoadClass(env, jni.ClassLoaderFor(env, jni.Object(d.config.Context)), "com/inkeliz/inapppay_android/inapppay_android")
		if err != nil {
			return err
		}

		obj, err := jni.NewObject(env, class, jni.GetMethodID(env, class, "<init>", "()V"))
		if err != nil {
			return err
		}

		d.clientObject = jni.NewGlobalRef(env, obj)
		d.clientClass = jni.Class(jni.NewGlobalRef(env, jni.Object(class)))

		// Map generic methods to Google implementation
		// Note signature update: takes Arrays, no JSON
		// listProducts(Context, String[], long)
		d.listProductsMethod = jni.GetMethodID(env, d.clientClass, "listProductsGoogle", "(Landroid/view/View;[Ljava/lang/String;J)V")

		// purchase(Activity, String, String, long)
		d.purchaseMethod = jni.GetMethodID(env, d.clientClass, "purchaseGoogle", "(Landroid/view/View;Ljava/lang/String;Ljava/lang/String;IJ)V")

		return nil
	})
}

func (d *driver) listProducts(productIDs []string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	return jni.Do(jni.JVMFor(d.config.VM), func(env jni.Env) error {
		if d.clientObject == 0 {
			return ErrNotConfigured
		}

		// Convert []string to ObjectArray
		strClass := jni.FindClass(env, "java/lang/String")
		jArray := jni.NewObjectArray(env, jni.Size(len(productIDs)), strClass, jni.Object(0))
		for i, id := range productIDs {
			jni.SetObjectArrayElement(env, jArray, jni.Size(i), jni.Object(jni.JavaString(env, id)))
		}

		return jni.CallVoidMethod(env, d.clientObject, d.listProductsMethod,
			jni.Value(d.config.View),
			jni.Value(jArray),
			jni.Value(d.cgoHandler),
		)
	})
}

func (d *driver) purchase(productID string, developerPayload string, isPersonalized bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	personalized := int32(0)
	if isPersonalized {
		personalized = 1
	}

	return jni.Do(jni.JVMFor(d.config.VM), func(env jni.Env) error {
		if d.clientObject == 0 {
			return ErrNotConfigured
		}

		return jni.CallVoidMethod(env, d.clientObject, d.purchaseMethod,
			jni.Value(d.config.View),
			jni.Value(jni.JavaString(env, productID)),
			jni.Value(jni.JavaString(env, developerPayload)),
			jni.Value(personalized),
			jni.Value(d.cgoHandler),
		)
	})
}
