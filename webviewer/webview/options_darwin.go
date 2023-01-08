package webview

/*
#cgo CFLAGS: -xobjective-c -fmodules -fobjc-arc

#import <Foundation/Foundation.h>

extern void enableDebug(CFTypeRef config);
*/
import "C"

func (r *driver) setDebug() {
	options.Lock()
	d := options.debug
	options.Unlock()
	if d <= 0 {
		return
	}

	C.enableDebug(r.webviewConfig)
}
