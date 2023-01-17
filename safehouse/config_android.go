package safehouse

import (
	"git.wow.st/gmp/jni"
)

// Config is the configuration for a WebView.
type Config struct {
	// App is the name of the app, which is used, in some
	// OSes to identify the app who creates the
	// credentials.
	App string

	// VM is the Java VM.
	VM jni.JVM

	// Context is the Android Context.
	Context jni.Object

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
