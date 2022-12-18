# explorer

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/explorer.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/explorer)

Opens the native file-dialog/file-picker.

--------------

## Usage

To open an single file, you can use `explorer.OpenFileOp`.

That will open native File Dialog/File Picker. Once the file is selected by the end-user,
one `explorer.OpenFileEvent` will be sent to the given `Tag`. 

```go
explorer.OpenFileOp{
    Tag: yourTag, 
    Mimetype: []mimetype.MimeType{
      {Extension: "png", Type: "image", Subtype: "png"},
      {Extension: "jpg", Type: "image", Subtype: "jpeg"},
      {Extension: "jpeg", Type: "image", Subtype: "jpeg"},
      {Extension: "gif", Type: "image", Subtype: "gif"},
      {Extension: "webp", Type: "image", Subtype: "webp"},
    },
}.Add(gtx.Ops)
```

### Operations:

Operations must be added with `.Add(gtx.Ops)` method. The operation will be executed at the end of the frame.

- `explorer.OpenFileOp`:
  - Opens the native file dialog to open/import a single file.
- `explorer.SaveFileOp`:
  - Open the native file dialog to save/export a single file.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `explorer.OpenFileEvent`:
  - Sent to `Tag` when the user chooses the file to be read/open. That event contains one io.ReadCloser.
- `explorer.SaveFileEvent`:
  - Sent to `Tag` when the user chooses the file to save/replace. That event contains one io.WriteCloser.
- `explorer.ErrorEvent`:
  - Sent to `Tag` when some error occurs.
- `explorer.CancelEvent`: 
  - Sent to `Tag` when the user closes the file-dialog or not select one valid file.

## Features

| Features | Windows | Android | MacOS | iOS | WebAssembly | FreeBSD |  Linux |
| -- | -- | -- | -- | -- | -- |  -- |  -- |
| Import File |✔|✔|✔|✔|✔*|❌|❌|
| Export File |✔|✔|✔|✔|✔*|❌|❌|

On WASM it contains two implementations, one using File System Access API and another
using the basic HTMLInputElement and such.