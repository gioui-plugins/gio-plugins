package deeplink

import "strings"

// Schemes is a list of schemes that the app will listen to, that
// must be separated by a comma, for example: "myapp,anotherapp".
var schemes string
var schemeList []string

func init() {
	if schemes == "" {
		panic(`deeplink: no schemes defined, you must use -ldflags "-X github.com/gioui-plugins/gio-plugins/deeplink.schemes=yourscheme,anotherscheme"`)
	}
	schemeList = strings.Split(schemes, ",")
}
