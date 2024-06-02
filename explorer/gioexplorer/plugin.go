package gioexplorer

import (
	"errors"
	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"reflect"
)

func init() {
	plugin.Register(func(w *app.Window, p *plugin.Plugin) plugin.Handler {
		return &explorerPlugin{window: w, plugin: p}
	})
}

var (
	// ErrUserDecline is returned when the user doesn't select the file.
	ErrUserDecline = errors.New("user exited the file selector without selecting a file")
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")
)

type explorerPlugin struct {
	window *app.Window
	plugin *plugin.Plugin

	config explorer.Config
	client *explorer.Explorer
}

// TypeOp implements the plugin.Handler interface.
func (e *explorerPlugin) TypeOp() []reflect.Type { return nil }

// TypeCommand implements the plugin.Handler interface.
func (e *explorerPlugin) TypeCommand() []reflect.Type { return wantCommands }

// TypeEvent implements the plugin.Handler interface.
func (e *explorerPlugin) TypeEvent() []reflect.Type { return wantEvents }

// Op implements the plugin.Handler interface.
func (e *explorerPlugin) Op(op interface{}) {}

// Execute implements the plugin.Handler interface.
func (e *explorerPlugin) Execute(cmd interface{}) {
	switch cmd := cmd.(type) {
	case OpenFileCmd:
		cmd.execute(e)
	case SaveFileCmd:
		cmd.execute(e)
	}
}

// Event implements the plugin.Handler interface.
func (e *explorerPlugin) Event(evt event.Event) {
	if e == nil {
		return
	}
	switch evt := evt.(type) {
	case app.ViewEvent:
		UpdateConfigFromViewEvent(&e.config, e.window, evt)
		if e.client == nil {
			e.client = explorer.NewExplorer(e.config)
		} else {
			e.client.Configure(e.config)
		}
	}
}
