package safedata

// Config is the configuration for a SafeData.
type Config struct {
	// App is the name of the app, which is used, in some
	// OSes to identify the app who creates the
	// credentials.
	App string

	// VM is the Java VM.
	VM uintptr

	// Context is the Android Context.
	Context uintptr

	// Folder is the place where the data is encrypted stored.
	Folder string

	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
