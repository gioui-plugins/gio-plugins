package explorer

import (
	"io"
	"strings"
	"syscall/js"

	"github.com/gioui-plugins/gio-plugins/explorer/mimetype"
)

type driver struct{}

func attachDriver(house *Explorer, config Config) {
	house.driver = driver{}
	configureDriver(&house.driver, config)
}

func configureDriver(driver *driver, config Config) {}

func (e *driver) saveFile(name string, mime mimetype.MimeType) (io.WriteCloser, error) {
	if js.Global().Get("showSaveFilePicker").IsUndefined() {
		return newFileWriterLegacy(name), nil
	}

	filter := js.Global().Get("Object").New()
	filterTypes := js.Global().Get("Array").New()
	filterMime := js.Global().Get("Object").New()
	filterMime.Set(mime.Type+"/*", js.Global().Get("Array").New())
	filter.Set("types", filterTypes)
	filter.Set("suggestedName", name)

	pickerArgs, ok := await(js.Global().Call("showSaveFilePicker", filter))
	if !ok || len(pickerArgs) == 0 {
		return nil, ErrUserDecline
	}
	pickerSelection := pickerArgs[0]

	writable, ok := await(pickerSelection.Call("createWritable"))
	if !ok || len(writable) == 0 {
		return nil, ErrUserDecline
	}
	return newFileWriter(writable[0]), nil
}

func (e *driver) openFile(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	if js.Global().Get("showOpenFilePicker").IsUndefined() {
		return e.openFileLegacy(mimes)
	}

	filter := js.Global().Get("Object").New()
	filterAccept := js.Global().Get("Object").New()

	for _, v := range mimes {
		t := v.Type + "/*"
		if filterAccept.Get(t).IsUndefined() {
			filterAccept.Set(t, js.Global().Get("Array").New())
		}
		list := filterAccept.Get(t)
		list.SetIndex(list.Length(), v.Extension)
	}

	filter.Set("Accept", filterAccept)

	pickerArgs, ok := await(js.Global().Call("showOpenFilePicker", filter))
	if !ok || len(pickerArgs) == 0 || pickerArgs[0].Length() == 0 {
		return nil, ErrUserDecline
	}
	pickerSelection := pickerArgs[0].Index(0)

	file, ok := await(pickerSelection.Call("getFile"))
	if !ok || len(file) == 0 {
		return nil, ErrUserDecline
	}
	return newFileReader(file[0]), nil
}

func (e *driver) openFileLegacy(mimes []mimetype.MimeType) (io.ReadCloser, error) {
	res := make(chan result)
	callback := openCallbackLegacy(res)

	extensions := stringBuilderPool.Get().(*strings.Builder)
	for i, v := range mimes {
		if i > 0 {
			extensions.WriteString(",")
		}
		extensions.WriteString(v.Extension)
	}
	defer stringBuilderPool.Put(extensions)
	defer extensions.Reset()

	document := js.Global().Get("document")
	input := document.Call("createElement", "input")
	input.Call("addEventListener", "change", callback)
	input.Call("addEventListener", "cancel", callback)
	input.Set("type", "file")
	input.Set("style", "display:none;")
	if extensions.Len() > 0 {
		input.Set("accept", extensions.String())
	}
	document.Get("body").Call("appendChild", input)
	input.Call("click")

	file := <-res
	if file.error != nil {
		return nil, file.error
	}
	return file.file.(io.ReadCloser), nil
}

func await(value js.Value) ([]js.Value, bool) {
	res := make(chan []js.Value, 1)

	var s, f js.Func
	s = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer s.Release()

		res <- args
		return nil
	})
	f = js.FuncOf(func(this js.Value, args []js.Value) any {
		defer f.Release()

		res <- nil
		return nil
	})

	go value.Call("then", s).Call("catch", f)

	r := <-res
	if r == nil {
		return nil, false
	}
	return r, true
}

func openCallbackLegacy(r chan result) js.Func {
	var fn js.Func
	fn = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		defer fn.Release()

		files := args[0].Get("target").Get("files")
		if files.Length() <= 0 {
			r <- result{error: ErrUserDecline}
			return nil
		}
		r <- result{file: newFileReader(files.Index(0))}
		return nil
	})
	return fn
}

var (
	_ io.ReadCloser  = (*FileReader)(nil)
	_ io.WriteCloser = (*FileWriterLegacy)(nil)
	_ io.WriteCloser = (*FileWriter)(nil)
)
