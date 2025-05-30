package authlayout

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout/internal"
	"image/color"
)

// DefaultLightGoogleButtonStyle is the default style for Google buttons.
var DefaultLightGoogleButtonStyle = ButtonStyle{
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 60, G: 64, B: 67, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IconVector:      internal.VectorGoogleLogo,
	IconPadding:     24,
	Format:          FormatRounded,
}

// DefaultDarkGoogleButtonStyle is the default style for Google buttons.
var DefaultDarkGoogleButtonStyle = ButtonStyle{
	TextSize:            unit.Dp(16),
	TextFont:            font.Font{},
	TextShaper:          internal.ShaperGoogleRoboto,
	TextColor:           color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	TextAlignment:       layout.Middle,
	IconAlignment:       layout.Start,
	BackgroundColor:     color.NRGBA{R: 66, G: 133, B: 244, A: 255},
	BackgroundIconColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IconVector:          internal.VectorGoogleLogo,
	IconPadding:         24,
	Format:              FormatRounded,
}

// DefaultGoogleTextContinue is the default text for Google buttons.
var DefaultGoogleTextContinue = [2]string{
	"Continue with Google",
	"Sign in",
}

func NewGoogleButton() *Button {
	return &Button{
		ButtonStyle: DefaultLightGoogleButtonStyle,
		ButtonTexts: DefaultGoogleTextContinue,
	}
}

func NewGoogleButtonDark() *Button {
	return &Button{
		ButtonStyle: DefaultDarkGoogleButtonStyle,
		ButtonTexts: DefaultGoogleTextContinue,
	}
}
