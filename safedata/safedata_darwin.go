package safedata

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

extern uint8_t setSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len);
extern uint8_t updateSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len);
extern CFTypeRef getSecret(char * identifier, uint32_t * retLength);
extern uint8_t getSecretAt(CFTypeRef array, uint32_t index, char ** retId, char ** retDesc, uint8_t ** retData, uint32_t * sizeData);
extern uint8_t removeSecret(char * identifier);
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

type driver struct {
	prefix string
}

func attachDriver(house *SafeData, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.prefix = "_" + config.App + "_"
}

func (d driver) setSecret(secret Secret) error {
	id := C.CString(d.keyFor(secret.Identifier))
	defer C.free(unsafe.Pointer(id))

	desc := C.CString(secret.Description)
	defer C.free(unsafe.Pointer(desc))

	if len(secret.Data) == 0 {
		secret.Data = []byte{0x00}
	}
	defer runtime.KeepAlive(secret.Data)

	if err := d.getSecret(secret.Identifier, nil); err == nil {
		if err := C.updateSecret(id, desc, (*C.uint8_t)(unsafe.Pointer(&secret.Data[0])), C.uint64_t(uint64(len(secret.Data)))); err != 0 {
			return ErrUnsupported
		}
	} else {
		if err := C.setSecret(id, desc, (*C.uint8_t)(unsafe.Pointer(&secret.Data[0])), C.uint64_t(uint64(len(secret.Data)))); err != 0 {
			return ErrUnsupported
		}
	}
	return nil
}

func (d driver) listSecret(looper Looper) error {
	var (
		count C.uint32_t
	)
	ref := C.getSecret(nil, &count)
	if ref == 0 || count == 0 {
		return ErrNotFound
	}
	defer C.CFRelease(ref)

	for i := uint32(0); i < uint32(count); i++ {
		var idd *C.char
		ok := C.getSecretAt(ref, C.uint32_t(i), &idd, nil, nil, nil)
		if ok != 0 {
			return ErrNotFound
		}

		id := C.GoString(idd)
		if d.isOwned(id) {
			looper(d.rawKeyFor(id))
		}
	}

	return nil
}

func (d driver) getSecret(identifier string, secret *Secret) error {
	id := C.CString(d.keyFor(identifier))
	defer C.free(unsafe.Pointer(id))

	var (
		count C.uint32_t
	)
	ref := C.getSecret(id, &count)
	if ref == 0 || count == 0 {
		return ErrNotFound
	}
	defer C.CFRelease(ref)

	var (
		idd     *C.char
		desc    *C.char
		data    *C.uint8_t
		dataLen C.uint32_t
	)

	ok := C.getSecretAt(ref, 0, &idd, &desc, &data, &dataLen)
	if ok != 0 {
		return ErrNotFound
	}

	if secret == nil {
		return nil
	}

	secret.Identifier = C.GoString(idd)
	secret.Description = C.GoString(desc)

	if int(dataLen) > cap(secret.Data) {
		secret.Data = make([]byte, int(dataLen))
	} else {
		secret.Data = secret.Data[:int(dataLen)]
	}

	h := [3]uintptr{uintptr(unsafe.Pointer(data)), uintptr(dataLen), uintptr(dataLen)}
	copy(secret.Data, *(*[]byte)(unsafe.Pointer(&h)))

	return nil
}

func (d driver) removeSecret(identifier string) error {
	id := C.CString(d.keyFor(identifier))
	defer C.free(unsafe.Pointer(id))

	if err := C.removeSecret(id); err != 0 {
		return ErrNotFound
	}
	return nil
}

func (d driver) keyFor(s string) string {
	return d.prefix + s
}

func (d driver) rawKeyFor(s string) string {
	if len(s) < len(d.prefix) {
		return s
	}
	return s[len(d.prefix):]
}

func (d driver) isOwned(s string) bool {
	return strings.HasPrefix(s, d.prefix)
}
