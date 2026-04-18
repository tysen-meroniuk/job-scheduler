package web

import (
	"embed"
	"io/fs"
	"net/http"
)

// staticFS embeds the built React dashboard. The Makefile / Dockerfile copy
// dashboard/dist/ into ./static before `go build`.
//
//go:embed all:static
var staticFS embed.FS

// Handler returns an http.Handler that serves the embedded dashboard.
func Handler() http.Handler {
	sub, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}
	return http.FileServer(http.FS(sub))
}
