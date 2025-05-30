package authlayout

import (
	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout/internal"
	"image/color"
)

// DefaultLightAppleButtonStyle is the default style for Apple buttons.
var DefaultLightAppleButtonStyle = ButtonStyle{
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	IconColor:       color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IconVector:      internal.VectorAppleLogo,
	IconPadding:     24,
	Format:          FormatRounded,
}

// DefaultDarkAppleButtonStyle is the default style for Apple buttons.
var DefaultDarkAppleButtonStyle = ButtonStyle{
	TextSize:        unit.Dp(16),
	TextFont:        font.Font{},
	TextShaper:      internal.ShaperGoogleRoboto,
	TextColor:       color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	TextAlignment:   layout.Middle,
	IconAlignment:   layout.Start,
	BackgroundColor: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	IconColor:       color.NRGBA{A: 255},
	IconVector:      internal.VectorAppleLogo,
	IconPadding:     24,
	Format:          FormatRounded,
}

// DefaultAppleTextContinue is the default text for Google buttons.
var DefaultAppleTextContinue = [2]string{
	"Continue with Apple",
	"Sign in",
}

func NewAppleButton() *Button {
	return &Button{
		ButtonStyle: DefaultLightAppleButtonStyle,
		ButtonTexts: DefaultAppleTextContinue,
	}
}

func NewAppleButtonDark() *Button {
	return &Button{
		ButtonStyle: DefaultDarkAppleButtonStyle,
		ButtonTexts: DefaultAppleTextContinue,
	}
}
