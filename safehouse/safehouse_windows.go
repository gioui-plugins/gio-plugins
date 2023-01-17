package safehouse

import (
	"golang.org/x/sys/windows"
	"runtime"
	"strings"
	"time"
	"unsafe"
)

var (
	_wincred = windows.NewLazySystemDLL("Advapi32.dll")

	_write  = _wincred.NewProc("CredWriteW")
	_read   = _wincred.NewProc("CredReadW")
	_list   = _wincred.NewProc("CredEnumerateW")
	_delete = _wincred.NewProc("CredDeleteW")
)

var (
	_Cred_Type_Generic uint32 = 0x01

	_Cred_Persist_Local_Machine uint32 = 0x02
)

type _Credential_AttributeW struct {
	Keyword   *uint8
	Flags     uint32
	ValueSize uint32
	Value     *byte
}

type _CredentialW struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint8
	Comment            *uint8
	LastWritten        uint64
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttributeCount     uint32
	Attributes         *_Credential_AttributeW
	TargetAlias        *uint8
	UserName           *uint8
}

type driver struct {
	prefix string
}

func attachDriver(house *SafeHouse, config Config) {
	house.driver = driver{prefix: config.App}
}

func (d driver) setSecret(secret Secret) error {
	id := d.keyFor(secret.Identifier) + "\x00"
	if len(id) > 32767 {
		return ErrIdentifierMaxLength
	}

	if len(secret.Data) > 5*512 {
		return Err
	}

	var iw *uint8
	if len(id) > 0 {
		iw = ((*[2]*uint8)(unsafe.Pointer(&id)))[0]
	}
	defer runtime.KeepAlive(&id)

	var dw *byte
	if len(secret.Data) > 0 {
		dw = ((*[3]*byte)(unsafe.Pointer(&secret.Data)))[0]
	}
	defer runtime.KeepAlive(&secret.Data)

	cred := _CredentialW{
		Flags:              0,
		Type:               _Cred_Type_Generic,
		TargetName:         iw,
		Comment:            nil,
		LastWritten:        uint64(time.Now().UnixNano()),
		CredentialBlobSize: uint32(len(secret.Data)),
		CredentialBlob:     dw,
		Persist:            _Cred_Persist_Local_Machine,
		AttributeCount:     0,
		Attributes:         nil,
		TargetAlias:        nil,
		UserName:           nil,
	}
	defer runtime.KeepAlive(&cred)

	_write.Call(uintptr(unsafe.Pointer(&cred)), uintptr(0))
}

func (d driver) listSecret(looper Looper) error {

}

func (d driver) getSecret(identifier string, secret *Secret) error {

}

func (d driver) removeSecret(identifier string) error {
	//TODO implement me
	panic("implement me")
}

func (d driver) keyFor(s string) string {
	return d.prefix + s
}

func (d driver) rawKeyFor(s string) string {
	if strings.HasPrefix(s, d.prefix) && len(s) > len(d.prefix) {
		return d.prefix[len(d.prefix):]
	}
	return s
}
