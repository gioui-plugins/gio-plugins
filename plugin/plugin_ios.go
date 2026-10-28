//go:build darwin && ios

package plugin

import "gioui.org/op/paint"

func NewEmbedOp(view ...uintptr) paint.EmbedOp {
	return paint.UIKitView{ViewController: view[0]}.Op()
}
