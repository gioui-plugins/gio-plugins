//go:build ios

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
	"net/url"
	"unsafe"
)

type driver struct{}

func attachDriver(house *Hyperlink, config Config) {}

func configureDriver(driver *driver, config Config) {}

func (*driver) open(u *url.URL) error {
	u.RawQuery = u.Query().Encode()
	cURL := C.CString(u.String())
	C.openLink(cURL)
	C.free(unsafe.Pointer(cURL))
	return nil
}
