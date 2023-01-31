//go:build !android && !darwin && !ios && !windows && !(js && wasm)

package webview

type webview struct{}

func newWebview(config Config) (*webview, error) {
	return nil, ErrNotSupported
}
