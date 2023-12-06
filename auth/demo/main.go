package main

import (
	"fmt"
	"gioui.org/app"
	"gioui.org/font"
	"gioui.org/io/system"
	"gioui.org/io/transfer"
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
	"github.com/gioui-plugins/gio-plugins/plugin"
	"image"
	"image/color"
)

var sharper = text.NewShaper()

func main() {
	w := app.NewWindow(app.Size(500, 500))
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
			select {
			case evt := <-w.Events():
				plugin.Install(w, evt)

				switch evt := evt.(type) {
				case system.DestroyEvent:
					return
				case system.FrameEvent:
					gtx := layout.NewContext(ops, evt)
					{
						s := clip.RRect{Rect: image.Rectangle{Max: gtx.Constraints.Max}}.Push(gtx.Ops)
						paint.ColorOp{Color: color.NRGBA{R: 240, G: 240, B: 240, A: 255}}.Add(gtx.Ops)
						paint.PaintOp{}.Add(gtx.Ops)
						s.Pop()
					}

					transfer.SchemeOp{Tag: w}.Add(gtx.Ops)

					if buttonGoogle.Clicked(gtx) {
						gioauth.RequestOp{Tag: w, Provider: google.IdentifierGoogle, Nonce: "nonce"}.Add(gtx.Ops)
					}

					if buttonApple.Clicked(gtx) {
						gioauth.RequestOp{Tag: w, Provider: apple.IdentifierApple, Nonce: "nonce"}.Add(gtx.Ops)
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

					gioauth.ListenOp{Tag: w}.Add(gtx.Ops)
					for _, evt := range gtx.Events(w) {
						switch evt := evt.(type) {
						case gioauth.AuthEvent:
							last = auth.AuthenticatedEvent(evt)
						case transfer.URLEvent:
							last.Provider = providers.Identifier(evt.URL.String())
						}
					}

					op.Offset(image.Pt(0, gtx.Dp(60))).Push(gtx.Ops)
					widget.Label{}.Layout(gtx, sharper, font.Font{Weight: font.Light}, 14, fmt.Sprintf("%v", last), cc)

					evt.Frame(ops)
				}
			}
		}
	}()

	app.Main()
}

func init() {
	gioauth.DefaultProviders = []providers.Provider{
		&google.Provider{
			WebClientID:     "295108043302-vha4imqq2ojrj8e5pjdbvqqq295c4v49.apps.googleusercontent.com",
			DesktopClientID: "295108043302-pulkf3bn908k63vapq80u084emiukfgg.apps.googleusercontent.com",
			RedirectURL:     "",
		},
		&apple.Provider{
			ServiceIdentifier: "InsteLikes",
			RedirectURL:       "https://instelikes.com/api/apple-oauth",
		},
	}
}
