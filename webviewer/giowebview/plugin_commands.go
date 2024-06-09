package giowebview

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"github.com/gioui-plugins/gio-plugins/webviewer/webview"
	"net/url"
	"reflect"
)

var wantCommands = []reflect.Type{
	reflect.TypeOf(NavigateCmd{}),
	reflect.TypeOf(DestroyCmd{}),
	reflect.TypeOf(SetCookieCmd{}),
	reflect.TypeOf(SetCookieArrayCmd{}),
	reflect.TypeOf(ListCookieCmd{}),
	reflect.TypeOf(RemoveCookieCmd{}),
	reflect.TypeOf(SetStorageCmd{}),
	reflect.TypeOf(ListStorageCmd{}),
	reflect.TypeOf(RemoveStorageCmd{}),
	reflect.TypeOf(ClearCacheCmd{}),
	reflect.TypeOf(ExecuteJavascriptCmd{}),
	reflect.TypeOf(InstallJavascriptCmd{}),
	reflect.TypeOf(MessageReceiverCmd{}),
}

// NavigateCmd redirects the last Display to the
// given URL. If the URL have unknown protocols,
// or malformed URL may lead to unknown behaviors.
type NavigateCmd struct {
	// View is the webview to redirect.
	View event.Tag
	// URL is the URL to redirect to.
	URL string
}

// DestroyCmd destroys the webview.
type DestroyCmd struct {
	View event.Tag
}

// SetCookieCmd sets given cookie in the webview.
type SetCookieCmd struct {
	// View is the webview to set the cookie.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Cookie is the cookie to set.
	Cookie webview.CookieData
}

// SetCookieArrayCmd sets the given list of cookies in the webview.
type SetCookieArrayCmd struct {
	// View is the webview to set the cookie.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Cookie is the cookie to set.
	Cookie []webview.CookieData
}

// RemoveCookieCmd sets given cookie in the webview.
type RemoveCookieCmd struct {
	// View is the webview to remove the cookie.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Cookie is the cookie to remove.
	Cookie webview.CookieData
}

// ListCookieCmd lists all cookies in the webview.
// The response in sent via CookiesEvent using the
// provided Tag.
type ListCookieCmd struct {
	// View is the webview to list the cookies.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Buffer is the buffer to use for the response,
	// that may prevent allocations.
	Buffer []webview.CookieData
}

// StorageType is the type of storage.
type StorageType int

const (
	// StorageTypeLocal is the local storage.
	StorageTypeLocal StorageType = iota
	// StorageTypeSession is the session storage.
	StorageTypeSession
)

// SetStorageCmd sets given Storage in the webview.
type SetStorageCmd struct {
	// View is the webview to set the storage.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Local is the type of storage.
	Local StorageType

	// Content is the data to set.
	Content webview.StorageData
}

// RemoveStorageCmd sets given Storage in the webview.
type RemoveStorageCmd struct {
	// View is the webview to remove the storage.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Local is the type of storage.
	Local StorageType

	// Content is the data to remove, it
	// may match the key or the key and value.
	Content webview.StorageData
}

// ListStorageCmd lists all Storage in the webview.
//
// The response in sent via StorageEvent using the
// provided Tag.
type ListStorageCmd struct {
	// View is the webview to list the storage.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag

	// Local is the type of storage.
	Local StorageType

	// Buffer is the buffer to use for the response,
	// that may prevent allocations.
	Buffer []webview.StorageData
}

// ClearCacheCmd clears the cache of the webview.
//
// The response in sent via ErrorEvent using the
// webview as tag. Also, one
type ClearCacheCmd struct {
	// View is the webview to clear the cache.
	View event.Tag

	// Tag is the tag to send the response.
	Tag event.Tag
}

// ExecuteJavascriptCmd executes given JavaScript in the webview.
type ExecuteJavascriptCmd struct {
	// View is the webview to execute the script.
	View event.Tag

	// Script is the Javascript to execute.
	Script string
}

// InstallJavascriptCmd installs given JavaScript in the webview, executing
// it every time the webview loads a new page. The script is executed before
// the page is fully loaded.
type InstallJavascriptCmd struct {
	// View is the webview to install the script.
	View event.Tag

	// Script is the Javascript to install,
	// which will be executed every time the
	// webview loads a new page.
	Script string
}

// MessageReceiverCmd receives a message from the webview,
// and sends it to the provided Tag. The message is sent
// as a string.
//
// You can use this to communicate with the webview, by using:
//
//	window.callback.<name>(<message>);
//
// Consider that <name> is the provided Name of the callback,
// and <message> is the message to send to Tag. The Tag will
// receive the message as a string, with MessageEvent.
//
// For further information, see webview.JavascriptManager.
type MessageReceiverCmd struct {
	// View is the webview to receive the message.
	View event.Tag

	// Tag is the tag to send the message.
	Tag event.Tag

	// Name is the name of the callback.
	Name string
}

func (o NavigateCmd) ImplementsCommand()          {}
func (o DestroyCmd) ImplementsCommand()           {}
func (o SetCookieCmd) ImplementsCommand()         {}
func (o SetCookieArrayCmd) ImplementsCommand()    {}
func (o RemoveCookieCmd) ImplementsCommand()      {}
func (o ListCookieCmd) ImplementsCommand()        {}
func (o SetStorageCmd) ImplementsCommand()        {}
func (o RemoveStorageCmd) ImplementsCommand()     {}
func (o ListStorageCmd) ImplementsCommand()       {}
func (o ClearCacheCmd) ImplementsCommand()        {}
func (o ExecuteJavascriptCmd) ImplementsCommand() {}
func (o InstallJavascriptCmd) ImplementsCommand() {}
func (o MessageReceiverCmd) ImplementsCommand()   {}

func (o NavigateCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	u, err := url.Parse(o.URL)
	if err != nil {
		return
	}
	current.Navigate(u)
}

func (o SetCookieCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		err := manager.AddCookie(o.Cookie)
		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, SetCookieEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o SetCookieArrayCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		var err error
		for _, c := range o.Cookie {
			if err = manager.AddCookie(c); err != nil {
				break
			}
		}

		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, SetCookieArrayEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o RemoveCookieCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		err := manager.RemoveCookie(o.Cookie)
		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, RemoveCookieEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o ListCookieCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		evt := CookiesEvent{
			Cookies: o.Buffer,
		}
		err := manager.Cookies(func(c *webview.CookieData) bool {
			evt.Cookies = append(evt.Cookies, *c)
			return true
		})
		p.plugin.SendEvent(o.Tag, evt)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o SetStorageCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		var err error
		switch o.Local {
		case StorageTypeLocal:
			err = manager.AddLocalStorage(o.Content)
		case StorageTypeSession:
			err = manager.AddSessionStorage(o.Content)
		}

		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, SetStorageEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o RemoveStorageCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		var err error
		switch o.Local {
		case StorageTypeLocal:
			err = manager.RemoveLocalStorage(o.Content)
		case StorageTypeSession:
			err = manager.AddSessionStorage(o.Content)
		}

		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, RemoveStorageEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o ListStorageCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		evt := StorageEvent{
			Storage: o.Buffer,
		}

		fn := manager.LocalStorage
		if o.Local == StorageTypeSession {
			fn = manager.SessionStorage
		}

		err := fn(func(c *webview.StorageData) bool {
			evt.Storage = append(evt.Storage, *c)
			return true
		})

		p.plugin.SendEvent(o.Tag, evt)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o ClearCacheCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.DataManager()
	wvTag := o.View

	p.run(func() {
		err := manager.ClearAll()
		if o.Tag != nil {
			p.plugin.SendEvent(o.Tag, ClearCacheEvent{Error: err})
		}
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o ExecuteJavascriptCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.JavascriptManager()
	wvTag := o.View

	p.run(func() {
		err := manager.RunJavaScript(o.Script)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o InstallJavascriptCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.JavascriptManager()
	wvTag := o.View

	p.run(func() {
		err := manager.InstallJavascript(o.Script, webview.JavascriptOnLoadStart)
		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}

func (o MessageReceiverCmd) execute(_ *app.Window, p *webViewPlugin, _ app.FrameEvent) {
	current, ok := p.getWebView(o.View)
	if !ok || current == nil {
		return
	}

	manager := current.JavascriptManager()
	wvTag := o.View
	tag := o.Tag

	p.run(func() {
		err := manager.AddCallback(o.Name, func(msg string) {
			p.plugin.SendEvent(tag, MessageEvent{Message: msg})
		})

		if err != nil {
			p.plugin.SendEvent(wvTag, ErrorEvent{error: err})
		}
	})
}
