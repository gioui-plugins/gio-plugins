//go:build ios
// +build ios

package hyperlink

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

@import UIKit;

void openLink(char *u) {
	[[UIApplication sharedApplication] openURL:[NSURL URLWithString: @(u)] options:@{} completionHandler:nil];
}
*/
import "C"

import (
	"gioui.org/io/event"
	"net/url"
	"unsafe"
)

func (*hyperlinkPlugin) listenEvents(_ event.Event) {}

func (*hyperlinkPlugin) open(u *url.URL) error {
	u.RawQuery = u.Query().Encode()
	cURL := C.CString(u.String())
	C.openLink(cURL)
	C.free(unsafe.Pointer(cURL))
	return nil
}
