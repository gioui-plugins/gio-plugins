# pushnotification

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/pushnotification.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/pushnotification)

Retrieves push notification tokens to allow sending push notifications, the notification will be displayed even if the
app is not running, and the app can be opened by clicking on the notification. Each OS have their own provider (FCM for
Android, APNs for iOS/macOS, Web Push for WebAssembly), but the API is unified across all platforms. Your server needs
to distinguish between the different providers and send the notification to the correct endpoint.

It requires an external server to send the push notifications. That is a major difference from Gio-X Notify, which only
allows sending notifications when the app is running.

> [!TIP]
> Do NOT request the token on every app start, as that may cause the permission prompt to be shown at inappropriate
> times.

> [!IMPORTANT]
> The WebAssembly version requires a service worker to be registered, the gio-cmd2 will generate a service worker for
> you, but you need to copy that file to your web server (similar to how you copy .wasm and .js files).

> [!NOTE]
> Be sure to validate the endpoint before sending push notifications, otherwise your server is prone to SSRF attacks.

--------------

## Before you start

Android: you need to create a Firebase project and download the `google-services.json` file, and set `Config` based on
`google-services.json` information.

macOS/iOS: you need to create a Provisioning Profile with Push Notifications capability and use `gogio` with the same
profile.

WebAssembly: you need to register a service worker and need to use HTTPS to serve your web app. You also need to
generate a `VAPID` and provide it in the Config. The Service Worker expect the data (which is pushed by the server) to
be in "Declarative Web Push Format" (see https://notifications.spec.whatwg.org/#dictdef-notificationoptions).

Windows: you need to bundle it as MSIX and register it in Azure platform and pray for the best. It might not work.

## Usage

### Freestanding

- `pushnotification.NewPush`:
    - Creates a instance of Push struct, given the config.
- `pushnotification.Configure`:
    - Updates the current Push with the given config.
- `pushnotification.RequestToken`:
    - Synchronously requests the push token from the provider (FCM, APNs, Web Push), it also requests the permission
      if needed.

### Gio

#### Plugin:

You need to populate the `giopushnotification.DefaultProviders` with Firebase configuration, that is required for
Android, but can be defined on any platform.

To request a push token, you can use `giopushnotification.RequestTokenCmd`, it will return a
`giopushnotification.TokenReceivedEvent` with the token. It will also request the permission for push notifications
if needed.

## Features

| OS                  | Windows | Android                                                                | MacOS                                                                                   | iOS                                                                                     | WebAssembly                                                                 |
|---------------------|---------|------------------------------------------------------------------------|-----------------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------|-----------------------------------------------------------------------------|
| Get Token           | Testing | ✔                                                                      | ✔                                                                                       | ✔                                                                                       | ✔¹                                                                          |
| Subscribe to Topics | ❌       | ❌                                                                      | ❌                                                                                       | ❌                                                                                       | ❌                                                                           |
| API                 | --      | [Firebase Messaging](https://firebase.google.com/docs/cloud-messaging) | [UNUserNotificationCenter](https://developer.apple.com/documentation/usernotifications) | [UNUserNotificationCenter](https://developer.apple.com/documentation/usernotifications) | [PushManager](https://developer.mozilla.org/en-US/docs/Web/API/PushManager) |

- ❌ = Not supported (Yet).
- ✔ = Supported.

## Some Server-Side Examples (Go)

You need to use a external server (not included in this repo) to send push notifications.

### Android (FCM) & iOS (APNs via FCM)

Using [firebase-admin-go](https://firebase.google.com/docs/admin/setup#go).

```go
package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/option"
)

func sendFCM(token string) {
	ctx := context.Background()
	opt := option.WithCredentialsFile("path/to/serviceAccountKey.json")
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	message := &messaging.Message{
		Data: map[string]string{
			"deeplink": "https://example.com/action",
		},
		Notification: &messaging.Notification{
			Title: "Hello",
			Body:  "World",
		},
		Token: token,
	}

	response, err := client.Send(ctx, message)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Successfully sent message:", response)
}
```

### Web Push

Using [SherClockHolmes/webpush-go](https://github.com/SherClockHolmes/webpush-go).

```go
package main

import (
	"encoding/json"

	"github.com/SherClockHolmes/webpush-go"
)

func main() {
	subscription := webpush.Subscription{}
	if err := json.Unmarshal([]byte(`*** THE TOKEN VALUE, WHICH YOU GET FROM GIO-PLUGINS *** `), &subscription); err != nil {
		panic(err)
	}

	// The data is in "Declarative Web Push Format"
	data := `{
    "web_push": 8030,
    "notification": {
        "title": "Your Web Push",
        "body": "It's working'!",
        "navigate": "https://example.com/",
        "silent": false
    }`

	sendWebPush(&subscription, data)
}

func sendWebPush(s *webpush.Subscription, c []byte) {
	// You need to generate VAPID keys first
	vapidPublicKey := "Base64 encoded public key"
	vapidPrivateKey := "Base64 encoded private key"

	resp, err := webpush.SendNotification(c, s, &webpush.Options{
		Subscriber:      "https://127.0.0.1",
		VAPIDPublicKey:  vapidPublicKey,
		VAPIDPrivateKey: vapidPrivateKey,
		TTL:             120,
	})
	if err != nil {
		panic(err)
	}
	res, err := httputil.DumpResponse(resp, true)
	fmt.Println(string(res), err)
}
```