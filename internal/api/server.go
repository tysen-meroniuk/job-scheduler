package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/tysenmeroniuk/jobqueue/internal/queue"
	"github.com/tysenmeroniuk/jobqueue/internal/web"
)

type Server struct {
	q *queue.Queue
}

func New(q *queue.Queue) *Server {
	return &Server{q: q}
}

// Router wires up the REST endpoints and the embedded dashboard.
func (s *Server) Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api", func(r chi.Router) {
		r.Post("/jobs", s.enqueueJob)
		r.Get("/jobs", s.listJobs)
		r.Get("/jobs/{id}", s.getJob)
		r.Get("/metrics", s.getMetrics)
	})

	// Embedded React dashboard (and SPA fallback) served from /.
	r.Handle("/*", web.Handler())

	return r
}
