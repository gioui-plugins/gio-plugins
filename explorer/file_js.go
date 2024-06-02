package explorer

import (
	"io"
	"syscall/js"
)

type FileReader struct {
	buffer                   js.Value
	isClosed                 bool
	index                    uint32
	callback                 chan js.Value
	successFunc, failureFunc js.Func
}

func newFileReader(v js.Value) *FileReader {
	f := &FileReader{
		buffer:   v,
		callback: make(chan js.Value, 1),
	}
	f.successFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.callback <- args[0]
		return nil
	})
	f.failureFunc = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		f.callback <- js.Undefined()
		return nil
	})

	return f
}

func (f *FileReader) Read(b []byte) (n int, err error) {
	if f == nil || f.isClosed {
		return 0, io.ErrClosedPipe
	}

	go fileSlice(f.index, f.index+uint32(len(b)), f.buffer, f.successFunc, f.failureFunc)

	buffer := <-f.callback
	n32 := fileRead(buffer, b)
	if n32 == 0 {
		return 0, io.EOF
	}
	f.index += n32

	return int(n32), err
}

func (f *FileReader) Close() error {
	if f == nil || f.isClosed {
		return io.ErrClosedPipe
	}

	f.failureFunc.Release()
	f.successFunc.Release()
	f.isClosed = true
	return nil
}

type FileWriterLegacy struct {
	buffer   js.Value
	isClosed bool
	name     string
}

func newFileWriterLegacy(name string) *FileWriterLegacy {
	return &FileWriterLegacy{
		name:   name,
		buffer: js.Global().Get("Uint8Array").New(),
	}
}

func (f *FileWriterLegacy) Write(b []byte) (n int, err error) {
	if f == nil || f.isClosed {
		return 0, io.ErrClosedPipe
	}
	if len(b) == 0 {
		return 0, nil
	}

	fileWrite(f.buffer, b)
	return len(b), err
}

func (f *FileWriterLegacy) Close() error {
	if f == nil || f.isClosed {
		return io.ErrClosedPipe
	}
	f.isClosed = true
	return f.saveFile()
}

func (f *FileWriterLegacy) saveFile() error {
	config := js.Global().Get("Object").New()
	config.Set("type", "octet/stream")

	blob := js.Global().Get("Blob").New(
		js.Global().Get("Array").New().Call("concat", f.buffer),
		config,
	)

	document := js.Global().Get("document")
	anchor := document.Call("createElement", "a")
	anchor.Set("download", f.name)
	anchor.Set("href", js.Global().Get("URL").Call("createObjectURL", blob))
	document.Get("body").Call("appendChild", anchor)
	anchor.Call("click")

	return nil
}

type FileWriter struct {
	writable js.Value
	isClosed bool

	wait    chan bool
	success js.Func
	fail    js.Func
}

func newFileWriter(v js.Value) *FileWriter {
	wait := make(chan bool, 1)
	return &FileWriter{
		writable: v,
		wait:     wait,
		success: js.FuncOf(func(this js.Value, args []js.Value) any {
			wait <- true
			return nil
		}),
		fail: js.FuncOf(func(this js.Value, args []js.Value) any {
			wait <- false
			return nil
		}),
	}
}

func (f *FileWriter) Write(b []byte) (n int, err error) {
	if f.isClosed {
		return 0, io.ErrClosedPipe
	}

	go writableWrite(f.writable, f.success.Value, f.fail.Value, b)
	if ok := <-f.wait; !ok {
		return 0, io.ErrUnexpectedEOF
	}
	return len(b), nil
}

func (f *FileWriter) Close() error {
	if f.isClosed {
		return nil
	}
	f.success.Release()
	f.fail.Release()
	f.writable.Call("close")
	return nil
}

//go:wasmimport gojs github.com/gioui-plugins/gio-plugins/explorer.fileRead
func fileRead(value js.Value, b []byte) uint32

//go:wasmimport gojs github.com/gioui-plugins/gio-plugins/explorer.fileWrite
func fileWrite(value js.Value, b []byte)

//go:wasmimport gojs github.com/gioui-plugins/gio-plugins/explorer.fileSlice
func fileSlice(start, end uint32, value js.Value, success, failure js.Func)

//go:wasmimport gojs github.com/gioui-plugins/gio-plugins/explorer.writableWrite
func writableWrite(writable js.Value, success js.Value, failure js.Value, b []byte)
