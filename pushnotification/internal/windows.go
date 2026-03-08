//go:build windows

package internal

import (
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
)

// ── IPushNotificationChannelManagerStatics ────────────────────────────────────
//
// Windows.Networking.PushNotifications.PushNotificationChannelManager
//
// Vtable (IInspectableVtbl = QI+AddRef+Release+GetIids+GetRuntimeClassName+GetTrustLevel):
//  6  CreatePushNotificationChannelForApplicationAsync
//  7  CreatePushNotificationChannelForApplicationAsyncWithApplicationId  ← use this
//  8  CreatePushNotificationChannelForSecondaryTileAsync

var (
	iPushNotificationChannelManagerStaticsCLSID = "Windows.Networking.PushNotifications.PushNotificationChannelManager"
	iPushNotificationChannelManagerStaticsGUID  = ole.NewGUID("8BAF9B65-77A1-4588-BD19-861529A9DCF0")
)

type iPushNotificationChannelManagerStaticsVtbl struct {
	ole.IInspectableVtbl
	CreateChannelForAppAsync       uintptr // slot 6
	CreateChannelForAppAsyncWithId uintptr // slot 7 ← we use this
	CreateChannelForTileAsync      uintptr // slot 8
}

type iPushNotificationChannelManagerStatics struct {
	vtbl *iPushNotificationChannelManagerStaticsVtbl
}

func newIPushNotificationChannelManagerStatics(r **iPushNotificationChannelManagerStatics) error {
	ins, err := ole.RoGetActivationFactory(iPushNotificationChannelManagerStaticsCLSID, iPushNotificationChannelManagerStaticsGUID)
	if err != nil {
		return fmt.Errorf("RoGetActivationFactory: %w", err)
	}
	*r = (*iPushNotificationChannelManagerStatics)(unsafe.Pointer(ins))
	return nil
}

// ── IAsyncInfo ────────────────────────────────────────────────────────────────

var iAsyncInfoGUID = ole.NewGUID("00000036-0000-0000-C000-000000000046")

type iAsyncInfoVtbl struct {
	ole.IInspectableVtbl
	GetId        uintptr // slot 6
	GetStatus    uintptr // slot 7
	GetErrorCode uintptr // slot 8
	Cancel       uintptr // slot 9
	Close        uintptr // slot 10
}

type iAsyncInfo struct {
	vtbl *iAsyncInfoVtbl
}

type AsyncStatus int32

const (
	AsyncStatusStarted   AsyncStatus = 0
	AsyncStatusCompleted AsyncStatus = 1
	AsyncStatusCanceled  AsyncStatus = 2
	AsyncStatusError     AsyncStatus = 3
)

func (a *iAsyncInfo) status() (AsyncStatus, error) {
	var s AsyncStatus
	hr, _, _ := syscall.SyscallN(
		a.vtbl.GetStatus,
		uintptr(unsafe.Pointer(a)),
		uintptr(unsafe.Pointer(&s)),
	)
	return s, hrErr(hr)
}

func (a *iAsyncInfo) errorCode() error {
	var e uintptr
	hr, _, _ := syscall.SyscallN(
		a.vtbl.GetErrorCode,
		uintptr(unsafe.Pointer(a)),
		uintptr(unsafe.Pointer(&e)),
	)
	if err := hrErr(hr); err != nil {
		return err
	}
	return hrErr(e)
}

// ── IAsyncOperation<PushNotificationChannel> ──────────────────────────────────
//
// NOTE: This is IAsyncOperation<T>, NOT IAsyncOperationWithProgress<T,P>
// so there are NO put_Progress/get_Progress slots.
//
//  6  put_Completed
//  7  get_Completed
//  8  GetResults       ← returns IPushNotificationChannel directly

var iAsyncOperationPushChannelGUID = ole.NewGUID("904103F1-7A8C-4F68-BB55-B59FC3A2D43D")

type iAsyncOperationVtbl struct {
	ole.IInspectableVtbl
	PutCompleted uintptr // slot 6
	GetCompleted uintptr // slot 7
	GetResults   uintptr // slot 8
}

type iAsyncOperation struct {
	vtbl *iAsyncOperationVtbl
}

func (a *iAsyncOperation) queryAsyncInfo() (*iAsyncInfo, error) {
	var info *iAsyncInfo
	hr, _, _ := syscall.SyscallN(
		a.vtbl.QueryInterface,
		uintptr(unsafe.Pointer(a)),
		uintptr(unsafe.Pointer(iAsyncInfoGUID)),
		uintptr(unsafe.Pointer(&info)),
	)
	return info, hrErr(hr)
}

func (a *iAsyncOperation) getResults() (*iPushNotificationChannel, error) {
	var result *iPushNotificationChannel
	hr, _, _ := syscall.SyscallN(
		a.vtbl.GetResults,
		uintptr(unsafe.Pointer(a)),
		uintptr(unsafe.Pointer(&result)),
	)
	return result, hrErr(hr)
}

// ── IPushNotificationChannel ──────────────────────────────────────────────────
//
//  6  get_Uri
//  7  get_ExpirationTime

type iPushNotificationChannelVtbl struct {
	ole.IInspectableVtbl
	GetUri            uintptr // slot 6
	GetExpirationTime uintptr // slot 7
}

type iPushNotificationChannel struct {
	vtbl *iPushNotificationChannelVtbl
}

func (c *iPushNotificationChannel) uri() (string, error) {
	var h ole.HString
	hr, _, _ := syscall.SyscallN(
		c.vtbl.GetUri,
		uintptr(unsafe.Pointer(c)),
		uintptr(unsafe.Pointer(&h)),
	)
	if err := hrErr(hr); err != nil {
		return "", err
	}
	return h.String(), nil
}

// ── Public API ────────────────────────────────────────────────────────────────

type requestURI struct {
	done chan result
}

type result struct {
	uri string
	err error
}

var getURIChan = make(chan requestURI)

func init() {
	go func() {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		if err := ole.RoInitialize(1); err != nil {
			panic(fmt.Sprintf("RoInitialize: %v", err))
		}

		for req := range getURIChan {
			uri, err := getChannel()
			req.done <- result{uri: uri, err: err}
		}
	}()
}

// GetChannelURI returns the WNS channel URI for this device.
// Requires app to be packaged (MSIX) for package identity.
func GetChannelURI() (string, error) {
	req := requestURI{
		done: make(chan result, 1),
	}
	getURIChan <- req
	res := <-req.done
	return res.uri, res.err
}

func getChannel() (string, error) {
	// 1. Get activation factory
	var statics *iPushNotificationChannelManagerStatics
	if err := newIPushNotificationChannelManagerStatics(&statics); err != nil {
		return "", fmt.Errorf("newIPushNotificationChannelManagerStatics: %w", err)
	}

	// 2. Call CreateChannelForApplicationAsync (slot 6) — uses the current app
	var op *iAsyncOperation
	hr, _, _ := syscall.SyscallN(statics.vtbl.CreateChannelForAppAsync, uintptr(unsafe.Pointer(statics)), uintptr(unsafe.Pointer(&op)))
	if err := hrErr(hr); err != nil {
		return "", fmt.Errorf("CreateChannelForAppAsync: %w", err)
	}

	// 4. QI for IAsyncInfo to poll status
	info, err := op.queryAsyncInfo()
	if err != nil {
		return "", fmt.Errorf("QueryInterface IAsyncInfo: %w", err)
	}

	// 5. Poll until done
	for {
		s, err := info.status()
		if err != nil {
			return "", fmt.Errorf("get_Status: %w", err)
		}
		switch s {
		case AsyncStatusCompleted:
			break
		case AsyncStatusError:
			return "", fmt.Errorf("async failed: %w", info.errorCode())
		case AsyncStatusCanceled:
			return "", fmt.Errorf("async cancelled")
		default:
			runtime.Gosched()
			continue
		}

		// 6. Get results
		ch, err := op.getResults()
		if err != nil {
			return "", fmt.Errorf("GetResults: %w", err)
		}

		return ch.uri()
	}
}

func hrErr(hr uintptr) error {
	if int32(hr) >= 0 {
		return nil
	}
	return fmt.Errorf("HRESULT 0x%08X", uint32(hr))
}
