package explorer

import (
	"errors"
	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
	"io"
	"strings"
	"sync"
)

var (
	// ErrUserDecline is returned when the user doesn't select the file.
	ErrUserDecline = errors.New("user exited the file selector without selecting a file")
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")
)

type Explorer struct {
	// driver holds OS-Specific content, it varies for each OS.
	driver
}

func NewExplorer(config Config) *Explorer {
	house := &Explorer{}
	attachDriver(house, config)
	return house
}

func (e *Explorer) Configure(config Config) {
	configureDriver(&e.driver, config)
}

func (e *Explorer) OpenFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	return e.driver.openFile(mimes)
}

func (e *Explorer) SaveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	return e.driver.saveFile(name, mime)
}

var stringBuilderPool = sync.Pool{New: func() any { return &strings.Builder{} }}

type result struct {
	file  interface{}
	error error
}
