# share

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/share.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/share)

Opens the share dialog to share texts and websites.

--------------

## Usage

### Operations:

Operations must be added with `.Add(gtx.Ops)` method. The operation will be executed at the end of the frame.

- `share.TextOp`:
  - Shares text using the native share dialog.
- `share.WebsiteOp`:
  - Shares a website using the native share dialog.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `webviewer.ErrorEvent`:
  - Sent to `Tag` when it's not possible to share.

## Features

| OS | Windows | Android | MacOS | iOS | WebAssembly |
| -- | -- | -- | -- | -- | -- |
| Share Text |✔|✔|✔|✔|✔|
| Share Website |✔|✔|✔|✔|✔|
| Share Image |❌|❌|❌|❌|❌|
| -- | -- | -- | -- | -- | -- |
| API | [DataTransfer](https://learn.microsoft.com/en-us/uwp/api/windows.applicationmodel.datatransfer.datatransfermanager?view=winrt-22621) | [Intent.ACTION_SEND](https://developer.android.com/training/sharing/send) | [NSSharingServicePicker](https://developer.apple.com/documentation/appkit/nssharingservicepicker) | [UIActivityViewController](https://developer.apple.com/documentation/uikit/uiactivityviewcontroller?language=objc) | [Web Share API](https://developer.mozilla.org/en-US/docs/Web/API/Web_Share_API) |

- ❌ = Not supported.
- ✔ = Supported.

## Requirements

- Windows:
    - End-Users: must have Windows 10+.
    - Developers: must have Golang 1.18+ installed (no CGO required).
- WebAssembly:
    - End-Users: must have WebAssembly enabled browser (usually Safari 13+, Chrome 70+).
    - Developers: must have Golang 1.18+ installed (no CGO required).
- macOS:
    - End-Users: must have macOS 12+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- iOS:
    - End-Users: must have macOS 13+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- Android:
    - End-Users: must have Android 6+.
    - Developers: must have Golang 1.18+, OpenJDK 1.8, Android NDK, Android SDK 31 installed ([here for more information](https://gioui.org/doc/install/android)).

