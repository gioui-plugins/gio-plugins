//go:build darwin && !ios

package pushnotification

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework AppKit -framework UserNotifications
*/
import "C"

/* See _darwin.go */
