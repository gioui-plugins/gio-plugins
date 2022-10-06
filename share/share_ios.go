// SPDX-License-Identifier: Unlicense OR MIT

package share

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc
#include <UIKit/UIKit.h>
#include <stdint.h>

static void openShare(CFTypeRef viewController, NSArray * obj) {
	UIActivityViewController * activityController = [[UIActivityViewController alloc] initWithActivityItems:obj applicationActivities:nil];
	[(__bridge UIViewController *)viewController presentViewController:activityController animated:YES completion:nil];
}

static void shareText(CFTypeRef viewController, char * text) {
	NSString * content = [NSString stringWithUTF8String: text];
	openShare(viewController, @[content]);
}

static void shareWebsite(CFTypeRef viewController, char * link) {
	NSString * content = [NSString stringWithUTF8String: link];
	openShare(viewController, @[[NSURL URLWithString:content]]);
}

*/
import "C"
import (
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
)

type share struct {
	window         *app.Window
	viewController C.CFTypeRef
}

func newShare(w *app.Window) share {
	return share{
		window: w,
	}
}

func (e *sharePlugin) listenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		e.viewController = C.CFTypeRef(evt.ViewController)
	}
}

func (e *sharePlugin) shareText(op TextOp) error {
	go e.window.Run(func() {
		t := C.CString(op.Text)
		defer C.free(unsafe.Pointer(t))
		C.shareText(e.viewController, t)
	})
	return nil
}

func (e *sharePlugin) shareWebsite(op WebsiteOp) error {
	go e.window.Run(func() {
		l := C.CString(op.Link)
		defer C.free(unsafe.Pointer(l))
		C.shareWebsite(e.viewController, l)
	})
	return nil
}
