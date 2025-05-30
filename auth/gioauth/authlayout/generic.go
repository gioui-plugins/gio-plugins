package authlayout

import (
	"gioui.org/layout"
	"github.com/gioui-plugins/gio-plugins/auth/gioauth/authlayout/internal"
)

var (
	// ShaperGoogleRoboto is the default shaper for Google and Apple buttons,
	// it's exposed for external use, so you can use it in your own
	// widgets.
	ShaperGoogleRoboto = internal.ShaperGoogleRoboto
)

// Button is widget that display a button with a text and an icon,
// and can be clicked, you can get the click event using the Pointer.Clicked method.
type Button struct {
	ButtonStyle
	ButtonTexts
	Pointer
}

// Layout lays out the button, with the default text (from ButtonStyle).
func (g *Button) Layout(gtx layout.Context) layout.Dimensions {
	return g.LayoutText(gtx, g.ButtonTexts)
}

// LayoutText lays out the button with the given text.
func (g *Button) LayoutText(gtx layout.Context, texts ButtonTexts) layout.Dimensions {
	return g.ButtonStyle.LayoutText(gtx, &g.Pointer, texts)
}
