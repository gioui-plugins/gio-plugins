// SPDX-License-Identifier: Unlicense OR MIT

package explorer

import (
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode/utf16"
	"unsafe"

	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
	"golang.org/x/sys/windows"
)

var (
	// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/
	_Dialog32 = windows.NewLazySystemDLL("comdlg32.dll")

	_GetSaveFileName = _Dialog32.NewProc("GetSaveFileNameW")
	_GetOpenFileName = _Dialog32.NewProc("GetOpenFileNameW")

	// https://docs.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-openfilenamew
	_FlagFileMustExist   = uint32(0x00001000)
	_FlagForceShowHidden = uint32(0x10000000)
	_FlagOverwritePrompt = uint32(0x00000002)
	_FlagDisableLinks    = uint32(0x00100000)

	_FilePathLength       = uint32(65535)
	_OpenFileStructLength = uint32(unsafe.Sizeof(_OpenFileName{}))
)

type (
	// _OpenFileName is defined at https://docs.microsoft.com/pt-br/windows/win32/api/commdlg/ns-commdlg-openfilenamew
	_OpenFileName struct {
		StructSize      uint32
		Owner           uintptr
		Instance        uintptr
		Filter          *uint16
		CustomFilter    *uint16
		MaxCustomFilter uint32
		FilterIndex     uint32
		File            *uint16
		MaxFile         uint32
		FileTitle       *uint16
		MaxFileTitle    uint32
		InitialDir      *uint16
		Title           *uint16
		Flags           uint32
		FileOffset      uint16
		FileExtension   uint16
		DefExt          *uint16
		CustData        uintptr
		FnHook          uintptr
		TemplateName    *uint16
		PvReserved      uintptr
		DwReserved      uint32
		FlagsEx         uint32
	}
)

type explorer struct{}

type explorerPlugin struct {
}

func (e *explorerPlugin) listenEvents(evt event.Event) {
	// NO-OP
}

func (e *explorerPlugin) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	pathUTF16 := make([]uint16, _FilePathLength)
	copy(pathUTF16, windows.StringToUTF16(name))

	filterUTF16, err := buildFilter([]mimetype.MimeType{mime})
	if err != nil {
		return nil, err
	}

	var filter *uint16
	if len(filterUTF16) > 0 {
		filter = &filterUTF16[0]
	}

	open := _OpenFileName{
		File:          &pathUTF16[0],
		MaxFile:       _FilePathLength,
		Filter:        filter,
		FileExtension: uint16(strings.Index(name, filepath.Ext(name))),
		Flags:         _FlagOverwritePrompt,
		StructSize:    _OpenFileStructLength,
	}

	if r, _, _ := _GetSaveFileName.Call(uintptr(unsafe.Pointer(&open))); r == 0 {
		return nil, ErrUserDecline
	}

	path := windows.UTF16ToString(pathUTF16)
	if len(path) == 0 {
		return nil, ErrUserDecline
	}

	runtime.KeepAlive(open)
	runtime.KeepAlive(filterUTF16)
	return os.Create(path)
}

func (e *explorerPlugin) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	pathUTF16 := make([]uint16, _FilePathLength)

	filterUTF16, err := buildFilter(mimes)
	if err != nil {
		return nil, err
	}

	var filter *uint16
	if len(filterUTF16) > 0 {
		filter = &filterUTF16[0]
	}

	open := _OpenFileName{
		File:       &pathUTF16[0],
		MaxFile:    _FilePathLength,
		Filter:     filter,
		Flags:      _FlagFileMustExist | _FlagForceShowHidden | _FlagDisableLinks,
		StructSize: _OpenFileStructLength,
	}

	if r, _, _ := _GetOpenFileName.Call(uintptr(unsafe.Pointer(&open))); r == 0 {
		return nil, ErrUserDecline
	}

	path := windows.UTF16ToString(pathUTF16)
	if len(path) == 0 {
		return nil, ErrUserDecline
	}

	runtime.KeepAlive(open)
	runtime.KeepAlive(filterUTF16)
	return os.Open(path)
}

func buildFilter(mimes []mimetype.MimeType) ([]uint16, error) {
	if len(mimes) <= 0 {
		return nil, nil
	}

	b := stringBuilderPool.Get().(*strings.Builder)
	b.Reset()
	defer stringBuilderPool.Put(b)

	for i, v := range mimes {
		// Extension must have `*` wildcard, so `.jpg` must be `*.jpg`.
		if i > 0 {
			b.WriteString(";")
		}
		if !strings.HasPrefix(v.Extension, "*") {
			b.WriteByte('*')
		}
		b.WriteString(strings.ToUpper(v.Extension))
	}

	// That is a "string-pair", Windows have a Title and the Filter, for instance it could be:
	// Images\0*.JPG;*.PNG\0\0
	// Where `\0` means NULL
	return utf16.Encode([]rune(b.String() + "\x00" + b.String() + "\x00" + "\x00")), nil
}
