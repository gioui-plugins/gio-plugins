package webview

import (
	"github.com/inkeliz/go_inkwasm/inkwasm"
)

//inkwasm:func console.log
func log(s inkwasm.Object)

//inkwasm:func .querySelector
func querySelector(o inkwasm.Object, s string) inkwasm.Object

//inkwasm:func document.createElement
func createElement(s string) inkwasm.Object

//inkwasm:func .setAttribute
func setAttribute(o inkwasm.Object, s string, v string)

//inkwasm:set .innerText
func setInnerText(o inkwasm.Object, s string)

//inkwasm:func .prepend
func prepend(o inkwasm.Object, c inkwasm.Object)

//inkwasm:func .removeChild
func removeChild(o inkwasm.Object, c inkwasm.Object)

//inkwasm:set .style.width
func setStyleWidth(o inkwasm.Object, s string)

//inkwasm:set .style.height
func setStyleHeight(o inkwasm.Object, s string)

//inkwasm:set .style.position
func setStylePosition(o inkwasm.Object, s string)

//inkwasm:set .style.zIndex
func setStyleZIndex(o inkwasm.Object, s string)

//inkwasm:set .style.border
func setStyleBorder(o inkwasm.Object, s string)

//inkwasm:set .style.display
func setStyleDisplay(o inkwasm.Object, s string)

//inkwasm:set .style.top
func setStyleTop(o inkwasm.Object, s string)

//inkwasm:set .style.left
func setStyleLeft(o inkwasm.Object, s string)

//inkwasm:set .src
func setSrc(o inkwasm.Object, s string)
