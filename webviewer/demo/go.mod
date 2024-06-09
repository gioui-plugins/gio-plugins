module demo

go 1.22

replace gioui.org => ../../../gio

replace github.com/gioui-plugins/gio-plugins => ../../

require (
	gioui.org v0.6.1-0.20240607083507-1151eac07d84
	github.com/gioui-plugins/gio-plugins v0.0.0-00010101000000-000000000000
	golang.org/x/exp/shiny v0.0.0-20230425010034-47ecfdc1ba53
)

require (
	gioui.org/cmd v0.0.0-20240602111522-ddde16a09e12 // indirect
	gioui.org/cpu v0.0.0-20220412190645-f1e9e8c3b1f7 // indirect
	gioui.org/shader v1.0.8 // indirect
	git.wow.st/gmp/jni v0.0.0-20210610011705-34026c7e22d0 // indirect
	github.com/akavel/rsrc v0.10.1 // indirect
	github.com/go-text/typesetting v0.1.1 // indirect
	github.com/inkeliz/go_inkwasm v0.1.23-0.20240519174017-989fbe5b10f6 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	golang.org/x/tools v0.21.0 // indirect
)
