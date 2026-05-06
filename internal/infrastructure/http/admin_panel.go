package http

import (
	"errors"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/axarus/vectrag/admin"
)

type AdminHandlerProvider struct{}

func (AdminHandlerProvider) Handler() http.Handler {
	return AdminHandler()
}

// AdminHandler returns an HTTP handler that serves the admin panel.
// It serves static files from the embedded "dist" directory, and falls back to index.html
func AdminHandler() http.Handler {
	// Create a sub-filesystem pointing to the "dist" directory inside the embedded admin files
	distFS, err := fs.Sub(admin.Dist, "dist")
	if err != nil {
		// If there is an error initializing the filesystem, log it and return a handler that always responds with 500
		log.Printf("error initializing the file system: %v", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "internal server error", http.StatusInternalServerError)
		})
	}

	// Create a standard file server to serve static files from distFS
	fsHandler := http.FileServer(http.FS(distFS))

	// Return the actual HTTP handler function
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove leading "/" from the URL path
		path := strings.TrimPrefix(r.URL.Path, "/")

		if path != "" {
			// Attempt to open the requested file
			f, err := distFS.Open(path)
			if err == nil {
				// If the file exists, close it immediately (we just wanted to check existence)
				// and serve it using the file server
				f.Close()
				fsHandler.ServeHTTP(w, r)
				return
			} else if !errors.Is(err, fs.ErrNotExist) {
				// If there was an error other than "file not found", return 500
				http.Error(w, "error serving file", http.StatusInternalServerError)
				return
			}
		}

		// Read the index.html file from the embedded filesystem
		data, err := fs.ReadFile(distFS, "index.html")
		if err != nil {
			// If index.html cannot be read, return 500
			http.Error(w, "index.html not found", http.StatusInternalServerError)
			return
		}

		// Serve index.html with proper content type and 200 OK status
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	})
}
