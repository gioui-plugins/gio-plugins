module demo

go 1.19

replace gioui.org/cmd => ../../../gio-cmd

replace github.com/gioui-plugins/gio-plugins => ../../

require (
	gioui.org v0.0.0-20230429160049-0e5ec18a82e9
	github.com/gioui-plugins/gio-plugins v0.0.0-00010101000000-000000000000
)

require (
	gioui.org/cmd v0.0.0-20230701070152-940364d3e94a // indirect
	gioui.org/cpu v0.0.0-20220412190645-f1e9e8c3b1f7 // indirect
	gioui.org/shader v1.0.6 // indirect
	git.wow.st/gmp/jni v0.0.0-20210610011705-34026c7e22d0 // indirect
	github.com/akavel/rsrc v0.10.1 // indirect
	github.com/go-text/typesetting v0.0.0-20230413204129-b4f0492bf7ae // indirect
	golang.org/x/crypto v0.11.0 // indirect
	golang.org/x/exp v0.0.0-20221012211006-4de253d81b95 // indirect
	golang.org/x/exp/shiny v0.0.0-20220921164117-439092de6870 // indirect
	golang.org/x/image v0.5.0 // indirect
	golang.org/x/mod v0.8.0 // indirect
	golang.org/x/sync v0.1.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/text v0.11.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
)
