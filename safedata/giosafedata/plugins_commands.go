package giosafedata

import (
	"gioui.org/io/input"
	"reflect"

	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/safedata"
)

var (
	wantCommands = []reflect.Type{
		reflect.TypeOf(WriteSecretCmd{}),
		reflect.TypeOf(ReadSecretCmd{}),
		reflect.TypeOf(DeleteSecretCmd{}),
		reflect.TypeOf(ListSecretCmd{}),
	}
)

// ReadSecretCmd requests a secret by its identifier.
// It will issue a SecretsEvent or an ErrorEvent.
type ReadSecretCmd struct {
	Tag        event.Tag
	Identifier string
}

// WriteSecretCmd writes a secret.
// It will issue an ErrorEvent if an error occurs.
type WriteSecretCmd struct {
	Tag    event.Tag
	Secret safedata.Secret
}

// DeleteSecretCmd deletes a secret by its identifier.
// It will issue an ErrorEvent if an error occurs.
type DeleteSecretCmd struct {
	Tag        event.Tag
	Identifier string
}

// ListSecretCmd requests a list of secrets,
// the buffer must be a slice of safedata.Secret.
//
// If the buffer is too small, it will be resized.
// It will issue a SecretsEvent or an ErrorEvent.
type ListSecretCmd struct {
	Tag    event.Tag
	Buffer []safedata.Secret
}

func (o ReadSecretCmd) ImplementsCommand()   {}
func (o WriteSecretCmd) ImplementsCommand()  {}
func (o DeleteSecretCmd) ImplementsCommand() {}
func (o ListSecretCmd) ImplementsCommand()   {}

type internalCmd interface {
	input.Command
	execute(plugin *safedataPlugin)
}

func (o WriteSecretCmd) execute(plugin *safedataPlugin) {
	go func() {
		if err := plugin.client.Set(o.Secret); err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		}
	}()
}

func (o ReadSecretCmd) execute(plugin *safedataPlugin) {
	go func() {
		secret, err := plugin.client.Get(o.Identifier)
		if err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		} else {
			plugin.plugin.SendEvent(o.Tag, SecretsEvent{
				Secrets: []safedata.Secret{secret},
			})
		}
	}()
}

func (o DeleteSecretCmd) execute(plugin *safedataPlugin) {
	go func() {
		if err := plugin.client.Remove(o.Identifier); err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		}
	}()
}

func (o ListSecretCmd) execute(plugin *safedataPlugin) {
	go func() {
		buf := o.Buffer
		i := 0
		var nerr error
		err := plugin.client.List(func(identifier string) (next bool) {
			if i >= cap(buf) {
				buf = append(buf, safedata.Secret{})
			}

			if err := plugin.client.View(identifier, &buf[i]); err != nil {
				nerr = err
				return false
			}
			return true
		})

		if nerr != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: nerr})
			return
		}

		if err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		} else {
			plugin.plugin.SendEvent(o.Tag, SecretsEvent{Secrets: buf})
		}
	}()
}
