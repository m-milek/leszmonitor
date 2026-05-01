package api

import (
	"embed"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func newSPAHandler(staticFiles embed.FS) http.Handler {
	staticRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		panic("Failed to read embedded static files: " + err.Error())
	}

	fileServer := http.FileServer(http.FS(staticRoot))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.NotFound(w, r)
			return
		}

		requestPath := path.Clean("/" + r.URL.Path)
		relativePath := strings.TrimPrefix(requestPath, "/")

		// Set cache headers based on file type
		if strings.HasPrefix(relativePath, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else if relativePath != "" && relativePath != "." {
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		}

		// Try to serve the file directly
		if relativePath == "" || relativePath == "." || relativePath == "index.html" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			fileServer.ServeHTTP(w, r)
			return
		}

		// For actual file requests, check if file exists before falling back to index
		if file, err := staticRoot.Open(relativePath); err == nil {
			file.Close()
			// File exists, serve it normally
			fileServer.ServeHTTP(w, r)
			return
		}

		// File doesn't exist - fall back to index.html for SPA routing
		// (but don't intercept explicit asset 404s)
		if !strings.Contains(relativePath, ".") {
			// No file extension, likely a route
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
			http.ServeFileFS(w, r, staticRoot, "index.html")
			return
		}

		// File with extension not found, let it 404
		http.NotFound(w, r)
	})
}
