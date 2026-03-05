package pushnotification

/*
extern void setupSwizzling(void);
extern void requestPushToken(void* handler);
*/
import "C"
import (
	"errors"
	"runtime"
	"runtime/cgo"
	"sync"
	"unsafe"
)

type driver struct {
	config Config
	mutex  sync.Mutex
}

func attachDriver(push *Push, config Config) {
	d := &driver{}
	push.driver = d
	configureDriver(d, config)
	if config.RunOnMain != nil {
		go config.RunOnMain(func() {
			C.setupSwizzling()
		})
	}
}

func configureDriver(d *driver, config Config) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.config = config
}

func (d *driver) requestToken() (Token, error) {
	c := make(chan struct {
		Token Token
		Error error
	}, 1)
	fn := func(token Token, err error) {
		c <- struct {
			Token Token
			Error error
		}{
			Token: token,
			Error: err,
		}
	}

	h := cgo.NewHandle(fn)

	d.mutex.Lock()
	runOnMain := d.config.RunOnMain
	d.mutex.Unlock()

	go runOnMain(func() {
		C.requestPushToken(unsafe.Pointer(h))
	})

	r := <-c
	runtime.KeepAlive(h)
	h.Delete()

	return r.Token, r.Error
}

//export gioplugins_pushnotification_on_push_token_received
func gioplugins_pushnotification_on_push_token_received(h unsafe.Pointer, token *C.char, err *C.char) {
	handle := cgo.Handle(h)
	fn, ok := handle.Value().(func(Token, error))
	if !ok {
		return
	}

	if err != nil {
		fn(Token{}, errors.New(C.GoString(err)))
		return
	}

	fn(Token{
		Token:    C.GoString(token),
		Platform: PlatformMacOS,
	}, nil)
}
