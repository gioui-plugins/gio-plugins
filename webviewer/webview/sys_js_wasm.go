// Code generated by INKWASM BUILD; DO NOT EDIT

package webview

import (
	"runtime"

	"github.com/inkeliz/go_inkwasm/inkwasm"
)

func _log(s inkwasm.Object) {
	__log(s)

}
func __log(s inkwasm.Object)

func _querySelector(o inkwasm.Object, s string) (_ inkwasm.Object) {
	r0 := __querySelector(o, s)
	runtime.KeepAlive(s)

	return r0
}
func __querySelector(o inkwasm.Object, s string) (_ inkwasm.Object)

func _createElement(s string) (_ inkwasm.Object) {
	r0 := __createElement(s)
	runtime.KeepAlive(s)

	return r0
}
func __createElement(s string) (_ inkwasm.Object)

func _setAttribute(o inkwasm.Object, s string, v string) {
	__setAttribute(o, s, v)
	runtime.KeepAlive(s)
	runtime.KeepAlive(v)

}
func __setAttribute(o inkwasm.Object, s string, v string)

func _setInnerText(o inkwasm.Object, s string) {
	__setInnerText(o, s)
	runtime.KeepAlive(s)

}
func __setInnerText(o inkwasm.Object, s string)

func _prepend(o inkwasm.Object, c inkwasm.Object) {
	__prepend(o, c)

}
func __prepend(o inkwasm.Object, c inkwasm.Object)

func _removeChild(o inkwasm.Object, c inkwasm.Object) {
	__removeChild(o, c)

}
func __removeChild(o inkwasm.Object, c inkwasm.Object)

func _setStyleWidth(o inkwasm.Object, s string) {
	__setStyleWidth(o, s)
	runtime.KeepAlive(s)

}
func __setStyleWidth(o inkwasm.Object, s string)

func _setStyleHeight(o inkwasm.Object, s string) {
	__setStyleHeight(o, s)
	runtime.KeepAlive(s)

}
func __setStyleHeight(o inkwasm.Object, s string)

func _setStylePosition(o inkwasm.Object, s string) {
	__setStylePosition(o, s)
	runtime.KeepAlive(s)

}
func __setStylePosition(o inkwasm.Object, s string)

func _setStyleZIndex(o inkwasm.Object, s string) {
	__setStyleZIndex(o, s)
	runtime.KeepAlive(s)

}
func __setStyleZIndex(o inkwasm.Object, s string)

func _setStyleBorder(o inkwasm.Object, s string) {
	__setStyleBorder(o, s)
	runtime.KeepAlive(s)

}
func __setStyleBorder(o inkwasm.Object, s string)

func _setStyleDisplay(o inkwasm.Object, s string) {
	__setStyleDisplay(o, s)
	runtime.KeepAlive(s)

}
func __setStyleDisplay(o inkwasm.Object, s string)

func _setStyleTop(o inkwasm.Object, s string) {
	__setStyleTop(o, s)
	runtime.KeepAlive(s)

}
func __setStyleTop(o inkwasm.Object, s string)

func _setStyleLeft(o inkwasm.Object, s string) {
	__setStyleLeft(o, s)
	runtime.KeepAlive(s)

}
func __setStyleLeft(o inkwasm.Object, s string)

func _setSrc(o inkwasm.Object, s string) {
	__setSrc(o, s)
	runtime.KeepAlive(s)

}
func __setSrc(o inkwasm.Object, s string)
