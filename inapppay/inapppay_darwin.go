package inapppay

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Foundation/Foundation.h>
#include <stdint.h>
#include <stdlib.h>

typedef struct {
	char* id;
	char* title;
	char* description;
	char* price;
	char* currencyCode;
} gioplugins_inapppay_product_t;

typedef struct {
	char* productID;
	char* purchaseID;
	int status;
    char* developerPayload;
	char* originalJSON;
	char* signature;
} gioplugins_inapppay_purchase_result_t;

extern void gioplugins_inapppay_create(uintptr_t data);
extern void gioplugins_inapppay_list_products(uintptr_t data, char** productIDs, int count);
extern void gioplugins_inapppay_purchase(uintptr_t data, char* productID, char* developerPayload);
*/
import "C"
import (
	"errors"
	"runtime/cgo"
	"sync"
	"unsafe"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	send        func(event Event)
	cgoHandler  cgo.Handle
	initialized bool
}

func attachDriver(iap *InAppPay, config Config) {
	d := &driver{send: iap.sendResponse}
	d.cgoHandler = cgo.NewHandle(d)
	iap.driver = d
	configureDriver(iap.driver, config)
}

func configureDriver(d *driver, config Config) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	d.config = config
	if d.initialized == false && config.RunOnMain != nil {
		config.RunOnMain(func() {
			C.gioplugins_inapppay_create(C.uintptr_t(d.cgoHandler))
		})
		d.initialized = true
	}
}

func (d *driver) listProducts(productIDs []string) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cArray := make([]*C.char, len(productIDs))
	for i, id := range productIDs {
		cArray[i] = C.CString(id)
		defer C.free(unsafe.Pointer(cArray[i]))
	}

	ptr := (**C.char)(nil)
	if len(cArray) > 0 {
		ptr = (**C.char)(unsafe.Pointer(&cArray[0]))
	}

	C.gioplugins_inapppay_list_products(C.uintptr_t(d.cgoHandler), ptr, C.int(len(productIDs)))
	return nil
}

func (d *driver) purchase(productID string, developerPayload string, _ bool) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	cs := C.CString(productID)
	defer C.free(unsafe.Pointer(cs))

	dv := C.CString(developerPayload)
	defer C.free(unsafe.Pointer(dv))

	C.gioplugins_inapppay_purchase(C.uintptr_t(d.cgoHandler), cs, dv)
	return nil
}

//export gioplugins_inapppay_on_product_details
func gioplugins_inapppay_on_product_details(handle C.uintptr_t, cProducts *C.gioplugins_inapppay_product_t, count C.int) {
	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	// Iterate C array
	products := make([]Product, int(count))
	if count > 0 {
		// Using manual arithmetic to be safe with unknowns.
		size := unsafe.Sizeof(C.gioplugins_inapppay_product_t{})
		base := uintptr(unsafe.Pointer(cProducts))

		for i := 0; i < int(count); i++ {
			p := (*C.gioplugins_inapppay_product_t)(unsafe.Pointer(base + uintptr(i)*size))
			products[i] = Product{
				ID:           C.GoString(p.id),
				Title:        C.GoString(p.title),
				Description:  C.GoString(p.description),
				Price:        C.GoString(p.price),
				CurrencyCode: C.GoString(p.currencyCode),
			}
		}
	}

	d.send(ProductDetailsEvent{Products: products})
}

//export gioplugins_inapppay_on_purchase_result
func gioplugins_inapppay_on_purchase_result(handle C.uintptr_t, res C.gioplugins_inapppay_purchase_result_t) {
	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	result := PaymentResultEvent{
		ProductID:        C.GoString(res.productID),
		PurchaseID:       C.GoString(res.purchaseID),
		Status:           PaymentStatus(res.status),
		DeveloperPayload: C.GoString(res.developerPayload),
		OriginalJSON:     C.GoString(res.originalJSON),
		Signature:        C.GoString(res.signature),
	}

	d.send(result)
}

//export gioplugins_inapppay_report_error
func gioplugins_inapppay_report_error(handle C.uintptr_t, msg *C.char) {
	message := C.GoString(msg)
	d, ok := cgo.Handle(handle).Value().(*driver)
	if !ok {
		return
	}

	d.send(ErrorEvent{Error: errors.New(message)})
}
