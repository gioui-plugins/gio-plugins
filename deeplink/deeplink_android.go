package deeplink

//go:generate mkdir -p $TEMP/share/classes
//go:generate javac -source 8 -target 8  -bootclasspath $ANDROID_HOME/platforms/android-30/android.jar -d $TEMP/share/classes deeplink_android.java
//go:generate jar cf deeplink_android.jar -C $TEMP/share/classes .

/*
#cgo CFLAGS: -Werror
#cgo LDFLAGS: -landroid

#include <android/native_window_jni.h>
*/
import "C"
import (
	"fmt"
	"git.wow.st/gmp/jni"
	"unsafe"
)

//export Java_com_inkeliz_deeplink_deeplink_1android_ReceiveScheme
func Java_com_inkeliz_deeplink_deeplink_1android_ReceiveScheme(env *C.JNIEnv, _ C.jclass, scheme C.jstring) {
	fmt.Println("demo: Java_com_inkeliz_deeplink_1android_ReceiveScheme")
	_Deeplink.updateURL(jni.GoString(jni.EnvFor(uintptr(unsafe.Pointer(env))), jni.String(scheme)))
}
