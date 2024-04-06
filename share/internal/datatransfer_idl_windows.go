//go:build windows
// +build windows

package internal

import (
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

/*
DataTransferManagerInterop: https://docs.microsoft.com/en-us/windows/win32/api/shobjidlcore/nn-shobjidlcore-idatatransfermanagerinterop
*/
var (
	IDataTransferManagerInteropCLSID = "Windows.ApplicationModel.DataTransfer.DataTransferManager"
	IDataTransferManagerInteropGUID  = ole.NewGUID("3A3DCD6C-3EAB-43DC-BCDE-45671CE800C8")
)

type (
	IDataTransferManagerInterop struct {
		VTBL *IDataTransferManagerInteropVTBL
	}
	IDataTransferManagerInteropVTBL struct {
		ole.IUnknownVtbl
		GetForWindow         uintptr // HRESULT GetForWindow([in] HWND appWindow, [in] REFIID riid, [out, optional] dataTransferManager**);
		ShowShareUIForWindow uintptr // HRESULT ShowShareUIForWindow([in] HWND appWindow);
	}
)

func NewIDataTransferManagerInterop(r **IDataTransferManagerInterop) error {
	ins, err := ole.RoGetActivationFactory(IDataTransferManagerInteropCLSID, IDataTransferManagerInteropGUID)
	if err != nil {
		return err
	}
	*r = (*IDataTransferManagerInterop)(unsafe.Pointer(ins))
	return nil
}

func (i *IDataTransferManagerInterop) GetForWindow(hwnd uintptr, r **IDataTransferManager) (err error) {
	return call(i.VTBL.GetForWindow, uintptr(unsafe.Pointer(i)), hwnd, uintptr(unsafe.Pointer(IDataTransferManagerGUID)), uintptr(unsafe.Pointer(r)))
}

func (i *IDataTransferManagerInterop) ShowShareUIWindow(hwnd uintptr) error {
	return call(i.VTBL.ShowShareUIForWindow, uintptr(unsafe.Pointer(i)), hwnd)
}

/*
	IDataTransferManager: http://definitelytyped.org/docs/winrt--winrt/interfaces/windows.applicationmodel.datatransfer.idatatransfermanager.html
*/

var IDataTransferManagerGUID = ole.NewGUID("A5CAEE9B-8708-49D1-8D36-67D25A8DA00C") // GUID from WinSDK 10, see https://www.magnumdb.com/search?q=IDataTransferManager.

type (
	IDataTransferManager struct {
		VTBL *IDataTransferManagerVTBL
	}
	IDataTransferManagerVTBL struct {
		ole.IInspectableVtbl
		AddDataRequested              uintptr // [eventadd] HRESULT DataRequested([in] Windows.Foundation.TypedEventHandler<Windows.ApplicationModel.DataTransfer.DataTransferManager*, Windows.ApplicationModel.DataTransfer.DataRequestedEventArgs*>* eventHandler, [out] [retval] EventRegistrationToken* eventCookie);
		RemoveDataRequested           uintptr // [eventremove] HRESULT DataRequested([in] EventRegistrationToken eventCookie);
		AddTargetApplicationChosen    uintptr // [eventadd] HRESULT TargetApplicationChosen([in] Windows.Foundation.TypedEventHandler<Windows.ApplicationModel.DataTransfer.DataTransferManager*, Windows.ApplicationModel.DataTransfer.TargetApplicationChosenEventArgs*>* eventHandler, [out] [retval] EventRegistrationToken* eventCookie);
		RemoveTargetApplicationChosen uintptr // [eventremove] HRESULT TargetApplicationChosen([in] EventRegistrationToken eventCookie);
	}
)

func (i *IDataTransferManager) AddDataRequested(r *ITypedEventHandler) (err error) {
	var token uintptr
	return call(i.VTBL.AddDataRequested, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(r)), uintptr(unsafe.Pointer(&token)))
}

/*
	IDataRequestedEventArgs: http://definitelytyped.org/docs/microsoft-live-connect--microsoft-live-connect/interfaces/windows.applicationmodel.datatransfer.idatarequestedeventargs.html
*/

var IDataRequestedEventArgsGUID = ole.NewGUID("CB8BA807-6AC5-43C9-8AC5-9BA232163182")

type (
	IDataRequestedEventArgs struct {
		VTBL *IDataRequestedEventArgsVTBL
	}
	IDataRequestedEventArgsVTBL struct {
		ole.IInspectableVtbl
		GetRequest uintptr // [propget] HRESULT Request([out] [retval] Windows.ApplicationModel.DataTransfer.DataRequest** value);
	}
)

func (i *IDataRequestedEventArgs) GetRequest(r **IDataRequest) error {
	return call(i.VTBL.GetRequest, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(r)))
}

/*
	IDataRequest: http://definitelytyped.org/docs/winrt--winrt/interfaces/windows.applicationmodel.datatransfer.idatarequest.html
*/

var IDataRequestGUID = ole.NewGUID("4341AE3B-FC12-4E53-8C02-AC714C415A27")

type (
	IDataRequest struct {
		VTBL *IDataRequestVTBL
	}
	IDataRequestVTBL struct {
		ole.IInspectableVtbl
		GetData             uintptr // [propget] HRESULT Data([out] [retval] Windows.ApplicationModel.DataTransfer.DataPackage** value);
		SetData             uintptr // [propput] HRESULT Data([in] Windows.ApplicationModel.DataTransfer.DataPackage* value);
		Deadline            uintptr // [propget] HRESULT Deadline([out] [retval] Windows.Foundation.DateTime* value);
		FailWithDisplayText uintptr // HRESULT FailWithDisplayText([in] HSTRING value);
		GetDeferral         uintptr // HRESULT GetDeferral([out] [retval] Windows.ApplicationModel.DataTransfer.DataRequestDeferral** value);
	}
)

func (i *IDataRequest) GetData(r **IDataPackage) error {
	return call(i.VTBL.GetData, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(r)))
}

func (i *IDataRequest) SetData(r *IDataPackage) error {
	return call(i.VTBL.SetData, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(r)))
}

/*
	IDataRequestedEventArgs: http://definitelytyped.org/docs/microsoft-live-connect--microsoft-live-connect/interfaces/windows.applicationmodel.datatransfer.idatarequestedeventargs.html
*/

var (
	IDataPackageGUID  = ole.NewGUID("61EBF5C7-EFEA-4346-9554-981D7E198FFE")
	IDataPackage2GUID = ole.NewGUID("041C1FE9-2409-45E1-A538-4C53EEEE04A7")
)

type (
	IDataPackage struct {
		VTBL *IDataPackageVTBL
	}
	IDataPackageVTBL struct {
		ole.IInspectableVtbl
		GetView                  uintptr // HRESULT GetView [out] [retval] Windows.ApplicationModel.DataTransfer.DataPackageView** value);
		Properties               uintptr // [propget] HRESULT Properties [out] [retval] Windows.ApplicationModel.DataTransfer.DataPackagePropertySet** value);
		GetRequestedOperation    uintptr // [propget] HRESULT RequestedOperation [out] [retval] Windows.ApplicationModel.DataTransfer.DataPackageOperation* value);
		SetRequestedOperation    uintptr // [propput] HRESULT RequestedOperation [in] Windows.ApplicationModel.DataTransfer.DataPackageOperation value);
		AddOperationCompleted    uintptr // [eventadd] HRESULT  OperationCompleted [in] Windows.Foundation.TypedEventHandler<Windows.ApplicationModel.DataTransfer.DataPackage*, Windows.ApplicationModel.DataTransfer.OperationCompletedEventArgs*>* handler, [out] [retval] EventRegistrationToken* eventCookie);
		RemoveOperationCompleted uintptr // [eventremove] HRESULT  OperationCompleted [in] EventRegistrationToken eventCookie);
		AddDestroyed             uintptr // [eventadd] HRESULT  Destroyed [in] Windows.Foundation.TypedEventHandler<Windows.ApplicationModel.DataTransfer.DataPackage*, IInspectable*>* handler, [out] [retval] EventRegistrationToken* eventCookie);
		RemoveDestroyed          uintptr // [eventremove] HRESULT  Destroyed [in] EventRegistrationToken eventCookie);
		SetData                  uintptr // HRESULT SetData [in] HSTRING formatId, [in] IInspectable* value);
		SetDataProvider          uintptr // HRESULT SetDataProvider [in] HSTRING formatId, [in] Windows.ApplicationModel.DataTransfer.DataProviderHandler* delayRenderer);
		SetText                  uintptr // HRESULT SetText [in] HSTRING value);
		SetUri                   uintptr // [deprecated("SetUri may be altered or unavailable for releases after Windows Phone 'OSVersion' (TBD). Instead, use SetWebLink or SetApplicationLink.", deprecate, Windows.Foundation.UniversalApiContract, 1.0)] HRESULT  SetUri [in] Windows.Foundation.Uri* value);
		SetHtmlFormat            uintptr // HRESULT SetHtmlFormat [in] HSTRING value);
		ResourceMap              uintptr // [propget] HRESULT  ResourceMap [out] [retval] Windows.Foundation.Collections.IMap<HSTRING, Windows.Storage.Streams.RandomAccessStreamReference*>** value);
		SetRtf                   uintptr // HRESULT SetRtf [in] HSTRING value);
		SetBitmap                uintptr // HRESULT SetBitmap [in] Windows.Storage.Streams.RandomAccessStreamReference* value);
		SetStorageItemsReadOnly  uintptr // [overload("SetStorageItems")] HRESULT  SetStorageItemsReadOnly [in] Windows.Foundation.Collections.IIterable<Windows.Storage.IStorageItem*>* value);
		SetStorageItems          uintptr // [overload("SetStorageItems")] HRESULT  SetStorageItems [in] Windows.Foundation.Collections.IIterable<Windows.Storage.IStorageItem*>* value, [in] boolean readOnly);
	}
)

func (i *IDataPackage) GetProperties(r **IDataPackagePropertySet) error {
	return call(i.VTBL.Properties, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(r)))
}

func (i *IDataPackage) SetText(s string) error {
	sw, err := ole.NewHString(s)
	if err != nil {
		return err
	}
	return call(i.VTBL.SetText, uintptr(unsafe.Pointer(i)), uintptr(sw))
}

func (i *IDataPackage) SetUri(s string) error {
	return nil
}

func (i *IDataPackage) GetIDataPackage2(r **IDataPackage2) error {
	return call(i.VTBL.QueryInterface, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(IDataPackage2GUID)), uintptr(unsafe.Pointer(r)))
}

type (
	IDataPackage2 struct {
		VTBL *IDataPackage2VTBL
	}
	IDataPackage2VTBL struct {
		ole.IInspectableVtbl
		SetApplicationLink uintptr // HRESULT SetApplicationLink([in] Windows.Foundation.Uri* value);
		SetWebLink         uintptr // HRESULT SetWebLink([in] Windows.Foundation.Uri* value);
	}
)

func (i *IDataPackage2) SetWebLink(uri *IUriRuntimeClass) error {
	return call(i.VTBL.SetWebLink, uintptr(unsafe.Pointer(i)), uintptr(unsafe.Pointer(uri)))
}

/*
	IDataPackagePropertySet: http://definitelytyped.org/docs/winrt--winrt/interfaces/windows.applicationmodel.datatransfer.idatapackagepropertyset.html
*/

var IDataPackagePropertySetGUID = ole.NewGUID("CD1C93EB-4C4C-443A-A8D3-F5C241E91689")

type IDataPackagePropertySet struct {
	VTBL *IDataPackagePropertySetVTBL
}

type IDataPackagePropertySetVTBL struct {
	ole.IInspectableVtbl
	GetTitle                 uintptr // [propget] HRESULT Title([out] [retval] HSTRING* value);
	SetTitle                 uintptr // [propput] HRESULT Title([in] HSTRING value);
	GetDescription           uintptr // [propget] HRESULT Description([out] [retval] HSTRING* value);
	SetDescription           uintptr // [propput] HRESULT Description([in] HSTRING value);
	GetThumbnail             uintptr // [propget] HRESULT Thumbnail([out] [retval] Windows.Storage.Streams.IRandomAccessStreamReference** value);
	SetThumbnail             uintptr // [propput] HRESULT Thumbnail([in] Windows.Storage.Streams.IRandomAccessStreamReference* value);
	GetFileTypes             uintptr // [propget] HRESULT FileTypes([out] [retval] Windows.Foundation.Collections.IVector<HSTRING>** value);
	GetApplicationName       uintptr // [propget] HRESULT ApplicationName([out] [retval] HSTRING* value);
	SetApplicationName       uintptr // [propput] HRESULT ApplicationName([in] HSTRING value);
	GetApplicationListingUri uintptr // [propget] HRESULT ApplicationListingUri([out] [retval] Windows.Foundation.Uri** value);
	SetApplicationListingUri uintptr // [propput] HRESULT ApplicationListingUri([in] Windows.Foundation.Uri* value);
}

func (i *IDataPackagePropertySet) SetTitle(s string) error {
	sw, err := ole.NewHString(s)
	if err != nil {
		return err
	}
	// @TODO memory leak
	return call(i.VTBL.SetTitle, uintptr(unsafe.Pointer(i)), uintptr(sw))
}

/*
	ITypedEventHandler<ABI::Windows::ApplicationModel::DataTransfer::DataTransferManager*, ABI::Windows::ApplicationModel::DataTransfer::DataRequestedEventArgs*>: https://docs.microsoft.com/en-us/previous-versions/hh438424(v=vs.85)
*/

var ITypedEventHandlerDataTransferManagerDataRequestedEventArgsGUID = ole.NewGUID("ec6f9cc8-46d0-5e0e-b4d2-7d7773ae37a0")

/*
	ITypedEventHandler: https://docs.microsoft.com/en-us/previous-versions/hh438424(v=vs.85)
	It's a generic type.
*/

type (
	ITypedEventHandler struct {
		VTBL   *ITypedEventHandlerVTBL
		Invoke func(transfer *IDataTransferManager, args *IDataRequestedEventArgs) int
	}
	ITypedEventHandlerVTBL struct {
		ole.IUnknownVtbl
		Invoke uintptr
	}
)

var _ITypedEventHandlerVTBLImpl = &ITypedEventHandlerVTBL{
	IUnknownVtbl: ole.IUnknownVtbl{
		QueryInterface: syscall.NewCallback(func(source, idp unsafe.Pointer, result *unsafe.Pointer) uintptr {
			return ole.S_OK
		}),
		AddRef: syscall.NewCallback(func(_ uintptr) int {
			return 2 // That is wrong, but works.
		}),
		Release: syscall.NewCallback(func(_ uintptr) int {
			return 1 // That is wrong, but works.
		}),
	},
	Invoke: syscall.NewCallback(func(self *ITypedEventHandler, transfer *IDataTransferManager, args *IDataRequestedEventArgs) int {
		return self.Invoke(transfer, args)
	}),
}

func NewTypedEventHandler(fn func(transfer *IDataTransferManager, args *IDataRequestedEventArgs) int) *ITypedEventHandler {
	return &ITypedEventHandler{
		VTBL:   _ITypedEventHandlerVTBLImpl,
		Invoke: fn,
	}
}

func call(trap uintptr, args ...uintptr) error {
	hr, _, _ := syscall.SyscallN(trap, args...)
	if hr != 0 {
		return ole.NewError(hr)
	}
	return nil
}
