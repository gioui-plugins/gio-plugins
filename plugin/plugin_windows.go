package plugin

import "gioui.org/op/paint"

func NewEmbedOp(view ...uintptr) paint.EmbedOp {
	return paint.Win32View{HWND: view[0]}.Op()
}
