package main

import (
	"embed"
	_ "embed"
	"io/fs"
	"log"
	"net/http"
	"path"
)

// http://localhost:8080/web/index.html
func main2() {
	mux := http.DefaultServeMux
	mux.Handle("/web/", AssetHandler("/web/", Assets, "./resources"))

	log.Fatal(http.ListenAndServe(":8080", mux))
}

type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

// AssetHandler returns an http.Handler that will serve files from
// the Assets embed.FS. When locating a file, it will strip the given
// prefix from the request and prepend the root to the fileapp.
func AssetHandler(prefix string, assets embed.FS, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		// If we can't find the asset, fs can handle the error
		file, err := assets.Open(assetPath)
		if err != nil {
			return nil, err
		}

		// Otherwise assume this is a legitimate request routed correctly
		return file, err
	})

	return http.StripPrefix(prefix, http.FileServer(http.FS(handler)))
}
