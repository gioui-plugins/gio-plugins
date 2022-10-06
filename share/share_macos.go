// SPDX-License-Identifier: Unlicense OR MIT

//go:build darwin && !ios
// +build darwin,!ios

package share

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc
#include <AppKit/AppKit.h>
#include <Foundation/Foundation.h>
#include <stdint.h>

static void openShare(CFTypeRef view, uint64_t x, uint64_t y, NSArray * obj) {
	NSSharingServicePicker * picker = [[NSSharingServicePicker alloc] initWithItems:obj];
	NSView * nview = (__bridge NSView *)view;

  	NSRect rect = NSMakeRect(x/2, y/2, 1, 1);
	[picker showRelativeToRect:rect ofView:nview preferredEdge:NSMinYEdge];
}

static void shareText(CFTypeRef view, uint64_t x, uint64_t y, char * text) {
	NSString * content = [NSString stringWithUTF8String: text];
	openShare(view, x, y, @[content]);
}

static void shareWebsite(CFTypeRef view, uint64_t x, uint64_t y, char * link) {
	NSString * content = [NSString stringWithUTF8String: link];
	openShare(view, x, y, @[[NSURL URLWithString:content]]);
}

*/
import "C"
import (
	"unsafe"

	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
)

type share struct {
	window *app.Window
	view   C.CFTypeRef
	size   [2]int64
}

func newShare(w *app.Window) share {
	return share{
		window: w,
	}
}

func (e *sharePlugin) listenEvents(evt event.Event) {
	switch evt := evt.(type) {
	case app.ViewEvent:
		e.view = C.CFTypeRef(evt.View)
	case system.FrameEvent:
		e.size[0] = int64(float32(evt.Size.X) / evt.Metric.PxPerDp)
		e.size[1] = int64(float32(evt.Size.Y) / evt.Metric.PxPerDp)
	}
}

func (e *sharePlugin) shareText(op TextOp) error {
	go e.window.Run(func() {
		t := C.CString(op.Text)
		defer C.free(unsafe.Pointer(t))
		C.shareText(e.view, C.uint64_t(e.size[0]), C.uint64_t(e.size[1]), t)
	})
	return nil
}

func (e *sharePlugin) shareWebsite(op WebsiteOp) error {
	go e.window.Run(func() {
		l := C.CString(op.Link)
		defer C.free(unsafe.Pointer(l))

		C.shareWebsite(e.view, C.uint64_t(e.size[0]), C.uint64_t(e.size[1]), l)
	})
	return nil
}
