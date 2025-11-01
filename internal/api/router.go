package api

import (
	"embed"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handler, frontendFS embed.FS) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Compress(5))

	r.Route("/api", func(r chi.Router) {
		r.Get("/machines", h.GetMachines)
		r.Get("/ws/ssh/{machine}", h.SSHWebSocket)
	})

	frontendDist, err := fs.Sub(frontendFS, "frontend/dist")
	if err != nil {
		panic(err)
	}

	r.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if strings.HasPrefix(req.URL.Path, "/api") {
			http.NotFound(w, req)
			return
		}

		path := strings.TrimPrefix(req.URL.Path, "/")
		if path == "" {
			path = "index.html"
		}

		file, err := frontendDist.Open(path)
		if err != nil {
			indexFile, err := frontendDist.Open("index.html")
			if err != nil {
				http.NotFound(w, req)
				return
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html")
			http.ServeContent(w, req, "index.html", time.Time{}, indexFile.(io.ReadSeeker))
			return
		}
		defer file.Close()

		stat, err := file.Stat()
		if err != nil {
			http.NotFound(w, req)
			return
		}

		if stat.IsDir() {
			indexFile, err := frontendDist.Open(path + "/index.html")
			if err != nil {
				indexFile, err = frontendDist.Open("index.html")
				if err != nil {
					http.NotFound(w, req)
					return
				}
			}
			defer indexFile.Close()

			w.Header().Set("Content-Type", "text/html")
			http.ServeContent(w, req, "index.html", time.Time{}, indexFile.(io.ReadSeeker))
			return
		}

		http.ServeContent(w, req, path, stat.ModTime(), file.(io.ReadSeeker))
	}))

	return r
}
