# webviewer

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/webviewer.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/webviewer)

Give some Webview to your Gio and Golang application. 😎

**Currently, GioWebview doesn't work with Gio out-of-box and requires some patches to work.**

------------------

## Usage

### Example

```go
// Define what WebView is been used (the Tag specify the tag).
wv := webviewer.WebViewOp{Tag: &something}.Push(gtx.Ops)

// Show the WebView with the given size.
webviewer.RectOp{Size: f32.Point{X: 100, Y: 100}}.Add(gtx.Ops)

// Stop using this webview in the current frame.
wv.Pop(gtx.Ops)
```

### Operations:

Operations must be added with `.Add(gtx.Ops)` method, except `webviewer.WebViewOp` which must be added
with `.Push(gtx.Ops)` and `.Pop(gtx.Ops)`. All operations (except `webviewer.WebViewOp`) should be added inside the
stack of `webviewer.WebViewOp`.

The operation will be executed at the end of the frame.

**WebView:**

- `webviewer.WebViewOp{Tag: &something}`:
    - Limits operations to the WebView. You must re-use the same `Tag` to use the same WebView. Also,
      some events (like `NavigationEvent`, `ErrorEvent`...) are sent using the `Tag`.
- `webviewer.NavigateOp{URL: "https://example.com"}`:
    - Load the given URL.
- `webviewer.RectOp{Size: f32.Point{X: 100, Y: 100}}`:
    - Display the WebView in the current frame, with the specified size in pixels.
- `webviewer.OffsetOp{Offset: f32.Point{X: 10, Y: 10}}`:
    - Set the position of the WebView in the current frame, with the specified offset in pixels.

**Storage:**

- `webviewer.SetStorageCmd{Local: webviewer.StorageTypeLocal, Storage: webviewer.Storage{Key: "key", Value: "value"}}`:
    - Set storage data for the WebView.
- `webviewer.ListStorageCmd{Tag: &something, Local: webviewer.StorageTypeLocal}`:
    - List the storage for the WebView. The response will be sent using the `Tag` as `StorageEvent`.
- `webviewer.RemoveStorageCmd{Tag: &something, Local: webviewer.StorageTypeLocal, Content: webviewer.Storage{Key: "key"}}`:
    - Remove the storage for the WebView.

**Cookies:**

- `webviewer.SetCookieCmd{Cookie: webviewer.Cookie{Name: "name", Value: "value"}}`:
    - Set the cookie data for the WebView.
- `webviewer.ListCookieCmd{Tag: &something}`:
    - List the cookie for the WebView. The response will be sent using the `Tag` as `CookiesEvent`.
- `webviewer.RemoveCookieCmd{Tag: &something, Cookie: webviewer.Cookie{Name: "name"}}`:
    - Remove the cookie for the WebView.

**Cache:**

- `webviewer.ClearCacheCmd{}`:
  - Clear Cache, Cookies, LocalStorage, SessionStorage, WebSQL, IndexedDB of the current WebView.

**Javascript:**

- `webviewer.JavascriptCmd{Javascript: "console.log(\"Hello World\")"}`:
    - Execute the given Javascript code in the WebView.
- `webviewer.InstallJavascriptCmd{Javascript: "console.log(\"Hello World\")"}`:
    - Persistently define the given Javascript code in the WebView. This code will be executed in every
      page load.
- `webviewer.MessageReceiverCmd{Tag: &something, Name: "your_function_name"}`:
    - Defines an function which can be called from the WebView, the function can be called
      as `window.callback.your_function_name("some text")` on the WebView side. The response will be sent using
      the `Tag` as `MessageEvent`.

### Options:

Some settings can be used to customize the WebView implementation, they must be set before the
first `webviewer.WebViewOp`, it will affect all WebViews. It's recommended to set them in the `init`.

- `webviewer.SetProxy()`:
    - [Android/Windows] Set the proxy to use for the WebView.
- `webviewer.SetCustomCertificates()`:
    - [Android/Windows] Set the custom certificates to use for the WebView, which will be used to validate the SSL
      connections, additionally to the system CA certificates.

## Events

Events are response sent using the `Tag` and should be handled with `gioplugins.Event()`.

- `webviewer.NavigationEvent`:
    - Sent to `WebViewOp.Tag` when the WebView navigates to a new page.
- `webviewer.TittleEvent`:
    - Sent to `WebViewOp.Tag` when the WebView changes the page title.
- `webviewer.ErrorEvent`:
    - Sent to `WebViewOp.Tag` when the WebView encounters an error.

- `webviewer.StorageEvent`:
    - Sent to `ListStorageOp.Tag` when the WebView returns the storage list.
- `webiver.CookiesEvent`:
    - Sent to `ListCookieOp.Tag` when the WebView returns the cookie list.

- `webiver.MessageEvent`:
    - Sent to `MessageReceiverOp.Tag` when the WebView calls a function defined with `MessageReceiverOp`.

--------------

## Features

We are capable of more than just displaying one webpage.

| Features | Windows | Android | MacOS | iOS | WebAssembly |
| -- | -- | -- | -- | -- | -- |
| Basic Support |✔|✔|✔|✔|✔****|
| Setup: Custom Proxy |✔***|✔***|❌|❌|❌|
| Setup: Custom Certificate |✔***|✔***|❌|❌|❌|
| Cookies: Read |✔|✔|✔*|✔|❌|
| Cookies: Write |✔|✔|✔*|✔|❌|
| Cookies: Delete |✔|✔|❌*|✔|❌|
| LocalStorage: Read |✔|✔|✔|✔|❌|
| LocalStorage: Write |✔|✔|✔|✔|❌|
| LocalStorage: Delete |✔|✔|✔|✔|❌|
| SessionStorage: Write |✔|✔|✔|✔|❌|
| SessionStorage: Read |✔|✔|✔|✔|❌|
| SessionStorage: Delete |✔|✔|✔|✔|❌|
| Cache: Delete |✔|✔|✔|✔|❌|
| Javascript: Execute |✔|✔|✔|✔|❌|
| Javascript: Install |✔|✔|✔|✔|❌|
| Javascript: Callback |✔**|✔**|✔**|✔**|❌|
| Events: NavigationChange |✔|✔|✔|✔|❌|
| Events: TitleChange |✔|✔|✔|✔|❌|

- ❌ = Not supported.
- ✔ = Supported.

- \* = Cookies can be shared across multiple instances of the WebView. Information from the cookie can be incomplete and
  lack metadata.
- ** = Only accepts a string as argument (other types are not supported and might be encoded as text).
- *** = Must be defined before the WebView is created and is shared with all instances.
- **** = Only websites that accepts iframe is supported. 
 
# APIs

Each operating system has uniqueAPI. For Windows 10+, we use WebView2. For Android 6+, we use WebView. For MacOS and
iOS, we use WKWebView. For WebAssembly, the HTMLIFrameElement is used.

## Requirements

- Windows:
    - End-Users: must have Windows 7+ and WebView2 installed (you can install it on the user's machine using
      the `installview` package).
    - Developers: must have Golang 1.18+ installed (no CGO required).
- WebAssembly:
    - End-Users: must have WebAssembly enabled browser (usually Safari 13+, Chrome 70+).
    - Developers: must have Golang 1.18+ installed (no CGO required).
    - Contributors: must have InkWasm installed.
- macOS:
    - End-Users: must have macOS 11+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- iOS:
    - End-Users: must have macOS 11+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- Android:
    - End-Users: must have Android 6 or later.
    - Developers: must have Golang 1.18+, OpenJDK 1.8, Android NDK, Android SDK 31+
      installed ([here for more information](https://gioui.org/doc/install/android)).
    - Hacking: must have Android SDK 30 installed.

## Limitations

1. Currently, GioWebview is always the top-most view/window and can't be overlapped by any other draw operation in Gio.
2. Render multiple webviews at the same time might cause unexpected behaviour, related to z-indexes.
3. On Javascript/WebAssembly, it needs to be allowed to iframe the content, which most websites blocks such operation.
4. It's not possible to use WebView using custom shapes (e.g. rounded corners) or apply transformations (e.g. rotating).
5. Some dialogs (such as alerts, file-picker, etc.) may not be displayed correctly, or not displayed at all.