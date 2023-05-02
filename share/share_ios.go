// SPDX-License-Identifier: Unlicense OR MIT

package share

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc
#include <UIKit/UIKit.h>
#include <stdint.h>

static CFTypeRef openShare(CFTypeRef viewController, NSArray * obj) {
	UIActivityViewController * activityController = [[UIActivityViewController alloc] initWithActivityItems:obj applicationActivities:nil];
	[(__bridge UIViewController *)viewController presentViewController:activityController animated:YES completion:nil];

	return CFBridgingRetain(activityController);
}

static CFTypeRef shareText(CFTypeRef viewController, char * text) {
	NSString * content = [NSString stringWithUTF8String: text];
	return openShare(viewController, @[content]);
}

static CFTypeRef shareWebsite(CFTypeRef viewController, char * link) {
	NSString * content = [NSString stringWithUTF8String: link];
	return openShare(viewController, @[[NSURL URLWithString:content]]);
}

*/
import "C"
import (
	"sync"
	"unsafe"
)

type driver struct {
	mutex  sync.Mutex
	config Config

	shareObject C.CFTypeRef
}

func attachDriver(house *Share, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	driver.config = config
}

func (e *driver) shareText(title, text string) error {
	go e.config.RunOnMain(func() {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		t := C.CString(text)
		defer C.free(unsafe.Pointer(t))

		if e.shareObject != 0 {
			C.CFRelease(e.shareObject)
		}
		e.shareObject = C.shareText(C.CFTypeRef(e.config.View), t)
	})
	return nil
}

func (e *driver) shareWebsite(title, description, url string) error {
	go e.config.RunOnMain(func() {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		l := C.CString(url)
		defer C.free(unsafe.Pointer(l))

		if e.shareObject != 0 {
			C.CFRelease(e.shareObject)
		}
		e.shareObject = C.shareWebsite(C.CFTypeRef(e.config.View), l)
	})
	return nil
}
