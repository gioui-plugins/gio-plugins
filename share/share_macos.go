// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

package share

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc
#include <AppKit/AppKit.h>
#include <Foundation/Foundation.h>
#include <stdint.h>

static CFTypeRef openShare(CFTypeRef view, uint64_t x, uint64_t y, NSArray * obj) {
	NSSharingServicePicker * picker = [[NSSharingServicePicker alloc] initWithItems:obj];
	NSView * nview = (__bridge NSView *)view;

  	NSRect rect = NSMakeRect(x/2, y/2, 1, 1);
	[picker showRelativeToRect:rect ofView:nview preferredEdge:NSMinYEdge];

	return CFBridgingRetain(picker);
}

static CFTypeRef shareText(CFTypeRef view, uint64_t x, uint64_t y, char * text) {
	NSString * content = [NSString stringWithUTF8String: text];
	return openShare(view, x, y, @[content]);
}

static CFTypeRef shareWebsite(CFTypeRef view, uint64_t x, uint64_t y, char * link) {
	NSString * content = [NSString stringWithUTF8String: link];
	return openShare(view, x, y, @[[NSURL URLWithString:content]]);
}

*/
import "C"
import (
	"unsafe"
)

type driver struct {
	config Config

	shareObject C.CFTypeRef
}

func attachDriver(house *Share, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.config = config
	driver.config.Size = [2]float32{
		config.Size[0] / config.PxPerDp,
		config.Size[1] / config.PxPerDp,
	}
}

func (e *driver) shareText(title, text string) error {
	go e.config.RunOnMain(func() {
		t := C.CString(text)
		defer C.free(unsafe.Pointer(t))

		if e.shareObject != 0 {
			C.CFRelease(e.shareObject)
		}
		e.shareObject = C.shareText(C.CFTypeRef(e.config.View), C.uint64_t(uint64(e.config.Size[0])), C.uint64_t(uint64(e.config.Size[1])), t)
	})
	return nil
}

func (e *driver) shareWebsite(title, description, url string) error {
	go e.config.RunOnMain(func() {
		l := C.CString(url)
		defer C.free(unsafe.Pointer(l))

		if e.shareObject != 0 {
			C.CFRelease(e.shareObject)
		}
		e.shareObject = C.shareWebsite(C.CFTypeRef(e.config.View), C.uint64_t(uint64(e.config.Size[0])), C.uint64_t(uint64(e.config.Size[1])), l)
	})
	return nil
}
