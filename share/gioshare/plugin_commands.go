package gioshare

import (
	"gioui.org/io/event"
	"reflect"
)

var wantCommands = []reflect.Type{
	reflect.TypeOf(WebsiteCmd{}),
	reflect.TypeOf(TextCmd{}),
}

// TextCmd represents the text to be shared.
type TextCmd struct {
	Tag   event.Tag
	Title string
	Text  string
}

// WebsiteCmd represents the website/link to be shared.
type WebsiteCmd struct {
	Tag   event.Tag
	Title string
	Text  string
	Link  string
}

func (o TextCmd) ImplementsCommand()    {}
func (o WebsiteCmd) ImplementsCommand() {}
