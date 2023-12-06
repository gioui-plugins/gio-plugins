# auth

[![Go Reference](https://pkg.go.dev/badge/github.com/gioui-plugins/gio-plugins/safedata.svg)](https://pkg.go.dev/github.com/gioui-plugins/gio-plugins/auth)

Brings "Sign in with Apple" and "Sign in with Google" to your Gio app.

--------------

## Setup

That package requires to register your app on each provider (Apple/Google) and need to define the appropriate
information (like AppID, AppName, RedirectURL and your signing key). Each provider has its own requirements,
and you should follow the instructions below. If you are not using this package as Plugin (you are not
using `gioauth.ListenOp`), then you need to call `ProcessCustomSchemeCallback` manually, when new URL is received.

It's out of scope of this package to explain how to register your app on each provider.

Notice: On macOS, ASAuthorizationAppleIDProvider is not available for Developer ID signed apps (apps signed to be
distributed outside AppStore). By default, we assume that you are using Developer ID signed apps, if you are using
AppStore signed apps, you should compile with `-tags appstore`.

### Apple

Apple doesn't allow to use Localhost or CustomScheme as RedirectURL, you must use a domain that you own. Then,
you MUST redirect the user using a custom URL Scheme to your app. For example, if you own `example.com`, you
should use `https://example.com/auth/apple` as RedirectURL and redirect the user to `example://auth/apple`.

Your custom URL Scheme must be registered using `-schemes` flag in `gogio`.

### Google

No special requirements are needed. The Scheme is the reverse of your AppID. For example, if your AppID is
`com.example.app`, your Scheme should be `app.example.com`. If you unsure about your Scheme, you can use:

```go
fmt.Println((&google.Provider{
WebClientID:     "your.id.apps.googleusercontent.com",
DesktopClientID: "your.id.apps.googleusercontent.com",
}).Scheme())
```

That will print the Scheme that you should use. Your custom URL Scheme must be registered using `-schemes` flag
in `gogio`.

## Security

Not all providers supports PKCE (I'm aim at you, Apple). Consider as recommended to use "Nonce" as additional security
measure, not just for replay attacks.

### Recommended Flow

If you are connecting to your own backend, for authentication, you should use the following flow:

1. Client: Generate a random byte sequence (>= 32 bytes)
    - You may need to store the random byte sequence, if you are on a web environment (JS/WASM).
    - You may want to set additional cookies to prevent CSRF and can combine that if this random-sequence.
2. Client: Creates a Hash (using any secure PRP/Hash) and encodes it.
3. Client: Define the `Nonce` as the pre-image of the hash (the value of step 1) and to the request (using Open function
   from Auth).
4. Client: Once the response received from the provider, along with the pre-image to your backend.
5. Server: Validates the OpenID Connect signature and the `Nonce` (using the pre-image).

That will prevent replay attacks, since the `Nonce` must be unique for each request. Also, it will prevent someone
from stealing the `id_token` and using it on your backend, since the `Nonce` must match if the provided pre-image.

Additionally, it's possible to combine other information into the `Nonce`, and use similarly to `State`.

## Using as Gio-Plugin:

### Setup:

Add your provider to the `gioauth.DefaultProviders` list. For example:

```go
	gioauth.DefaultProviders = []providers.Provider{
		&google.Provider{
			WebClientID:     "YOUR-CODE.apps.googleusercontent.com",
			DesktopClientID: "YOUR-CODE.apps.googleusercontent.com",
			RedirectURL:     "",
		},
		&apple.Provider{
			ServiceIdentifier: "YOUR-APP",
			RedirectURL:       "https://your-call-back.com/path/",
		},
	}
```

You also need to use `gogio` to setup deeplinking, see above.

### Operations:

Operations must be added with `.Add(gtx.Ops)` method. The operation will be executed at the end of the frame.

- `gioauth.RequestOp`:
    - Requests login with Google/Apple.
- `gioauth.ListenOp`:
    - Listen to events.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `gioauth.AuthEvent`:
    - Sent to `Tag` with the tokens received from the provider.
- `giosafedata.SecretsEvent`:
    - Sent to `Tag` as response from `ReadSecretOp` or `ListSecretOp`.

## Features

| OS     | Windows                | Android                                                                                                                                     | MacOS                                                                                                                                                                                                                                                                                                | iOS                                                                                                                                                      | WebAssembly |
|--------|------------------------|---------------------------------------------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------|-------------|
| Google | ✔ <br/> (CustomScheme) | ✔ <br/> ([gms.auth.api.identity](https://developers.google.com/android/reference/com/google/android/gms/auth/api/identity/package-summary)) | ✔ <br/> ([ASWebAuthenticationSession](https://developer.apple.com/documentation/authenticationservices/aswebauthenticationsession?language=objc) or [ASAuthorizationAppleIDProvider](https://developer.apple.com/documentation/authenticationservices/asauthorizationappleidprovider?language=objc)) | ✔ <br/> ([ASAuthorizationAppleIDProvider](https://developer.apple.com/documentation/authenticationservices/asauthorizationappleidprovider?language=objc) | ✔           |
| Apple  | ✔ <br/> (CustomScheme) | ✔ <br/> ([Custom Chrome Tabs](https://developer.android.com/reference/androidx/browser/customtabs/CustomTabsIntent))                        | ✔ <br/> ([ASAuthorizationAppleIDProvider](https://developer.apple.com/documentation/authenticationservices/asauthorizationappleidprovider?language=objc)                                                                                                                                             | ✔ <br/> ([ASAuthorizationAppleIDProvider](https://developer.apple.com/documentation/authenticationservices/asauthorizationappleidprovider?language=objc) | ✔           |

- ❌ = Not supported.
- ✔ = Supported.

## Requirements

- Windows:
    - End-Users: must have Windows 10+.
    - Developers: must have Golang 1.20+ installed (no CGO required).
- WebAssembly:
    - End-Users: must have WebAssembly enabled browser (usually Safari 13+, Chrome 70+).
    - Developers: must have Golang 1.20+ installed (no CGO required).
- macOS:
    - End-Users: must have macOS 11+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- iOS:
    - End-Users: must have iOS 14+.
    - Developers: must have macOS device with Golang, Xcode, and CLang installed.
- Android:
    - End-Users: must have Android 6+.
    - Developers: must have Golang 1.18+, OpenJDK 1.8, Android NDK, Android SDK 30
      installed ([here for more information](https://gioui.org/doc/install/android)).

