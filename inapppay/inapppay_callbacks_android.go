package inapppay

/*
#cgo LDFLAGS: -landroid

#include <jni.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"unsafe"

	"git.wow.st/gmp/jni"
)

//export Java_com_inkeliz_inapppay_1android_inapppay_1android_nativeOnProductDetails
func Java_com_inkeliz_inapppay_1android_inapppay_1android_nativeOnProductDetails(
	env *C.JNIEnv,
	clazz C.jclass,
	handle C.jlong,
	ids C.jobjectArray,
	titles C.jobjectArray,
	descriptions C.jobjectArray,
	prices C.jobjectArray,
	currencies C.jobjectArray,
	count int32,
) {
	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	jEnv := jni.EnvFor(uintptr(unsafe.Pointer(env)))

	products := make([]Product, int(count))
	for i := 0; i < int(count); i++ {
		id, err := jni.GetObjectArrayElement(jEnv, jni.ObjectArray(ids), jni.Size(i))
		if err != nil {
			d.send(ErrorEvent{Error: errors.New("failed to get product ID: " + err.Error())})
			return
		}
		products[i].ID = jni.GoString(jEnv, jni.String(id))

		title, err := jni.GetObjectArrayElement(jEnv, jni.ObjectArray(titles), jni.Size(i))
		if err != nil {
			d.send(ErrorEvent{Error: errors.New("failed to get product title: " + err.Error())})
			return
		}
		products[i].Title = jni.GoString(jEnv, jni.String(title))

		desc, err := jni.GetObjectArrayElement(jEnv, jni.ObjectArray(descriptions), jni.Size(i))
		if err != nil {
			d.send(ErrorEvent{Error: errors.New("failed to get product description: " + err.Error())})
			return
		}
		products[i].Description = jni.GoString(jEnv, jni.String(desc))

		price, err := jni.GetObjectArrayElement(jEnv, jni.ObjectArray(prices), jni.Size(i))
		if err != nil {
			d.send(ErrorEvent{Error: errors.New("failed to get product price: " + err.Error())})
			return
		}
		products[i].Price = jni.GoString(jEnv, jni.String(price))

		currency, err := jni.GetObjectArrayElement(jEnv, jni.ObjectArray(currencies), jni.Size(i))
		if err != nil {
			d.send(ErrorEvent{Error: errors.New("failed to get product currency code: " + err.Error())})
			return
		}
		products[i].CurrencyCode = jni.GoString(jEnv, jni.String(currency))
	}

	d.send(ProductDetailsEvent{Products: products})
}

//export Java_com_inkeliz_inapppay_1android_inapppay_1android_NativeOnPurchaseResult
func Java_com_inkeliz_inapppay_1android_inapppay_1android_NativeOnPurchaseResult(
	env *C.JNIEnv,
	clazz C.jclass,
	handle C.jlong,
	productID C.jstring,
	purchaseID C.jstring,
	status C.jint,
	originalJSON C.jstring,
	signature C.jstring,
) {
	jEnv := jni.EnvFor(uintptr(unsafe.Pointer(env)))

	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	result := PaymentResultEvent{
		ProductID:    jni.GoString(jEnv, jni.String(productID)),
		PurchaseID:   jni.GoString(jEnv, jni.String(purchaseID)),
		Status:       PaymentStatus(status),
		OriginalJSON: jni.GoString(jEnv, jni.String(originalJSON)),
		Signature:    jni.GoString(jEnv, jni.String(signature)),
	}

	d.send(result)
}

//export Java_com_inkeliz_inapppay_1android_inapppay_1android_NativeOnReportError
func Java_com_inkeliz_inapppay_1android_inapppay_1android_NativeOnReportError(
	env *C.JNIEnv,
	clazz C.jclass,
	handle C.jlong,
	msg C.jstring,
) {
	jEnv := jni.EnvFor(uintptr(unsafe.Pointer(env)))

	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	d.send(ErrorEvent{Error: errors.New(jni.GoString(jEnv, jni.String(uintptr(unsafe.Pointer(msg)))))})
}
