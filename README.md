Gio-Plugins
-------

Gio plugins is a system to use _third-party_ plugins similar to Gio features, this package also holds one collection
of plugins for the [Gio](https://gioui.org). All plugins uses the same interface, and just one line of code must be
added in the main-loop, as explained below. _Furthermore, some packages can be used as standalone and doesn't require
the `plugin` package._

That package also serves as experimentation ground for "plugin system", which may land to Gio in the future.

### Motivation

Gio is a great GUI library, but it lacks some features that are important for some applications. I worked in some
extensions before, and I also maintained one fork of Gio, which contains some changes. However, that leads to many
issues: keeping the fork updated is problematic. Creating extensions that are similar to Gio is impossible.

I decided to create a plugin system to make it easier to use _third-party_
extensions. [I also proposed this idea to the Gio-core](https://lists.sr.ht/~eliasnaur/gio/%3Cfe3835f7-b4d4-4db9-81fb-dfd8ab06f2ed%40www.fastmail.com%3E)
.

### Usage

First, you must download the `plugin` package:

```bash
go get -u github.com/gioui-plugins/gio-plugins@latest
```

Now, you need to modify your event-loop, you must include `gioplugins.Hijack(window)` in your event-loop, before handling
events:

```diff
window := &app.Window{} // Gio window

for { 
+   evt := gioplugins.Hijack(window) // Gio main event loop

    switch evt := evt.(type) {
        // ...
    }
}
```

Each plugin has its own README.md file, explaining how to use it. In general, you can simple
use `nameOfPlugin.SomeOp{}.Add(gtx.Ops)`, similar of how you use `pointer.PassOp{}.Add(gtx.Ops)`, native 
from Gio. Beginning with Gio 0.6, it also introduces `Command`, which is executed using `gtx.Execute`, you can
also use Commands, if the plugin supports it, for instance `gioplugins.Execute(gtx, gioshare.TextCmd{Text: "Hello, World!"})`, 
similar to what you are familiar with `gtx.Execute(clipboard.WriteCmd{Text: "Hello, World!"})`.

Once `gioplugins.Event` is set, you can use the plugins as simple `op` operations or `cmd` commands. If you are unsure 
if the plugin is working, you can use the `pingpong` package, which will return on `PongEvent` to the given Tag:

```go
pingpong.PingOp{Tag: &something}.Add(gtx.Ops)

gioplugins.Execute(gtx, pingpong.PingCmd{Tag: &something})
```

> Note: It's uses `gioplugins.Execute` and not `gtx.Execute`!

You can receive responses using the `Tag`, similar Gio-core operations:

```go
for {
	evt, ok := gioplugins.Event(pingpong.Filter{Tag: &something})
	if !ok {
		break
    } 
	
    if evt, ok := evt.(pingpong.PongEvent); ok {
        fmt.Println(evt.Pong)
    } 
}
```

Of course, `pingpong` has no use in real-world applications, but it can be used to test if the plugin is working.

## Plugins

**We have few plugins available:**

| Name                                                                              | Description                                                                      | OS                                                          |
|-----------------------------------------------------------------------------------|----------------------------------------------------------------------------------|-------------------------------------------------------------|
| **[PingPong](https://github.com/gioui-plugins/gio-plugins/tree/main/pingpong)**   | Test if the plugin system is working.                                            | _Android, iOS, macOS, Windows, WebAssembly, Linux, FreeBSD_ |  
| **[Share](https://github.com/gioui-plugins/gio-plugins/tree/main/share)**         | Share text/links using the native share dialog.                                  | _Android, iOS, macOS, Windows, WebAssembly_                 |  
| **[WebViewer](https://github.com/gioui-plugins/gio-plugins/tree/main/webviewer)** | Display in-app webview using the native webview implementation on each platform. | _Android, iOS, macOS, Windows, WebAssembly_                 |  
| **[Hyperlink](https://github.com/gioui-plugins/gio-plugins/tree/main/hyperlink)** | Open hyperlinks in the default browser.                                          | _Android, iOS, macOS, Windows, WebAssembly_                 |  
| **[Explorer](https://github.com/gioui-plugins/gio-plugins/tree/main/explorer)**   | Opens the native file-dialog, to read/write files.                               | _Android, iOS, macOS, Windows, WebAssembly_                 |  
| **[Safedata](https://github.com/gioui-plugins/gio-plugins/tree/main/safedata)**   | Read/Write files into the secure storage of the device.                          | _Android, iOS, macOS, Windows, WebAssembly_                 |

**We have few plugins planned:**

Some plugins are planned, but not yet implemented, follow the development at https://github.com/orgs/gioui-plugins/projects/1. Also, 
consider send some 👍 on issues which mentions features that you like.

If you want to help, please open an issue or a PR! If you want to suggest a plugin, please open an issue.

-----------

### Creating a new plugin

If you want to create a new plugin, you can check the `pingpog` package, which is the simplest plugin available.
Generally, you need implement `plugin.Handler`, and call `plugin.Register` in your `init()` function. Your code will get
specific events and ops, which you define in your `Handler` implementation.

#### Limitations

- Android: XML/Manifest: There's no direct integration with `gogio`. Consequently, your plugin cannot require any
  additional XML file or changes in the manifest. That may limit some plugins, but it's a limitation of `gogio` itself.
- General: Position: There's no way to get relative position of each operation, or mimic the `paint.PaintOp{}`, so if
  your plugin is adding views to the screen, you must use absolute positions. This is a limitation of Gio itself, which
  doesn't easily expose the relative position of each operation.

### Compatibility

Most packages are compatible with the latest version of Gio, and only the latest version should be supported. Beware
that internal changes in the Gio API can break the compatibility with `plugin`. None of the packages has stable API,
and breaking changes can happen at any time.

Most plugins are compatible with Android 5+, iOS 13+, MacOS 12+, Windows 10+ and WebAssembly. Currently, we consider
the Windows and Android as high-priority, and the WebAssembly as medium-priority, MacOS and iOS as low-priority.
Furthermore, we don't have any plans to support Linux and FreeBSD due to the low market-share and the lack of API
standards.

| Priority | OS          | Arch         |
|----------|-------------|--------------|
| High     | Windows     | AMD64        |
| High     | Android     | ARM64, ARMv7 |
| Medium   | WebAssembly | WASMv1       |
| Low      | MacOS       | ARM64, AMD64 |
| Low      | iOS         | ARM64        |
| Ignored  | FreeBSD     | -            |
| Ignored  | Linux       | -            |

Currently, some package might not work on some platforms, and some features might not be available on some platforms.

### Security

This package heavily uses `unsafe`, and as it suggest: it can be unsafe to use. We are not responsible for any damage
caused by this package. Some plugins also use `unsafe` and CGO to interact with the native platform. While we try to
keep the code safe, we can't guarantee that it is safe enough for your use-case. Also, is impossible to verify the
integrity of native-APIs, so we can't guarantee that the native-APIs will have the expected behavior.

### Testing

Since we have limited resources and devices, we can't test all plugins on all platforms and devices. Currently, we have
a few devices available and with limited range of OS versions. Plugins are usually tested on those devices:

- **Android**: Motorola Droid Max, Motorola E6, Xiaomi A7, Xiaomi Note 9, (+ BrowserStack)
- **iOS**: iPhone SE 2Gen (2020), (+ BrowserStack)
- **WASM**: Chrome, Firefox, Safari, (+ BrowserStack)
- **Windows**: Custom (Ryzen 3900X+RX 5700XT), Proxmox VM (EPYC 7501P)
- **MacOS**: MacBook Air (M1, 2020)
- **Linux**: Proxmox VM (EPYC 7501P)
- **FreeBSD**: Proxmox VM (EPYC 7501P)

Tests are performed manually, since most features interact with the native platform, and automated tests are not
easy to implement. We are open to suggestions on how to improve the testing process.

Please, if you find any bug, open an issue or a PR!

### Hacking

Each OS has its own way to interact with the native APIs. For example, on Android, you need to create a Java class and
call it using JNI. On iOS and MacOS you need to write some Objective-C code and call it using CGO. On Windows, you need
to write some code using `syscall` to each DLL, which may use COM API, some APIs uses WinRT instead of Win32, which can
be harder to use. On WebAssembly, you need to use `syscall/js` to interact with the browser APIs, or use InkWasm, which
is faster. On Linux/FreeBSD you may need to use C and CGO to interact with the native APIs.

### License

This package is licensed under the MIT License, some pre-compiled files may have other license. See
the [LICENSE](LICENSE) file for details.
