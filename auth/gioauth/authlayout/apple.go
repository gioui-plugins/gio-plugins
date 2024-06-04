package authlayout

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout/internal"
	"github.com/inkeliz/giosvg"
	"image/color"
)

// DefaultLightAppleButtonStyle is the default style for Apple buttons.
var DefaultLightAppleButtonStyle = ButtonStyle{
	Text:            "Continue with Apple",
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	IconColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	Format:          FormatRounded,
}

// DefaultDarkAppleButtonStyle is the default style for Apple buttons.
var DefaultDarkAppleButtonStyle = ButtonStyle{
	Text:            "Continue with Apple",
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IconColor:       color.NRGBA{A: 255},
	Format:          FormatRounded,
}

// AppleDummyButton is a button that can be used to sign in with Apple.
// It doesn't perform any action, it just displays a button.
type AppleDummyButton struct {
	ButtonStyle
	Pointer
	icon *giosvg.Icon
}

// Layout lays out the button, with the default text (from ButtonStyle).
func (g *AppleDummyButton) Layout(gtx layout.Context) layout.Dimensions {
	return g.LayoutText(gtx, g.Text)
}

// LayoutText lays out the button with the given text.
func (g *AppleDummyButton) LayoutText(gtx layout.Context, text string) layout.Dimensions {
	if g.icon == nil {
		g.icon = giosvg.NewIcon(internal.VectorAppleLogo)
	}
	return g.layoutText(gtx, g.icon, &g.Pointer, text, 0, gtx.Dp(24))
}
