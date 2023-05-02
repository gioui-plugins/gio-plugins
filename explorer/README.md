# explorer

[![Go Reference](https://pkg.go.dev/badge/github.com/gio-plugins/gio-plugin/explorer.svg)](https://pkg.go.dev/github.com/gio-plugins/gio-plugin/explorer)

Opens the native file-dialog/file-picker.

--------------

## Usage

## Freestanding

- `exporer.NewExplorer`:
  - Creates a instance of Explorer struct, given the config.
- `explorer.Configure`:
  - Updates the current Explorer with the given config.
- `explorer.OpenFile`:
  - Opens the native file dialog to open a single file.
- `explorer.SaveFile`:
  - Opens the native file dialog to save a single file.

## Gio

## Non-Plugin:

If you want to use it without plugin, read the Freestanding instructions. We provide some helper functions, such as NewConfigFromViewEvent and such.

## Plugin:

To open an single file, you can use `explorer.OpenFileOp`.

That will open native File Dialog/File Picker. Once the file is selected by the end-user,
one `explorer.OpenFileEvent` will be sent to the given `Tag`. 

```go
gioexplorer.OpenFileOp{
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

- `gioexplorer.OpenFileOp`:
  - Opens the native file dialog to open/import a single file.
- `gioexplorer.SaveFileOp`:
  - Open the native file dialog to save/export a single file.

## Events:

Events are response sent using the `Tag` and should be handled with `gtx.Events()`.

- `gioexplorer.OpenFileEvent`:
  - Sent to `Tag` when the user chooses the file to be read/open. That event contains one io.ReadCloser.
- `gioexplorer.SaveFileEvent`:
  - Sent to `Tag` when the user chooses the file to save/replace. That event contains one io.WriteCloser.
- `gioexplorer.ErrorEvent`:
  - Sent to `Tag` when some error occurs.
- `gioexplorer.CancelEvent`: 
  - Sent to `Tag` when the user closes the file-dialog or not select one valid file.

## Features

| Features | Windows | Android | MacOS | iOS | WebAssembly | FreeBSD |  Linux |
| -- | -- | -- | -- | -- | -- |  -- |  -- |
| Import File |✔|✔|✔|✔|✔*|❌|❌|
| Export File |✔|✔|✔|✔|✔*|❌|❌|

On WASM it contains two implementations, one using File System Access API and another
using the basic HTMLInputElement and such.