package explorer

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Foundation/Foundation.h>

@interface explorer_file:NSObject
@property NSFileHandle* handler;
@property NSError* err;
@property NSURL* url;
@end

extern CFTypeRef newFile(CFTypeRef url);
extern uint64_t fileRead(CFTypeRef file, uint8_t *b, uint64_t len);
extern bool fileWrite(CFTypeRef file, uint8_t *b, uint64_t len);
extern bool fileClose(CFTypeRef file);
extern char* getError(CFTypeRef file);

*/
import "C"
import (
	"errors"
	"io"
	"unsafe"
)

type File struct {
	file   C.CFTypeRef
	url    C.CFTypeRef
	closed bool
}

func newFile(url C.CFTypeRef) (*File, error) {
	file := C.newFile(url)
	if err := getError(file); err != nil {
		return nil, err
	}
	return &File{file: file, url: url}, nil
}

func (f *File) Read(b []byte) (n int, err error) {
	if f.file == 0 || f.closed {
		return 0, io.ErrClosedPipe
	}

	buf := (*C.uint8_t)(unsafe.Pointer(&b[0]))
	length := C.uint64_t(uint64(len(b)))

	if n = int(int64(C.fileRead(f.file, buf, length))); n == 0 {
		if err := getError(f.file); err != nil {
			return n, err
		}
		return n, io.EOF
	}
	return n, nil
}

func (f *File) Write(b []byte) (n int, err error) {
	if f.file == 0 || f.closed {
		return 0, io.ErrClosedPipe
	}

	if len(b) == 0 {
		return 0, nil
	}

	buf := (*C.uint8_t)(unsafe.Pointer(&b[0]))
	length := C.uint64_t(int64(len(b)))

	if ok := bool(C.fileWrite(f.file, buf, length)); !ok {
		if err := getError(f.file); err != nil {
			return 0, err
		}
		return 0, errors.New("unknown error")
	}

	return len(b), nil
}

func (f *File) Close() error {
	defer func() {
		C.CFRelease(f.file)
		C.CFRelease(f.url)
	}()

	if ok := bool(C.fileClose(f.file)); !ok {
		return getError(f.file)
	}
	f.closed = true
	return nil
}

func getError(file C.CFTypeRef) error {
	// file will be 0 if the current device doesn't match with @available (i.e older than iOS 13).
	if file == 0 {
		return ErrNotAvailable
	}
	if err := C.GoString(C.getError(file)); len(err) > 0 {
		return errors.New(err)
	}
	return nil
}

// Exported function is required to create cgo header.
//
//export file_darwin
func file_darwin() {}
