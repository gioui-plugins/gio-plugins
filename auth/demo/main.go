package main

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/clipboard"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/widget"
	"github.com/gioui-plugins/gio-plugins/auth"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout"
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/auth/providers/apple"
	"github.com/gioui-plugins/gio-plugins/auth/providers/google"
	"github.com/gioui-plugins/gio-plugins/plugin/gioplugins"
	"image"
	"image/color"
	"io"
	"strings"
)

var sharper = text.NewShaper()

func main() {
	w := &app.Window{}
	ops := new(op.Ops)

	var last auth.AuthenticatedEvent

	var buttonGoogle = authlayout.GoogleDummyButton{
		ButtonStyle: authlayout.DefaultLightGoogleButtonStyle,
	}
	var buttonApple = authlayout.AppleDummyButton{
		ButtonStyle: authlayout.DefaultLightAppleButtonStyle,
	}

	go func() {
		for {
			evt := gioplugins.Hijack(w)

			switch evt := evt.(type) {
			case app.DestroyEvent:
				return
			case app.FrameEvent:
				gtx := app.NewContext(ops, evt)
				{
					s := clip.RRect{Rect: image.Rectangle{Max: gtx.Constraints.Max}}.Push(gtx.Ops)
					paint.ColorOp{Color: color.NRGBA{R: 240, G: 240, B: 240, A: 255}}.Add(gtx.Ops)
					paint.PaintOp{}.Add(gtx.Ops)
					s.Pop()
				}

				if buttonGoogle.Clicked(gtx) {
					gioplugins.Execute(gtx, gioauth.OpenCmd{Provider: google.IdentifierGoogle, Nonce: nonce()})
				}

				if buttonApple.Clicked(gtx) {
					gioplugins.Execute(gtx, gioauth.OpenCmd{Provider: apple.IdentifierApple, Nonce: nonce()})
				}

				gtx.Constraints.Max.X -= gtx.Dp(20)
				op.Offset(image.Pt(10, 10)).Push(gtx.Ops)
				buttonGoogle.Layout(gtx)

				op.Offset(image.Pt(0, gtx.Dp(60))).Push(gtx.Ops)
				buttonApple.Layout(gtx)

				c := op.Record(gtx.Ops)
				paint.ColorOp{Color: color.NRGBA{A: 255}}.Add(gtx.Ops)
				cc := c.Stop()

				op.Offset(image.Pt(0, gtx.Dp(60))).Push(gtx.Ops)
				widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, "Last", cc)

				for {
					evt, ok := gioplugins.Event(gtx, gioauth.Filter{})
					if !ok {
						break
					}

					switch evt := evt.(type) {
					case gioauth.AuthEvent:
						last = auth.AuthenticatedEvent(evt)
						gtx.Execute(clipboard.WriteCmd{Type: "text", Data: io.NopCloser(strings.NewReader(evt.IDToken))})
					}
				}

				op.Offset(image.Pt(0, gtx.Dp(60))).Push(gtx.Ops)
				layout.Flex{Axis: layout.Vertical}.Layout(gtx,
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, fmt.Sprintf("%v", last), cc)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, fmt.Sprintf("token: %v", last.IDToken), cc)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, fmt.Sprintf("provider: %v", last.Provider), cc)
					}),
					layout.Rigid(func(gtx layout.Context) layout.Dimensions {
						return widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, fmt.Sprintf("code: %v", last.Code), cc)
					}),
				)

				evt.Frame(ops)
			}
		}
	}()

	app.Main()
}

func init() {
	gioauth.DefaultProviders = []providers.Provider{
		&google.Provider{
			WebClientID:     "413620002498-30v79n5m1f715fgr89df8c646a3gliu5.apps.googleusercontent.com",
			DesktopClientID: "413620002498-sv2sftp4fk3ttkdr9uampsi03eg5c66h.apps.googleusercontent.com",
			RedirectURL:     "https://gio-demo.inkeliz.com",
		},
		&apple.Provider{
			ServiceIdentifier: "GioDemo",
			RedirectURL:       "https://gio-demo.inkeliz.com",
			SchemeURL:         "giodemo.oauth",
		},
	}

	fmt.Println((&google.Provider{
		WebClientID:     "413620002498-30v79n5m1f715fgr89df8c646a3gliu5.apps.googleusercontent.com",
		DesktopClientID: "413620002498-sv2sftp4fk3ttkdr9uampsi03eg5c66h.apps.googleusercontent.com",
	}).Scheme())
}

func nonce() string {
	id := make([]byte, 16)
	rand.Read(id)
	return hex.EncodeToString(id)
}
