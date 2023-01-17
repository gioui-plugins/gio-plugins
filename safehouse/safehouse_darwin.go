package safehouse

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#include <stdint.h>
#import <Foundation/Foundation.h>

extern void setSecret(char * identifier, char * desc, uint8_t * value, uint64_t value_len);
extern CFTypeRef getSecret(char * identifier, uint8_t ** ret, uint32_t * size, char ** desc, uint8_t multiple);
extern uint8_t removeSecret(char * identifier);
*/
import "C"
import "unsafe"

type driver struct{}

func attachDriver(house *SafeHouse) {
	house.driver = driver{}
}

func (d driver) setSecret(secret Secret) error {
	id := C.CString(secret.Identifier)
	defer C.free(unsafe.Pointer(id))

	desc := C.CString(secret.Description)
	defer C.free(unsafe.Pointer(desc))

	if len(secret.Data) == 0 {
		secret.Data = []byte{0x00}
	}

	C.setSecret(id, desc, (*C.uint8_t)(unsafe.Pointer(&secret.Data[0])), C.uint64_t(uint64(len(secret.Data))))
	return nil
}

func (d driver) listSecret(looper Looper) error {
	return nil
}

func (d driver) getSecret(identifier string, secret *Secret) error {
	id := C.CString(identifier)
	defer C.free(unsafe.Pointer(id))

	var dData *C.char
	var cData *C.uint8_t
	var cSize C.uint32_t

	ref := C.getSecret(id, &cData, &cSize, &dData, 0)
	if ref == 0 {
		return ErrNotFound
	}
	defer func() {
		C.CFRelease(ref)
	}()

	if cData == nil {
		return ErrNotFound
	}

	secret.Description = C.GoString(dData)
	secret.Identifier = identifier

	if len(secret.Data) < int(cSize) {
		secret.Data = make([]byte, int(cSize))
	}

	h := [3]uintptr{uintptr(unsafe.Pointer(cData)), uintptr(cSize), uintptr(cSize)}
	copy(secret.Data, *(*[]byte)(unsafe.Pointer(&h)))

	return nil
}

func (d driver) removeSecret(identifier string) error {
	id := C.CString(identifier)
	defer C.free(unsafe.Pointer(id))

	if err := C.removeSecret(id); err != 0 {
		return ErrNotFound
	}
	return nil
}
