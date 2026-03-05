//go:build js && wasm

package pushnotification

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"syscall/js"
)

type driver struct {
	config Config
	push   *Push
}

func attachDriver(push *Push, config Config) {
	d := &driver{push: push}
	push.driver = d
	configureDriver(d, config)
}

func configureDriver(d *driver, config Config) {
	d.config = config
}

func (d *driver) requestToken() (Token, error) {
	if d.config.VAPIDPublicKey == "" {
		return Token{}, errors.New("missing VAPID public key")
	}

	window := js.Global()
	navigator := window.Get("navigator")
	if navigator.IsUndefined() {
		return Token{}, ErrNotAvailable
	}

	serviceWorker := navigator.Get("serviceWorker")
	if serviceWorker.IsUndefined() {
		return Token{}, ErrNotAvailable
	}

	result := make(chan struct {
		Token Token
		Error error
	}, 1)

	var (
		success    js.Func
		subSuccess js.Func
		fail       js.Func
	)

	go func() {
		// Prepare options
		options := map[string]interface{}{
			"userVisibleOnly":      true,
			"applicationServerKey": d.config.VAPIDPublicKey,
		}

		// We assume a SW is registered.
		// We need to get the registration first.

		promise := serviceWorker.Call("getRegistration")

		fail = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			err := args[0]
			result <- struct {
				Token Token
				Error error
			}{
				Token: Token{},
				Error: errors.New(err.Get("message").String()),
			}
			return nil
		})

		success = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			registration := args[0]
			if registration.IsUndefined() || registration.IsNull() {
				result <- struct {
					Token Token
					Error error
				}{
					Token: Token{},
					Error: errors.New("no service worker registration found"),
				}
				return nil
			}

			pushManager := registration.Get("pushManager")
			if pushManager.IsUndefined() {
				result <- struct {
					Token Token
					Error error
				}{
					Token: Token{},
					Error: errors.New("push manager not available in service worker registration"),
				}
				return nil
			}

			subPromise := pushManager.Call("subscribe", options)

			subSuccess = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				subscription := args[0]

				auth := js.Global().Get("Uint8Array").New(subscription.Get("getKey").Call("call", subscription, "auth"))
				p256dh := js.Global().Get("Uint8Array").New(subscription.Get("getKey").Call("call", subscription, "p256dh"))

				authBytes := make([]byte, auth.Get("byteLength").Int())
				js.CopyBytesToGo(authBytes, auth)

				p256dhBytes := make([]byte, p256dh.Get("byteLength").Int())
				js.CopyBytesToGo(p256dhBytes, p256dh)

				webtoken := WebPushToken{
					Endpoint: subscription.Get("endpoint").String(),
					Keys: WebPushKeys{
						Auth:   base64.RawURLEncoding.EncodeToString(authBytes),
						P256dh: base64.RawURLEncoding.EncodeToString(p256dhBytes),
					},
				}

				webpushJson, err := json.Marshal(webtoken)

				result <- struct {
					Token Token
					Error error
				}{
					Token: Token{
						Token:    string(webpushJson),
						Platform: PlatformWeb,
					},
					Error: err,
				}
				return nil
			})

			subPromise.Call("then", subSuccess).Call("catch", fail)

			return nil
		})

		promise.Call("then", success).Call("catch", fail)
	}()

	r := <-result

	success.Release()
	subSuccess.Release()
	fail.Release()

	return r.Token, r.Error
}
