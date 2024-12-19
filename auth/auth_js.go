package auth

import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/auth/providers/apple"
	"github.com/gioui-plugins/gio-plugins/hyperlink"
	"net/url"
	"syscall/js"
)

type driver struct {
	hp *hyperlink.Hyperlink

	appleSuccess js.Func
	appleFailure js.Func
}

func attachDriver(house *Auth, config Config) {
	house.driver = driver{
		hp: hyperlink.NewHyperlink(hyperlink.Config{}),
		appleSuccess: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			if len(args) == 0 {
				return nil
			}

			authorization := args[0].Get("authorization")
			if authorization.IsNull() {
				return nil
			}

			code := authorization.Get("code").String()
			idToken := authorization.Get("id_token").String()

			house.sendResponse(AuthenticatedEvent{
				Provider: apple.IdentifierApple,
				Code:     code,
				IDToken:  idToken,
			})
			return nil
		}),
		appleFailure: js.FuncOf(func(this js.Value, args []js.Value) interface{} {
			return nil
		}),
	}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.hp.Configure(config.Config)
}

func (d *driver) open(provider providers.Provider, nonce string) error {
	switch provider := provider.(type) {
	case *apple.Provider:
		return d.openAppleSDK(provider, nonce)
	default:
		return d.openAny(provider, nonce)
	}
}

// Requires using external SDK
// <script type="text/javascript" src="https://appleid.cdn-apple.com/appleauth/static/jsapi/appleid/1/en_US/appleid.auth.js" async></script>
func (d *driver) openAppleSDK(provider *apple.Provider, nonce string) error {
	if !d.isAvailableAppleSDK() {
		return d.openAny(provider, nonce)
	}

	appleID := js.Global().Get("AppleID")
	obj := js.Global().Get("Object").New()
	obj.Set("clientId", provider.ClientID())
	if !provider.DisabledEmailAndName {
		obj.Set("scope", "email name")
	}
	obj.Set("redirectURI", provider.RedirectURL)
	obj.Set("nonce", nonce)
	obj.Set("state", provider.Identifier().State())
	obj.Set("usePopup", true)

	auth := appleID.Get("auth")
	auth.Call("init", obj)

	go auth.Call("signIn", nil).Call("then", d.appleSuccess).Call("catch", d.appleFailure)

	return nil
}

func (d *driver) isAvailableAppleSDK() bool {
	return js.Global().Get("AppleID").Truthy()
}

func (d *driver) openAny(provider providers.Provider, nonce string) error {
	u, err := url.Parse(provider.URL(nonce))
	if err != nil {
		return err
	}
	return d.hp.Open(u)
}
