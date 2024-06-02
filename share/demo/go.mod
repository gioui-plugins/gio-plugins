module demo

go 1.22

replace gioui.org => ../../../gio

replace github.com/gioui-plugins/gio-plugins => ../../

require (
	gioui.org v0.6.0
	github.com/gioui-plugins/gio-plugins v0.0.0-00010101000000-000000000000
)

require (
	gioui.org/cpu v0.0.0-20220412190645-f1e9e8c3b1f7 // indirect
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20210610011705-34026c7e22d0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-text/typesetting v0.1.1 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/exp/shiny v0.0.0-20220921164117-439092de6870 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)
