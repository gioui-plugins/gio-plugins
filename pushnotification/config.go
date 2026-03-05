package pushnotification

// ExternalConfig represents the external config.
// That is not locked to each OS to make easier to setup, otherwise
// you would need to create a config for each OS (yourfile_android.go, yourfile_ios.go, etc)
// and fill the config for each OS.
type ExternalConfig interface {
	implementsExternalConfig()
}

// AndroidFirebaseConfig represents the Firebase config for Android.
//
// Most of the fields are from "google-services.json", which you get when creating a Firebase project.
type AndroidFirebaseConfig struct {
	// AppID is your "mobilesdk_app_id" from "google-services.json".
	AppID string
	// ProjectID is your "project_id" from "google-services.json".
	ProjectID string
	// APIKey is your "api_key" from "google-services.json".
	APIKey string
	// SenderID is your "project number" from the Firebase console (or "project_info" from "google-services.json").
	SenderID string
}

// BrowserConfig represents the config for the browser.
type BrowserConfig struct {
	// VAPIDPublicKey is a self-generated public key in Base64
	VAPIDPublicKey string
}

// WindowsAzureConfig represents the config for Azure ObjectADD.
type WindowsAzureConfig struct {
	// ObjectID is the ObjectID of the Azure ObjectADD.
	ObjectID string
}

func (a AndroidFirebaseConfig) implementsExternalConfig() {}
func (b BrowserConfig) implementsExternalConfig()         {}
func (w WindowsAzureConfig) implementsExternalConfig()    {}
