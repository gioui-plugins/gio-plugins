# inapppay

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/inapppay.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/inapppay)

Launches the In-App purchase flow.

This plugins makes possible to offer in-app purchases to your users, relies on
the [In-App Billing API](https://developer.android.com/google/play/billing/billing_reference)
and [Apple's App Store Connect](https://developer.apple.com/in-app-purchase/).

--------------

## Usage

### Freestanding

- `inapppay.NewInAppPay`:
    - Creates a instance of the InAppPay plugin.
- `inapppay.Configure`:
    - Configures the plugin.
- `inapppay.ListProducts`:
    - Lists the available products for purchase.
- `inapppay.Purchase`:
    - Launches the purchase flow for a specific product. You can attach additional data to the purchase, to identify the
      user in the backend.
- `pushnotification.ListMemeberships`:
    - Lists the available subscriptions for purchase.
- `pushnotification.Subscribe`:
    - Launches the subscription flow for a specific subscription. You can attach additional data to the subscription, to
      identify the user in the backend.

### Gio

#### Plugin:

You can get products by `gioinapppay.ListProductsCmd` and subscriptions by `gioinapppay.ListSubscriptionsCmd`.
The result is sent as `gioinapppay.PurchaseResultEvent, `gioinapppay.SubscribeResultEvent`

You can initiate a purchase by `gioinapppay.PurchaseCmd` and subscribe by `gioinapppay.SubscribeCmd`.
The result is sent as `gioinapppay.ProductDetailsEvent` and `gioinapppay.SubscriptionDetailsEvent`.

## Features

| OS                            | Windows | Android | MacOS | iOS | WebAssembly |
|-------------------------------|---------|---------|-------|-----|-------------|
| List One-Purchase Offers      | ❌       | ✔       | ?     | ✔   | ❌           |
| Display payment checkout      | ❌       | ✔       | ?     | ✔   | ❌           |
| List Subscriptions Offers     | ❌       | ❌       | ❌     | ❌   | ❌           |
| Display subscription checkout | ❌       | ❌       | ❌     | ❌   | ❌           |

- ❌ = Not supported (Yet).
- ✔ = Supported.