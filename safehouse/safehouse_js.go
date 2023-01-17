package safehouse

import (
	"sync"
	"syscall/js"
)

var _jsKeyPrefix = "_gioplugins_safehouse_"

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
	internal safeHouse
	mutex    sync.Mutex
}

func attachDriver(house *SafeHouse) {
	house.driver = driver{}
}

func (c *driver) setSecret(secret Secret) error {
	return c.driver().setSecret(secret)
}

func (c *driver) listSecret(looper Looper) error {
	return c.driver().listSecret(looper)
}

func (c *driver) getSecret(identifier string) (Secret, error) {
	return c.driver().getSecret(identifier)
}

func (c *driver) removeSecret(identifier string) error {
	return c.driver().removeSecret(identifier)
}

func (c *driver) driver() safeHouse {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.internal != nil {
		return c.internal
	}

	switch {
	case _localstorage.IsUndefined():
		c.internal = driverUnsupported{}
	case _passwordCredential.IsUndefined():
		c.internal = driverLocalStorage{}
	case _credentials.IsUndefined():
		c.internal = driverLocalStorage{}
	default:
		c.internal = driverCredentialsManagement{}
	}

	return c.internal
}

type driverCredentialsManagement struct{}

func (d driverCredentialsManagement) setSecret(secret Secret) error {
	obj := js.Global().Get("Object").New()
	obj.Set("name", secret.Identifier)
	obj.Set("password", secret.Data)

	res := make(chan error, 1)
	success := js.FuncOf(func(this interface{}, args []interface{}) any {
		res <- nil
		return nil
	})
	failure := js.FuncOf(func(this interface{}, args []interface{}) any {
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
	success := js.FuncOf(func(this interface{}, args []interface{}) any {
		if len(args) == 0 {
			err <- ErrUserRefused
			return nil
		}

		cred, ok := args[0].(js.Value)
		if !ok {
			err <- ErrUserRefused
			return nil
		}

		pass := cred.Get("password")
		name := cred.Get("name")

		for _, x := range []js.Value{pass, name} {
			if x.IsUndefined() || x.IsNull() {
				err <- ErrUserRefused
				return nil
			}
		}

		secret.Data = []byte(pass.String())
		secret.Identifier = name.String()

		err <- nil
		return nil
	})
	failure := js.FuncOf(func(this interface{}, args []interface{}) any {
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

type driverLocalStorage struct{}

func (d driverLocalStorage) setSecret(secret Secret) error {
	_localstorage.Call("setItem", secret.Identifier, secret.Data)
	return nil
}

func (d driverLocalStorage) listSecret(looper Looper) error {
	keys := _reflect.Call("ownKeys", _localstorage)
	if keys.IsUndefined() {
		return nil
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

func (d driverLocalStorage) getSecret(identifier string, secret *Secret) error {
	content := _localstorage.Call("getItem", d.keyFor(identifier))
	if !content.Truthy() {
		return ErrNotFound
	}

	secret.Identifier = identifier
	secret.Data = []byte(content.String())

	return nil
}

func (d driverLocalStorage) removeSecret(identifier string) error {
	_localstorage.Call("removeItem", d.keyFor(identifier))
	return nil
}

func (d driverLocalStorage) keyFor(id string) string {
	return _jsKeyPrefix + id
}

func (d driverLocalStorage) rawKeyFor(id string) string {
	if len(id) < len(_jsKeyPrefix) {
		return id
	}
	return id[len(_jsKeyPrefix):]
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
