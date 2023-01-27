package safedata

import (
	"encoding/base64"
	"hash"
	"runtime"
	"time"
	"unsafe"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/sys/windows"
)

var (
	_wincred = windows.NewLazySystemDLL("advapi32.dll")

	_write  = _wincred.NewProc("CredWriteW")
	_read   = _wincred.NewProc("CredReadW")
	_list   = _wincred.NewProc("CredEnumerateW")
	_delete = _wincred.NewProc("CredDeleteW")
	_free   = _wincred.NewProc("CredFree")
)

var (
	_Cred_Type_Generic uint32 = 0x01

	_Cred_Persist_Local_Machine uint32 = 0x02

	_Cred_Max_Credential_Blob_Size = 5 * 512
	_Cred_Max_Comment_Length       = 256
	_Cred_Max_Target_Length        = 32767
)

type _CredentialW struct {
	Flags              uint32
	Type               uint32
	TargetName         *uint16
	Comment            *uint16
	LastWritten        uint64
	CredentialBlobSize uint32
	CredentialBlob     *byte
	Persist            uint32
	AttributeCount     uint32
	Attributes         uintptr
	TargetAlias        *uint16
	UserName           *uint16
}

type driver struct {
	prefix string
}

func attachDriver(house *SafeData, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {
	driver.prefix = config.App
}

func (d driver) setSecret(secret Secret) error {
	sizeExtension := uint64(len(secret.Data) / _Cred_Max_Credential_Blob_Size)
	if len(secret.Data)%_Cred_Max_Credential_Blob_Size == 0 {
		sizeExtension -= 1
	}

	if sizeExtension > 0 {
		if err := d.deleteExtended(secret.Identifier, sizeExtension-1); err != nil && err != ErrNotFound {
			return err
		}
	}

	id, err := d.encodeUTF16(d.keyFor(secret.Identifier), _Cred_Max_Target_Length)
	if err != nil {
		return err
	}

	desc, err := d.encodeUTF16(secret.Description, _Cred_Max_Comment_Length)
	if err != nil {
		return err
	}

	defer runtime.KeepAlive(id)
	defer runtime.KeepAlive(desc)

	dw := secret.Data
	if len(secret.Data) == 0 {
		dw = []byte{0x00}
	}

	defer runtime.KeepAlive(secret.Data)

	size := len(secret.Data)
	if size > _Cred_Max_Credential_Blob_Size {
		size = _Cred_Max_Credential_Blob_Size
	}

	now := (time.Now().UTC().Unix() / 100) + 116444736000000000
	cred := _CredentialW{
		Flags:              0,
		Type:               _Cred_Type_Generic,
		TargetName:         id,
		LastWritten:        uint64(now&0xffffffff | (now >> 32 & 0xffffffff)),
		Comment:            desc,
		CredentialBlobSize: uint32(size),
		CredentialBlob:     &dw[0],
		Persist:            _Cred_Persist_Local_Machine,
	}

	defer runtime.KeepAlive(cred)

	r, _, err := _write.Call(uintptr(unsafe.Pointer(&cred)), uintptr(0))
	if r == 0 {
		return err
	}

	hash := d.hasherFor(secret.Identifier)
	for i := uint64(0); i < sizeExtension; i++ {
		hash.Reset()
		hash.Write((*[8]byte)(unsafe.Pointer(&i))[:])
		id := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		cred.TargetName, _ = d.encodeUTF16("_"+d.keyFor(id), _Cred_Max_Target_Length)
		cred.CredentialBlob = &dw[uint64(_Cred_Max_Credential_Blob_Size)*(i+1)]
		cred.CredentialBlobSize = uint32(len(secret.Data)) - uint32(uint64(_Cred_Max_Credential_Blob_Size)*(i+1))
		if cred.CredentialBlobSize > uint32(_Cred_Max_Credential_Blob_Size) {
			cred.CredentialBlobSize = uint32(_Cred_Max_Credential_Blob_Size)
		}

		r, _, err := _write.Call(uintptr(unsafe.Pointer(&cred)), uintptr(0))
		if r == 0 {
			return err
		}
	}

	return nil
}

func (d driver) listSecret(looper Looper) error {
	id, err := d.encodeUTF16(d.keyFor("*"), _Cred_Max_Target_Length)
	if err != nil {
		return err
	}

	defer runtime.KeepAlive(&id)

	credentials := make([]*_CredentialW, 0, 0)
	uc := (*[3]uintptr)(unsafe.Pointer(&credentials))

	r, _, err := _list.Call(
		uintptr(unsafe.Pointer(id)),
		0,
		uintptr(unsafe.Pointer(&uc[2])),
		uintptr(unsafe.Pointer(&uc[0])),
	)

	if r == 0 || cap(credentials) == 0 || credentials == nil {
		return ErrNotFound
	}
	credentials = credentials[:cap(credentials)]

	defer _free.Call(uc[0])

	for i := 0; i < len(credentials); i++ {
		looper(d.rawKeyFor(windows.UTF16PtrToString(credentials[i].TargetName)))
	}

	return nil
}

func (d driver) getSecret(identifier string, secret *Secret) error {
	if err := d.read(d.keyFor(identifier), secret); err != nil {
		return err
	}

	hash := d.hasherFor(identifier)
	for i := uint64(0); true; i++ {
		hash.Reset()
		hash.Write((*[8]byte)(unsafe.Pointer(&i))[:])

		if l := len(secret.Data); _Cred_Max_Credential_Blob_Size > cap(secret.Data)-l {
			secret.Data = append(secret.Data, make([]byte, _Cred_Max_Credential_Blob_Size)...)
			secret.Data = secret.Data[:l]
		}

		id := base64.URLEncoding.EncodeToString(hash.Sum(nil))
		s := Secret{Data: secret.Data[len(secret.Data):]}

		if err := d.read("_"+d.keyFor(id), &s); err != nil {
			if err == ErrNotFound {
				return nil
			}
			return err
		}
		secret.Data = secret.Data[:len(secret.Data)+len(s.Data)]
	}

	return nil
}

func (d driver) read(identifier string, secret *Secret) error {
	id, err := d.encodeUTF16(identifier, _Cred_Max_Target_Length)
	if err != nil {
		return err
	}
	defer runtime.KeepAlive(id)

	var cred *_CredentialW
	r, _, err := _read.Call(uintptr(unsafe.Pointer(id)), uintptr(_Cred_Type_Generic), 0, uintptr(unsafe.Pointer(&cred)))
	if r == 0 {
		return ErrNotFound
	}
	defer _free.Call(uintptr(unsafe.Pointer(cred)))

	secret.Identifier = identifier
	secret.Description = windows.UTF16PtrToString(cred.Comment)
	if int(cred.CredentialBlobSize) > cap(secret.Data) {
		secret.Data = make([]byte, cred.CredentialBlobSize)
	}

	h := [3]uintptr{
		uintptr(unsafe.Pointer(cred.CredentialBlob)),
		uintptr(cred.CredentialBlobSize),
		uintptr(cred.CredentialBlobSize),
	}
	secret.Data = secret.Data[:cred.CredentialBlobSize]
	copy(secret.Data, *(*[]byte)(unsafe.Pointer(&h)))

	return nil
}

func (d driver) removeSecret(identifier string) error {
	if err := d.delete(d.keyFor(identifier)); err != nil {
		return err
	}
	return d.deleteExtended(identifier, 0)
}

func (d driver) deleteExtended(identifier string, start uint64) error {
	hash := d.hasherFor(identifier)
	for i := start; true; i++ {
		hash.Reset()
		hash.Write((*[8]byte)(unsafe.Pointer(&i))[:])

		id := base64.URLEncoding.EncodeToString(hash.Sum(nil))

		if err := d.delete("_" + d.keyFor(id)); err != nil {
			if err == ErrNotFound {
				return nil
			}
			return err
		}
	}
	return nil
}

func (d driver) delete(identifier string) error {
	id, err := d.encodeUTF16(identifier, _Cred_Max_Target_Length)
	if err != nil {
		return err
	}
	defer runtime.KeepAlive(id)

	r, _, err := _delete.Call(uintptr(unsafe.Pointer(id)), uintptr(_Cred_Type_Generic), 0)
	if r == 0 {
		return ErrNotFound
	}
	return nil
}

func (d driver) keyFor(s string) string {
	return d.prefix + s
}

func (d driver) rawKeyFor(s string) string {
	if len(s) < len(d.prefix) {
		return s
	}
	return s[len(d.prefix):]
}

func (d driver) hasherFor(s string) hash.Hash {
	key := blake2b.Sum384([]byte(s))
	hash, err := blake2b.New(48, key[:])
	if err != nil {
		panic(err)
	}
	return hash
}

func (d driver) encodeUTF16(s string, max int) (*uint16, error) {
	id, err := windows.UTF16FromString(s)
	if err != nil {
		return nil, ErrMalformedMetadata
	}

	if len(id) > max {
		return nil, ErrMetadataMaxLength
	}
	return &id[0], nil
}
