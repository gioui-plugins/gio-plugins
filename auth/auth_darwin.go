package auth

/*
#cgo CFLAGS: -Werror -xobjective-c -fmodules -fobjc-arc

#import <Foundation/Foundation.h>

#if TARGET_OS_IOS
#import <UIKit/UIKit.h>
#else
#import <Appkit/AppKit.h>
#endif

extern CFTypeRef gioplugins_auth_createContextProvider(CFTypeRef viewRef, uintptr_t id);
extern uintptr_t gioplugins_auth_general_open(CFTypeRef viewRef, char * url, char * scheme, uintptr_t id);
extern CFTypeRef gioplugins_auth_apple_open(CFTypeRef viewRef);

*/
import "C"
import (
	"github.com/gioui-plugins/gio-plugins/auth/providers"
	"github.com/gioui-plugins/gio-plugins/auth/providers/apple"
	"runtime/cgo"
	"sync"
)

type driver struct {
	config Config
	mutex  sync.Mutex

	sendURL  func(url string) error
	sendResp func(r Event)
	context  C.CFTypeRef

	cgoHandler cgo.Handle
}

func attachDriver(house *Auth, config Config) {
	house.driver = driver{
		sendURL:  house.ProcessCustomSchemeCallback,
		sendResp: house.sendResponse,
	}
	house.driver.cgoHandler = cgo.NewHandle(&house.driver)
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.mutex.Lock()
	defer driver.mutex.Unlock()

	old := driver.config
	driver.config = config

	if old.View != config.View {
		if old.View != 0 {
			C.CFRelease(C.CFTypeRef(old.View))
		}
		driver.context = C.gioplugins_auth_createContextProvider(C.CFTypeRef(config.View), C.uintptr_t(driver.cgoHandler))
	}
}

func (d *driver) open(provider providers.Provider, nonce string) error {
	switch provider.(type) {
	case *apple.Provider:
		return d.openApple(provider, nonce)
	default:
		return d.openAny(provider, nonce)
	}
}

func (d *driver) openAny(provider providers.Provider, nonce string) error {
	err := C.gioplugins_auth_general_open(
		d.context,
		C.CString(provider.URL(nonce)),
		C.CString(provider.Scheme()),
		C.uintptr_t(d.cgoHandler),
	)
	if err != 0 {
		return ErrProviderNotAllowed
	}
	return nil
}

func (d *driver) openApple(provider providers.Provider, nonce string) error {
	if !isOnAppStore() {
		return d.openAny(provider, nonce)
	}

	C.gioplugins_auth_apple_open(d.context)
	return nil
}

//export auth_general_callback
func auth_general_callback(u *C.char, id C.uintptr_t) {
	url := C.GoString(u)
	driver, ok := cgo.Handle(uintptr(id)).Value().(*driver)
	if !ok {
		return
	}

	driver.sendURL(url)
}

//export auth_apple_callback
func auth_apple_callback(c *C.char, t *C.char, id C.uintptr_t) {
	code := C.GoString(c)
	token := C.GoString(t)

	driver, ok := cgo.Handle(uintptr(id)).Value().(*driver)
	if !ok {
		return
	}

	driver.sendResp(AuthenticatedEvent{
		Provider: apple.IdentifierApple,
		Code:     code,
		IDToken:  token,
	})
}

//export auth_report_error
func auth_report_error(code C.uintptr_t, id C.uintptr_t) {
	driver, ok := cgo.Handle(uintptr(id)).Value().(*driver)
	if !ok {
		return
	}

	const (
		_ASWebAuthenticationSessionErrorCodeCanceledLogin                  = 1
		_ASWebAuthenticationSessionErrorCodePresentationContextNotProvided = 2
		_ASWebAuthenticationSessionErrorCodePresentationContextInvalid     = 3
	)

	switch int(code) {
	case _ASWebAuthenticationSessionErrorCodeCanceledLogin:
		driver.sendResp(ErrorEvent{Error: ErrUserCancelled})
	default:
		driver.sendResp(ErrorEvent{Error: ErrNotConfigured})
	}
}

//export auth_apple_report_error
func auth_apple_report_error(code C.uintptr_t, id C.uintptr_t) {
	driver, ok := cgo.Handle(uintptr(id)).Value().(*driver)
	if !ok {
		return
	}

	const (
		_ASAuthorizationErrorCanceled        = 1001
		_ASAuthorizationErrorFailed          = 1004
		_ASAuthorizationErrorInvalidResponse = 1002
		_ASAuthorizationErrorNotHandled      = 1003
		_ASAuthorizationErrorUnknown         = 1000
		_ASAuthorizationErrorNotInteractive  = 1005
	)

	switch int(code) {
	case _ASAuthorizationErrorCanceled:
		driver.sendResp(ErrorEvent{Error: ErrUserCancelled})
	case _ASAuthorizationErrorInvalidResponse:
		driver.sendResp(ErrorEvent{Error: ErrProviderNotAllowed})
	default:
		driver.sendResp(ErrorEvent{Error: ErrNotConfigured})
	}
}
