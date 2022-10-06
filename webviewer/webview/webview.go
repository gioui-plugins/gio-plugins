package webview

//go:generate javac -source 8 -target 8 -bootclasspath $ANDROID_HOME\platforms\android-30\android.jar -d $TEMP\gowebview\classes sys_android.java
//go:generate jar cf sys_android.jar -C $TEMP\gowebview\classes .

import (
	"net/url"
)

// WebView is a webview.
type WebView interface {
	// Configure sets the configuration for the webview.
	Configure(config Config)

	// Resize resizes the webview to the specified size.
	// To make it invisible, set the size to (0, 0).
	// To make it visible, set the size to a non-zero value.
	//
	// The offset is the position of the webview relative to window coordinates,
	// assuming 0,0 as the top-left corner of the window.
	Resize(size Point, offset Point)

	// Navigate navigates to the specified URL.
	Navigate(url *url.URL)

	// Close closes and terminate the webview.
	Close()

	// Events returns actions that occur on the WebView.
	Events() chan Event

	// DataManager returns the DataManager for the webview.
	DataManager() DataManager

	// JavascriptManager returns the JavascriptManager for the webview.
	JavascriptManager() JavascriptManager
}

type internalWebView interface {
	attach(w *webview) error
	configure(w *webview, config Config)
	resize(w *webview, pos [4]float32)
	navigate(w *webview, url *url.URL)
	close(w *webview)
}

// Point is a point in the coordinate system.
type Point struct {
	X, Y float32
}

// NewWebView creates a new webview.
func NewWebView(config Config) (WebView, error) {
	return newWebview(config)
}

const (
	idMethodStart int = iota
	idMethodResize
	idMethodConfig
	idMethodNavigate
	idMethodClose
)
