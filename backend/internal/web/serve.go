package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"strings"
)

//go:embed all:dist
var frontendFS embed.FS

func NewHandler() http.Handler {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem: %v", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if _, err := fs.ReadFile(distFS, strings.TrimPrefix(path, "/")); err == nil {
			http.FileServer(http.FS(distFS)).ServeHTTP(w, r)
			return
		}

		if index, err := fs.ReadFile(distFS, "index.html"); err == nil {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write(index)
			return
		}

		http.NotFound(w, r)
	})
}
