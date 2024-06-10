package share

// Config is the configuration for a Share.
type Config struct {
	// RunOnMain is a function that runs on the main UI thread.
	RunOnMain func(f func())
}
