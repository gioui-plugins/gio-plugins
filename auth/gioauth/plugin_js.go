package gioauth

import "syscall/js"

func startURL() string {
	return js.Global().Get("location").Get("href").String()
}
