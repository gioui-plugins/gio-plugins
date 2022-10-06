# hyperlink

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/share.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/share)

Open a hyperlink in the default browser.

--------------

## Usage

To open one link, you can use the `Open` operation.

That will open the link in the default browser. You can use `OpenURL` to open a `*url.URL`:

```go
hyperlink.OpenOp{URI: &url.URL{
    Scheme: "https",
    Host:   "github.com",
}}.Add(gtx.Ops)
```

### Operations:

Operations must be added with `.Add(gtx.Ops)` method. The operation will be executed at the end of the frame.

- `share.OpenOp`:
    - Opens the link in the default browser. Currently, only supports `http` and `https` schemes.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `webviewer.ErrorEvent`:
    - Sent to `Tag` when it's not possible to open the hyperlink.

## Features

| Features | Windows | Android | MacOS | iOS | WebAssembly | FreeBSD |  Linux |
| -- | -- | -- | -- | -- | -- |  -- |  -- |
| HTTP |✔|✔|✔|✔|✔|✔|✔|
| HTTPS |✔|✔|✔|✔|✔|✔|✔|


By default, only HTTP and HTTPS links are allowed, but you can change that by changing `InsecureIgnoreScheme` to `true`,
you should validate the URL and scheme on your own.