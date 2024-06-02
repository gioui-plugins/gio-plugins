# safedata

[![Go Reference](https://pkg.go.dev/badge/github.com/gioui-plugins/gio-plugins/safedata.svg)](https://pkg.go.dev/github.com/gioui-plugins/gio-plugins/safedata)

Store credentials (and arbitrary data) securely and persistent. This is useful to store tokens
and other types of access credentials.

--------------

## Usage

This package can be used as standalone (with/without Gio) and also as plugin for Gio.

## Using as Standalone:

```go
config := Config{
App: "MyApp"
// ...
}

sh := safedata.NewSafeData(config)

secret := safedata.Secret{
Identifier:  "AccessToken",
Description: "AccessToken for MyApp",
Data:        []byte{0xDE, 0xAD, 0xBE, 0xEF},
}

// Inserting/updating data:
if err := sh.Set(secret); err != nil {
// ...
}

// Retrieve data:
myToken, err := sh.Get("AccessToken")
if err != nil {
// ...
}
```

Note: `safedata.Config` varies for each OS, and you should create `yourfile_{os}.go` for each
supported OS. If you are using Gio, it's also possible to create one `safedata.Config`
using `giosafedata.NewConfigFromViewEvent`.

## Using as Gio-Plugin:

### Operations:

Operations must be added with `gtx.Execute` method. The operation will be executed at the end of the frame.

- `giosafedata.WriteSecretCmd`:
    - Writes a Secret.
- `giosafedata.ReadSecretCmd`:
    - Reads a Secret using the provided Identifier, the response is sent to the given Tag.
- `giosafedata.DeleteSecretCmd`:
    - Deletes a Secret using the provided Identifier.
- `giosafedata.ListSecretCmd`:
    - List all Secret which belongs to the current app.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `giosafedata.ErrorEvent`:
    - Sent to `Tag` when it's not possible to write/read/list/delete.
- `giosafedata.SecretsEvent`:
    - Sent to `Tag` as response from `ReadSecretCmd` or `ListSecretCmd`.

## Features

| OS     | Windows                                                                 | Android                                                                      | MacOS                                                                                                   | iOS                                                                                                     | WebAssembly                                                                          |
|--------|-------------------------------------------------------------------------|------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| Write  | ✔                                                                       | ✔                                                                            | ✔                                                                                                       | ✔                                                                                                       | ✔                                                                                    |
| Read   | ✔                                                                       | ✔                                                                            | ✔                                                                                                       | ✔                                                                                                       | ✔                                                                                    |
| List   | ✔                                                                       | ✔                                                                            | ✔                                                                                                       | ✔                                                                                                       | ✔                                                                                    |
| Delete | ✔                                                                       | ✔                                                                            | ✔                                                                                                       | ✔                                                                                                       | ✔                                                                                    |
| API    | [WinCred](https://learn.microsoft.com/en-us/windows/win32/api/wincred/) | [Android Keystore](https://developer.android.com/training/articles/keystore) | [Keychain Services](https://developer.apple.com/documentation/security/keychain_services?language=objc) | [Keychain Services](https://developer.apple.com/documentation/security/keychain_services?language=objc) | [LocalStorage](https://developer.mozilla.org/pt-BR/docs/Web/API/Window/localStorage) |

- ❌ = Not supported.
- ✔ = Supported.

## Security/Notes

This package uses what is available on the OS to safely store
credentials and any sensible data. However, not all OSes provides
such function or have heavily limitations. That is the list of
known issues/vulnerability:

- [Darwin] You must sign your app.
- [Darwin] Credentials are visible cross-application after user authorization.
- [Android] Credentials may lost after app uninstall or update.
- [Windows] Large data is split into multiple credentials, due to maximum size for each credential.
- [Windows] Credential Storage have a very low capacity, preventing from storing large data or too many credentials.
- [Windows] Credentials are visible cross-application, without restrictions.
- [WebAssembly] Credentials are visible to any script in the page, which is vulnerable to XSS.
- [WebAssembly] Credentials may lose after cache clear (Clear-Site-Data header or invoked by the end-user).

- This package doesn't check the integrity of the data (you should add your own checksum).
- Credentials can be modified or deleted externally (usually on device settings and similar).

## Background

Since it's a security-related package, I'm listing how it works behind the scenes.

- [Android] It creates files on the folder specified by Config, that file is encrypted
  using Android KeyStore. The IV/Nonce is stored into the file. It uses AES-CBC as encryption,
  since we don't guarantee integrity on any OS.
- [Darwin] It stores the data into Keychain, directly.
- [Windows] It creates new credentials using WinCred, as Generic Credentials. Each credential
  supports upto 512*5 bytes. If the data exceeds the maximum size, new credentials are created
  using Blake2 derivation for names, starting from index 0, if the data is larger.
- [WebAssembly] It stores the data into LocalStorage, directly.

## Requirements

- Windows:
    - End-Users: must have Windows 7+.
    - Developers: must have Golang 1.18+ installed (no CGO required).
- WebAssembly:
    - End-Users: must have WebAssembly enabled browser (usually Safari 13+, Chrome 70+).
    - Developers: must have Golang 1.18+ installed (no CGO required).
- macOS:
    - End-Users: must have macOS 10+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- iOS:
    - End-Users: must have iOS 10+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- Android:
    - End-Users: must have Android 6+.
    - Developers: must have Golang 1.18+, OpenJDK 1.8, Android NDK, Android SDK 31
      installed ([here for more information](https://gioui.org/doc/install/android)).

