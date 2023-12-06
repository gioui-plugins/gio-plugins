package authlayout

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout/internal"
	"github.com/inkeliz/giosvg"
	"image/color"
)

// DefaultLightGoogleButtonStyle is the default style for Google buttons.
var DefaultLightGoogleButtonStyle = ButtonStyle{
	Text:            "Continue with Google",
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 60, G: 64, B: 67, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	Format:          FormatRounded,
}

// DefaultDarkGoogleButtonStyle is the default style for Google buttons.
var DefaultDarkGoogleButtonStyle = ButtonStyle{
	Text:                "Continue with Google",
	TextSize:            unit.Dp(16),
	TextFont:            font.Font{},
	TextShaper:          internal.ShaperGoogleRoboto,
	TextColor:           color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	TextAlignment:       layout.Middle,
	IconAlignment:       layout.Start,
	BackgroundColor:     color.NRGBA{R: 66, G: 133, B: 244, A: 255},
	BackgroundIconColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	Format:              FormatRounded,
}

// GoogleDummyButton is a button that can be used to sign in with Google.
// It doesn't perform any action, it just displays a button.
type GoogleDummyButton struct {
	ButtonStyle
	Pointer
	icon *giosvg.Icon
}

// Layout lays out the button, with the default text (from ButtonStyle).
func (g *GoogleDummyButton) Layout(gtx layout.Context) layout.Dimensions {
	return g.LayoutText(gtx, g.Text)
}

// LayoutText lays out the button with the given text.
func (g *GoogleDummyButton) LayoutText(gtx layout.Context, text string) layout.Dimensions {
	if g.icon == nil {
		g.icon = giosvg.NewIcon(internal.VectorGoogleLogo)
	}
	return g.layoutText(gtx, g.icon, &g.Pointer, text, 0, gtx.Dp(24))
}
