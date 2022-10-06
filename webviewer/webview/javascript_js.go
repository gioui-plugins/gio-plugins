package webview

import (
	"sync"
)

type javascriptManager struct {
	callbacks sync.Map
}

func newJavascriptManager(w *webview) *javascriptManager {
	return new(javascriptManager)
}

// RunJavaScript implements the JavascriptManager interface.
func (*javascriptManager) RunJavaScript(_ string) error {
	return ErrNotSupported
}

// InstallJavascript implements the JavascriptManager interface.
func (*javascriptManager) InstallJavascript(_ string, _ JavascriptInstallationTime) error {
	return ErrNotSupported
}

// AddCallback implements the JavascriptManager interface.
func (*javascriptManager) AddCallback(_ string, _ func(message string)) error {
	return ErrNotSupported
}
