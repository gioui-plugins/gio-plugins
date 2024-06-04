package internal

import (
	_ "embed"
	"gioui.org/font"
	"gioui.org/font/opentype"
	"gioui.org/text"
)

var ShaperGoogleRoboto = text.NewShaper(text.NoSystemFonts(), text.WithCollection([]text.FontFace{
	{Font: font.Font{}, Face: register(FontGoogleRoboto)},
}))

//go:embed google_font.ttf
var FontGoogleRoboto []byte

func register(ttf []byte) opentype.Face {
	face, _ := opentype.Parse(ttf)
	return face
}
