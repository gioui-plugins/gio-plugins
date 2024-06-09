# hyperlink

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/share.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/share)

Open a hyperlink in the default browser.

--------------

## Usage

## Freestanding

- `hyperlink.NewShare`:
  - Creates a instance of Share struct, given the config.
- `hyperlink.Configure`:
  - Updates the current Share with the given config.
- `hyperlink.Open`:
  - Opens the link in the default browser. 
  - By default, only supports `http` and `https` schemes.

## Gio

## Non-Plugin:

If you want to use it without plugin, read the Freestanding instructions. We provide some helper functions, such as NewConfigFromViewEvent and such.

To open one link, you can use the `Open` operation.

That will open the link in the default browser. You can use `OpenURL` to open a `*url.URL`:

```go
gioplugins.Execute(gtx, giohyperlink.OpenOp{URI: &url.URL{
    Scheme: "https",
    Host:   "github.com",
}})
```

### Operations:

Operations must be added with `gtx.Execute` method. The operation will be executed at the end of the frame.

- `giohyperlink.OpenCmd`:
    - Opens the link in the default browser. Currently, only supports `http` and `https` schemes.

## Events:

Events are response sent using the `Tag` and should be handled with `gioplugins.Event()`.

- `giohyperlink.ErrorEvent`:
    - Sent to `Tag` when it's not possible to open the hyperlink.

## Features

| Features | Windows | Android | MacOS | iOS | WebAssembly | FreeBSD |  Linux |
| -- | -- | -- | -- | -- | -- |  -- |  -- |
| HTTP |✔|✔|✔|✔|✔|✔|✔|
| HTTPS |✔|✔|✔|✔|✔|✔|✔|


By default, only HTTP and HTTPS links are allowed, but you can change that by changing `InsecureIgnoreScheme` to `true`,
you should validate the URL and scheme on your own.