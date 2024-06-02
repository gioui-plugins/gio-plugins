package giowebview

import (
	"gioui.org/app"
	"github.com/gioui-plugins/gio-plugins/plugin"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
	"reflect"
)

var wantEvent = []reflect.Type{
	reflect.TypeOf(plugin.ViewEvent{}),
	reflect.TypeOf(app.FrameEvent{}),
	reflect.TypeOf(app.DestroyEvent{}),
	reflect.TypeOf(plugin.EndFrameEvent{}),
}

// CookiesEvent is the event sent when ListCookieCmd is executed.
type CookiesEvent struct {
	Cookies []webview.CookieData
}

// StorageEvent is the event sent when ListStorageCmd is executed.
type StorageEvent struct {
	Storage []webview.StorageData
}

// MessageEvent is the event sent when receiving a message,
// from previously defined MessageReceiverCmd.
type MessageEvent struct {
	Message string
}

// NavigationEvent is issued when the webview change the URL.
type NavigationEvent webview.NavigationEvent

// TitleEvent is issued when the webview change the title.
type TitleEvent webview.TitleEvent

// ErrorEvent is issued when the webview encounters an error.
type ErrorEvent struct {
	error
}

// ClearCacheEvent is issued when clearing the storage is completed.
type ClearCacheEvent struct {
	Error error
}

// SetCookieEvent is issued when setting a cookie is completed.
type SetCookieEvent struct {
	Error error
}

// SetCookieArrayEvent is issued when setting a cookie is completed.
type SetCookieArrayEvent struct {
	Error error
}

// RemoveCookieEvent is issued when removing a cookie is completed.
type RemoveCookieEvent struct {
	Error error
}

// SetStorageEvent is issued when setting a storage is completed.
type SetStorageEvent struct {
	Error error
}

// RemoveStorageEvent is issued when removing a storage is completed.
type RemoveStorageEvent struct {
	Error error
}

func (c StorageEvent) ImplementsEvent()      {}
func (c CookiesEvent) ImplementsEvent()      {}
func (c MessageEvent) ImplementsEvent()      {}
func (NavigationEvent) ImplementsEvent()     {}
func (TitleEvent) ImplementsEvent()          {}
func (ErrorEvent) ImplementsEvent()          {}
func (ClearCacheEvent) ImplementsEvent()     {}
func (SetCookieEvent) ImplementsEvent()      {}
func (SetCookieArrayEvent) ImplementsEvent() {}
func (RemoveCookieEvent) ImplementsEvent()   {}
func (SetStorageEvent) ImplementsEvent()     {}
func (RemoveStorageEvent) ImplementsEvent()  {}
