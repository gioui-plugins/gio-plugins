package internal

import (
	"github.com/go-ole/go-ole"
	"unsafe"
)

/*
IUriRuntimeClassFactory: https://docs.microsoft.com/en-us/cpp/cppcx/wrl/how-to-activate-and-use-a-windows-runtime-component-using-wrl?view=msvc-170
*/
var (
	IUriRuntimeClassFactoryCLSID = "Windows.Foundation.Uri"
	IUriRuntimeClassFactoryGUID  = ole.NewGUID("44A9796F-723E-4FDF-A218-033E75B0C084")
)

type (
	IUriRuntimeClassFactory struct {
		VTBL *IUriRuntimeClassFactoryVTBL
	}
	IUriRuntimeClassFactoryVTBL struct {
		ole.IInspectableVtbl
		CreateUri             uintptr // HRESULT CreateUri([in] HSTRING uri, [out] [retval] Windows.Foundation.Uri** instance);
		CreateWithRelativeUri uintptr // HRESULT CreateWithRelativeUri([in] HSTRING baseUri, [in] HSTRING relativeUri, [out] [retval] Windows.Foundation.Uri** instance);
	}
)

func NewIUriRuntimeClassFactory(r **IUriRuntimeClassFactory) error {
	ins, err := ole.RoGetActivationFactory(IUriRuntimeClassFactoryCLSID, IUriRuntimeClassFactoryGUID)
	if err != nil {
		return err
	}
	*r = (*IUriRuntimeClassFactory)(unsafe.Pointer(ins))
	return nil
}

func (i *IUriRuntimeClassFactory) CreateUri(s string, r **IUriRuntimeClass) (err error) {
	sw, err := ole.NewHString(s)
	if err != nil {
		return err
	}
	//defer ole.DeleteHString(sw)

	return call(i.VTBL.CreateUri, uintptr(unsafe.Pointer(i)), uintptr(sw), uintptr(unsafe.Pointer(r)))
}

/*
	IUriRuntimeClass: https://docs.microsoft.com/en-us/cpp/cppcx/wrl/how-to-activate-and-use-a-windows-runtime-component-using-wrl?view=msvc-170
*/

var (
	IUriRuntimeClassGUID = ole.NewGUID("9E365E57-48B2-4160-956F-C7385120BBFC")
)

type (
	IUriRuntimeClass struct {
		VTBL *interface{}
	}
)
