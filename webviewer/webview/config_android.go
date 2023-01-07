package webview

import (
	"git.wow.st/gmp/jni"
)

// Config is the configuration for a WebView.
type Config struct {
	// View is the Android View.
	View jni.Class

	// VM is the Java VM.
	VM jni.JVM

	// Context is the Android Context.
	Context jni.Object

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())

	// PxPerDp represents how many pixels per each dp.
	PxPerDp float32
}
