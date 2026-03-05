module demo

go 1.24.1

toolchain go1.24.3

replace github.com/gioui-plugins/gio-plugins => ../../

replace gioui.org/cmd => ../../../inste/gio-cmd

require (
	gioui.org v0.9.1-0.20251215212054-7bcb315ee174
	github.com/gioui-plugins/gio-plugins v0.1.0
)

require (
	gioui.org/cmd v0.9.0 // indirect
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20260127013417-d142949d346a // indirect
	github.com/akavel/rsrc v0.10.1 // indirect
	github.com/go-ole/go-ole v1.3.0 // indirect
	github.com/go-text/typesetting v0.3.0 // indirect
	golang.org/x/crypto v0.48.0 // indirect
	golang.org/x/exp/shiny v0.0.0-20250408133849-7e4ce0ab07d0 // indirect
	golang.org/x/image v0.26.0 // indirect
	golang.org/x/mod v0.32.0 // indirect
	golang.org/x/net v0.49.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
	golang.org/x/text v0.34.0 // indirect
	golang.org/x/tools v0.41.0 // indirect
)
