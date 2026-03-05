//go:build android

package inapppay

// Config holds the configuration for the InAppPay.
type Config struct {
	// View is the Android View, used to perform the purchase.
	View uintptr
	// VM is the Java VM, used to call Java methods.
	VM uintptr
	// Context is the Android Context, used to load classes.
	Context uintptr
	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(func())
}
