// SPDX-License-Identifier: Unlicense OR MIT

package share

import (
	"github.com/gioui-plugins/gio-plugins/share/internal"
	"github.com/go-ole/go-ole"
	"sync"
)

type driver struct {
	mutex  sync.Mutex
	config Config

	mode      uint8
	shareable [3]string

	// The definition of those field lives at share_windows_idl.go:
	// It's  important to keep those values here to prevent the content to be freed
	// by GC, so it must live here "forever".
	_IDataTransferManagerInterop *internal.IDataTransferManagerInterop
	_IDataTransferManager        *internal.IDataTransferManager
	_ITypedEventHandler          *internal.ITypedEventHandler

	_IUriRuntimeClassFactory *internal.IUriRuntimeClassFactory
}

func attachDriver(house *Share, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
	house.driver.init()
}

func configureDriver(driver *driver, config Config) {
	driver.config = config
}

func (e *driver) init() {
	if err := ole.RoInitialize(1); err != nil {
		return
	}

	if err := internal.NewIDataTransferManagerInterop(&e._IDataTransferManagerInterop); err != nil {
		return
	}

	if err := internal.NewIUriRuntimeClassFactory(&e._IUriRuntimeClassFactory); err != nil {
		return
	}

	if e._IDataTransferManagerInterop == nil || e._IUriRuntimeClassFactory == nil {
		return
	}

	if err := e._IDataTransferManagerInterop.GetForWindow(e.config.HWND, &e._IDataTransferManager); err != nil {
		return
	}

	callback := func(transfer *internal.IDataTransferManager, args *internal.IDataRequestedEventArgs) int {
		e.mutex.Lock()
		defer e.mutex.Unlock()

		var dataRequest *internal.IDataRequest
		if err := args.GetRequest(&dataRequest); err != nil {
			return ole.E_FAIL
		}

		var dataPackage *internal.IDataPackage
		if err := dataRequest.GetData(&dataPackage); err != nil {
			return ole.E_FAIL
		}

		var dataPackage2 *internal.IDataPackage2
		if err := dataPackage.GetIDataPackage2(&dataPackage2); err != nil {
			return ole.E_FAIL
		}

		var dataProperty *internal.IDataPackagePropertySet
		if err := dataPackage.GetProperties(&dataProperty); err != nil {
			return ole.E_FAIL
		}

		switch e.mode {
		case 0:
			dataProperty.SetTitle(e.shareable[0])
			dataPackage.SetText(e.shareable[1])
		case 1:
			dataProperty.SetTitle(e.shareable[0])
			dataPackage.SetText(e.shareable[1])

			var uri *internal.IUriRuntimeClass
			if err := e._IUriRuntimeClassFactory.CreateUri(e.shareable[2], &uri); err != nil {
				return ole.S_OK
			}

			dataPackage2.SetWebLink(uri)
		}

		dataRequest.SetData(dataPackage)
		return 1
	}

	e._ITypedEventHandler = internal.NewTypedEventHandler(func(transfer *internal.IDataTransferManager, args *internal.IDataRequestedEventArgs) int {
		var r int
		e.config.RunOnMain(func() {
			r = callback(transfer, args)
		})
		return r
	})

	if err := e._IDataTransferManager.AddDataRequested(e._ITypedEventHandler); err != nil {
		return
	}
}

func (e *driver) shareShareable() error {
	go e.config.RunOnMain(func() {
		e._IDataTransferManagerInterop.ShowShareUIWindow(e.config.HWND)
	})
	return nil
}

func (e *driver) shareText(title, text string) error {
	// Mutex prevents changes of shareable data when Window is triggering the callback.
	e.mutex.Lock()
	e.shareable = [3]string{title, text, ""}
	e.mode = 0
	e.mutex.Unlock()

	return e.shareShareable()
}

func (e *driver) shareWebsite(title, description, url string) error {
	// Mutex prevents changes of shareable data when Window is triggering the callback.
	e.mutex.Lock()
	e.shareable = [3]string{title, description, url}
	e.mode = 1
	e.mutex.Unlock()

	return e.shareShareable()
}
