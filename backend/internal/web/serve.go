package web

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"
)

//go:embed all:dist
var frontendFS embed.FS

func NewHandler() http.Handler {
	distFS, err := fs.Sub(frontendFS, "dist")
	if err != nil {
		log.Fatalf("Failed to create sub-filesystem: %v", err)
	}

	h, err := newHandlerWithFS(distFS)
	if err != nil {
		log.Fatalf("Failed to prepare handler: %v", err)
	}
	return h
}

// newHandlerWithFS builds the SPA handler from an arbitrary filesystem so it
// can be unit-tested without the embedded dist. It serves the config-injected
// index.html for the root and SPA fallback routes, and serves real static
// files (assets, favicon, etc.) directly.
func newHandlerWithFS(distFS fs.FS) (http.Handler, error) {
	indexHTML, err := buildIndexHTML(distFS)
	if err != nil {
		return nil, err
	}

	fileServer := http.FileServer(http.FS(distFS))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cleanPath := strings.TrimPrefix(r.URL.Path, "/")

		// Serve the config-injected index.html for the root path, a direct
		// /index.html request, and any path that does not map to a real
		// static file (SPA fallback).
		if cleanPath == "" || cleanPath == "index.html" {
			serveIndex(w, indexHTML)
			return
		}
		if info, err := fs.Stat(distFS, cleanPath); err != nil || info.IsDir() {
			serveIndex(w, indexHTML)
			return
		}

		fileServer.ServeHTTP(w, r)
	}), nil
}

func serveIndex(w http.ResponseWriter, indexHTML []byte) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(indexHTML)
}

// buildIndexHTML reads the embedded index.html and injects runtime
// configuration (read from the environment) as a global JS object so the
// SPA can access values that are not known at build time. The embedded SPA
// is a static bundle, so env vars cannot be read at runtime by the client
// unless the server injects them here.
func buildIndexHTML(distFS fs.FS) ([]byte, error) {
	indexBytes, err := fs.ReadFile(distFS, "index.html")
	if err != nil {
		return nil, err
	}

	config := map[string]string{
		"logoLinkUrl": os.Getenv("LOGO_LINK_URL"),
	}
	// json.Marshal escapes <, >, & by default, making the output safe to
	// embed directly inside an HTML <script> element.
	configJSON, err := json.Marshal(config)
	if err != nil {
		return nil, fmt.Errorf("marshal app config: %w", err)
	}

	script := []byte("<script>window.__APP_CONFIG__=" + string(configJSON) + "</script>")
	if bytes.Count(indexBytes, []byte("</head>")) == 0 {
		return nil, fmt.Errorf("index.html missing </head> tag, cannot inject config")
	}
	return bytes.Replace(indexBytes, []byte("</head>"), append(script, []byte("</head>")...), 1), nil
}
