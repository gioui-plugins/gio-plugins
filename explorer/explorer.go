package explorer

import (
	"errors"
	"io"
	"strings"
	"sync"

	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

var (
	// ErrUserDecline is returned when the user doesn't select the file.
	ErrUserDecline = errors.New("user exited the file selector without selecting a file")
	// ErrNotAvailable is return when the current OS isn't supported.
	ErrNotAvailable = errors.New("current OS not supported")
)

type driver struct { //todo need rewrite method
}

func (d driver) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	return nil, nil
}

func (d driver) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	return nil, nil
}

type Explorer struct {
	// driver holds OS-Specific content, it varies for each OS.
	driver
}

func NewExplorer(config Config) *Explorer {
	house := &Explorer{}
	attachDriver(house, config)
	return house
}

func attachDriver(house *Explorer, config Config) {

}

func (e *Explorer) Configure(config Config) {
	configureDriver(&e.driver, config)
}

func configureDriver(d *driver, config Config) {

}

func (e *Explorer) OpenFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	return e.driver.openFile(mimes)
}

func (e *Explorer) SaveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	return e.driver.saveFile(name, mime)
}

var stringBuilderPool = sync.Pool{New: func() any { return &strings.Builder{} }}

type result[T any] struct {
	file  T
	error error
}
