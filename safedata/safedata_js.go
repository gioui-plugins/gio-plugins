package safedata

import (
	"encoding/hex"
	"sync"
	"syscall/js"
)

var (
	_getOptions         js.Value
	_reflect            js.Value
	_localstorage       js.Value
	_credentials        js.Value
	_passwordCredential js.Value
)

func init() {
	_getOptions = js.Global().Get("Object").New()
	_getOptions.Set("password", true)

	_reflect = js.Global().Get("Reflect")
	_localstorage = js.Global().Get("localStorage")
	_credentials = js.Global().Get("navigator").Get("credentials")
	_passwordCredential = js.Global().Get("PasswordCredential")
}

type driver struct {
	internal safeData
	mutex    sync.Mutex
}

func attachDriver(house *SafeData, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	switch d := driver.driver().(type) {
	case driverLocalStorage:
		d.prefix = config.App
	default:
	}
}

func (c *driver) setSecret(secret Secret) error {
	return c.driver().setSecret(secret)
}

func (c *driver) listSecret(looper Looper) error {
	return c.driver().listSecret(looper)
}

func (c *driver) getSecret(identifier string, secret *Secret) error {
	return c.driver().getSecret(identifier, secret)
}

func (c *driver) removeSecret(identifier string) error {
	return c.driver().removeSecret(identifier)
}

func (c *driver) driver() safeData {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.internal == nil {
		switch {
		case _localstorage.IsUndefined():
			c.internal = &driverUnsupported{}
			/*
				case _passwordCredential.IsUndefined():
					c.internal = driverLocalStorage{}
				case _credentials.IsUndefined():
					c.internal = driverLocalStorage{}
			*/
		default:
			c.internal = &driverLocalStorage{}
		}
	}

	return c.internal
}

type driverCredentialsManagement struct{}

func (d driverCredentialsManagement) setSecret(secret Secret) error {
	obj := js.Global().Get("Object").New()
	obj.Set("name", secret.Identifier)
	obj.Set("password", hex.EncodeToString(secret.Data))

	res := make(chan error, 1)
	success := js.FuncOf(func(this js.Value, args []js.Value) any {
		res <- nil
		return nil
	})
	failure := js.FuncOf(func(this js.Value, args []js.Value) any {
		res <- ErrUserRefused
		return nil
	})

	go func() {
		_credentials.Call("store", _passwordCredential.New(obj)).Call("then", success, failure)
	}()

	return <-res
}

func (d driverCredentialsManagement) listSecret(looper Looper) error {
	looper("")
	return nil
}

func (d driverCredentialsManagement) getSecret(_ string, secret *Secret) error {
	err := make(chan error, 1)
	success := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			err <- ErrUserRefused
			return nil
		}

		cred := args[0]

		pass := cred.Get("password")
		name := cred.Get("name")

		for _, x := range []js.Value{pass, name} {
			if x.IsUndefined() || x.IsNull() {
				err <- ErrUserRefused
				return nil
			}
		}

		data, e := hex.DecodeString(pass.String())
		if e != nil {
			err <- e
			return nil
		}

		secret.Identifier = name.String()
		secret.Data = data

		err <- nil
		return nil
	})
	failure := js.FuncOf(func(this js.Value, args []js.Value) any {
		err <- ErrUserRefused
		return nil
	})

	_credentials.Call("get", _getOptions).Call("then", success, failure)
	return <-err
}

func (d driverCredentialsManagement) removeSecret(identifier string) error {
	_credentials.Call("preventSilentAccess")
	return nil
}

type driverLocalStorage struct {
	prefix string
}

func (d driverLocalStorage) setSecret(secret Secret) error {
	_localstorage.Call("setItem", secret.Identifier, hex.EncodeToString(secret.Data))
	return nil
}

func (d driverLocalStorage) listSecret(looper Looper) error {
	keys := _reflect.Call("ownKeys", _localstorage)
	if keys.IsUndefined() {
		return ErrNotFound
	}

	for i := 0; i < keys.Length(); i++ {
		name := keys.Index(i)
		if !name.Truthy() {
			continue
		}
		looper(d.rawKeyFor(name.String()))
	}
	return nil
}

func (d driverLocalStorage) getSecret(identifier string, secret *Secret) (err error) {
	content := _localstorage.Call("getItem", d.keyFor(identifier))
	if !content.Truthy() {
		return ErrNotFound
	}

	secret.Identifier = identifier
	secret.Data, err = hex.DecodeString(content.String())
	return err
}

func (d driverLocalStorage) removeSecret(identifier string) error {
	_localstorage.Call("removeItem", d.keyFor(identifier))
	return nil
}

func (d driverLocalStorage) keyFor(id string) string {
	return d.prefix + id
}

func (d driverLocalStorage) rawKeyFor(id string) string {
	if len(id) < len(d.prefix) {
		return id
	}
	return id[len(d.prefix):]
}

type driverUnsupported struct{}

func (d driverUnsupported) setSecret(_ Secret) error {
	return ErrUnsupported
}

func (d driverUnsupported) listSecret(looper Looper) error {
	return ErrUnsupported
}

func (d driverUnsupported) getSecret(identifier string, secret *Secret) error {
	return ErrUnsupported
}

func (d driverUnsupported) removeSecret(identifier string) error {
	return ErrUnsupported
}
