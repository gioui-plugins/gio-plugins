package giosafedata

import (
	"reflect"

	"gioui.org/io/event"
	"gioui.org/op"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/safedata"
)

var wantOps = []reflect.Type{
	reflect.TypeOf(&WriteSecretOp{}),
	reflect.TypeOf(&ReadSecretOp{}),
	reflect.TypeOf(&DeleteSecretOp{}),
	reflect.TypeOf(&ListSecretOp{}),
}

type internalOp interface {
	execute(plugin *safedataPlugin)
}

type WriteSecretOp struct {
	Tag    event.Tag
	Secret safedata.Secret
}

var poolWriteSecretOp = plugin.NewOpPool[WriteSecretOp]()

func (o WriteSecretOp) Add(ops *op.Ops) {
	poolWriteSecretOp.WriteOp(ops, o)
}

func (o WriteSecretOp) execute(plugin *safedataPlugin) {
	go func() {
		defer poolWriteSecretOp.Release(&o)

		if err := plugin.client.Set(o.Secret); err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		}
	}()
}

type ReadSecretOp struct {
	Tag        event.Tag
	Identifier string
}

var poolReadSecretOp = plugin.NewOpPool[ReadSecretOp]()

func (o ReadSecretOp) Add(ops *op.Ops) {
	poolReadSecretOp.WriteOp(ops, o)
}

func (o ReadSecretOp) execute(plugin *safedataPlugin) {
	go func() {
		defer poolReadSecretOp.Release(&o)

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

type DeleteSecretOp struct {
	Tag        event.Tag
	Identifier string
}

var poolDeleteSecretOp = plugin.NewOpPool[DeleteSecretOp]()

func (o DeleteSecretOp) Add(ops *op.Ops) {
	poolDeleteSecretOp.WriteOp(ops, o)
}

func (o DeleteSecretOp) execute(plugin *safedataPlugin) {
	go func() {
		defer poolDeleteSecretOp.Release(&o)

		if err := plugin.client.Remove(o.Identifier); err != nil {
			plugin.plugin.SendEvent(o.Tag, ErrorEvent{Error: err})
		}
	}()
}

type ListSecretOp struct {
	Tag    event.Tag
	Buffer []safedata.Secret
}

var poolListSecretOp = plugin.NewOpPool[ListSecretOp]()

func (o ListSecretOp) Add(ops *op.Ops) {
	poolListSecretOp.WriteOp(ops, o)
}

func (o ListSecretOp) execute(plugin *safedataPlugin) {
	go func() {
		defer poolListSecretOp.Release(&o)

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
