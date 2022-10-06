package webview

import (
	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
	"strconv"
)

type JavascriptManager interface {
	// RunJavaScript execute the given JavaScript code into the loaded page.
	//
	// The behavior of this function is undefined if the page is not loaded.
	RunJavaScript(js string) error

	// InstallJavascript installs the given JavaScript code as a function,
	// that will be injected into the page when it is loaded.
	//
	// The behavior of this function is undefined if the page is already loaded.
	InstallJavascript(js string, when JavascriptInstallationTime) error

	// AddCallback registers a callback for the given JavaScript function.
	//
	// The function can be called as follows:
	//
	// 		window.callback.<name>(<message>);
	//
	// The function will be called with the message provided by the JavaScript caller,
	// you probably want to encode (in JSON/Protobuf/Flatbuffers/Karmem/CapNProto or
	// another more efficient format). Due to different browser implementations and memory layouts, it's not possible to
	// pass a pointer to a struct to the JavaScript function or another type rather than
	// a string.
	//
	// It may introduce performance penalties, so use it wisely. In order to be compatible
	// with multiples webviews and OSes, it may introduce more indirection and use ProxyAPI
	// and ReflectionAPI and other slow APIs to achieve the same functionality
	//
	// If you want to return a result to the JavaScript caller, use RunJavaScript and
	// set the result into some array or global variable.
	//
	// The name must be unique and must not contain any dots, and must have a maximum
	// length of 255 characters.
	AddCallback(name string, fn func(message string)) error
}

// JavascriptInstallationTime defines when the JavaScript code is injected.
type JavascriptInstallationTime int64

const (
	// JavascriptOnLoadStart ensures that the JavaScript code is injected on page load start,
	// before the page is fully loaded.
	JavascriptOnLoadStart JavascriptInstallationTime = iota

	// JavascriptOnLoadFinish ensures that the JavaScript code is injected on page load end, when
	// all contents are loaded.
	JavascriptOnLoadFinish
)

// scriptCallback uses ProxyAPI to provide the function, which is called by JS.
// The "%s" is replaced by the native function, which varies across each OS
// implementation.
var scriptCallback = `
	globalThis.callback = new Proxy({}, {
		get(self, name) {
			return function(message) {
				let size = (name.toString().length).toString(16).toUpperCase().padStart(2, '0');
				%s(size + name.toString() + message);
			}
		},
	});
`

func receiveCallback(handler uintptr, in string) {
	if len(in) < 3 {
		return
	}

	size, err := strconv.ParseUint(in[:2], 16, 8)
	if err != nil {
		return
	}

	name := in[2 : 2+size]
	message := in[2+size:]

	j := internal.Handle(handler).Value().(*javascriptManager)
	if fn, ok := j.callbacks.Load(name); ok {
		fn.(func(message string))(message)
	}
}
