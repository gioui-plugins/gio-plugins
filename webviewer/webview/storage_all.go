//go:build ios || darwin || windows || js || android

package webview

import (
	"strconv"

	"github.com/gioui-plugins/gio-plugins/webviewer/webview/internal"
)

type dataManager struct {
	CookieManager
	StorageManager
}

func newDataManager(w *webview) DataManager {
	return &dataManager{
		CookieManager:  newCookieManager(w),
		StorageManager: newStorageManager(w),
	}
}

type storageManager struct {
	WebView
}

func newStorageManager(webview WebView) *storageManager {
	r := &storageManager{WebView: webview}
	r.JavascriptManager().AddCallback("_sendStorage", _getStorage)

	return r
}

func _getStorage(msg string) {
	if len(msg) < 48 {
		return
	}

	fr, err := strconv.ParseUint(msg[:16], 16, 64)
	if err != nil || !internal.Handle(fr).IsValid() {
		return
	}
	ksize, err := strconv.ParseUint(msg[16:32], 16, 64)
	if err != nil {
		return
	}

	vsize, err := strconv.ParseUint(msg[32:48], 16, 64)
	if err != nil {
		return
	}

	if len(msg) < int(48+ksize+vsize) {
		return
	}

	internal.Handle(fr).Value().(DataLooper[StorageData])(&StorageData{
		Key:   msg[48 : 48+ksize],
		Value: msg[48+ksize : 48+ksize+vsize],
	})
}

// LocalStorage implements the StorageManager interface.
func (s *storageManager) LocalStorage(fn DataLooper[StorageData]) (err error) {
	fr := internal.NewHandle(fn)
	defer fr.Delete()

	return s.JavascriptManager().RunJavaScript(`
	(function() {
		let fr = "` + strconv.FormatUint(uint64(fr), 16) + `".padStart(16, '0');
		for (i = 0; i < globalThis.localStorage.length; i++) {
			let key = globalThis.localStorage.key(i);
			let value = globalThis.localStorage.getItem(key);

			let [ksize, vsize] = [(key.length).toString(16).padStart(16, '0'), (value.length).toString(16).padStart(16, '0')];
			globalThis.callback._sendStorage(fr + ksize + vsize + key + value);
		}
	})();
	`)
}

// AddLocalStorage implements the StorageManager interface.
func (s *storageManager) AddLocalStorage(c StorageData) error {
	return s.JavascriptManager().RunJavaScript(`
	(function() {
		globalThis.localStorage.setItem("` + c.Key + `", "` + c.Value + `");
	})();
	`)
}

// RemoveLocalStorage implements the StorageManager interface.
func (s *storageManager) RemoveLocalStorage(c StorageData) error {
	return s.JavascriptManager().RunJavaScript(`
	(function() {
		globalThis.localStorage.removeItem("` + c.Key + `");
	})();
	`)
}

// SessionStorage implements the StorageManager interface.
func (s *storageManager) SessionStorage(fn DataLooper[StorageData]) (err error) {
	fr := internal.NewHandle(fn)
	defer fr.Delete()

	return s.JavascriptManager().RunJavaScript(`
	(function() {
		let fr = "` + strconv.FormatUint(uint64(fr), 16) + `".padStart(16, '0');
		for (i = 0; i < globalThis.sessionStorage.length; i++) {
			let key = globalThis.sessionStorage.key(i);
			let value = globalThis.sessionStorage.getItem(key);

			let [ksize, vsize] = [(key.length).toString(16).padStart(16, '0'), (value.length).toString(16).padStart(16, '0')];
			globalThis.callback._sendStorage(fr + ksize + vsize + key + value);
		}
	})();
	`)
}

// AddSessionStorage implements the StorageManager interface.
func (s *storageManager) AddSessionStorage(c StorageData) error {
	return s.JavascriptManager().RunJavaScript(`
 	(function() {
		globalThis.sessionStorage.setItem("` + c.Key + `", "` + c.Value + `");
	})();
	`)
}

// RemoveSessionStorage implements the StorageManager interface.
func (s *storageManager) RemoveSessionStorage(c StorageData) error {
	return s.JavascriptManager().RunJavaScript(`
	(function() {
		globalThis.sessionStorage.removeItem("` + c.Key + `");
	})();
	`)
}
