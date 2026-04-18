package api

import "net/http"

func (s *Server) enqueueJob(w http.ResponseWriter, r *http.Request) {
	// TODO: decode JSON body {queue?, type, payload, priority?, run_at?, max_attempts?, idempotency_key?}
	//       call s.q.Enqueue, return {id}
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) listJobs(w http.ResponseWriter, r *http.Request) {
	// TODO: paginated list with optional state/queue filters
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) getJob(w http.ResponseWriter, r *http.Request) {
	// TODO: fetch job by id, 404 if missing
	http.Error(w, "not implemented", http.StatusNotImplemented)
}

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	// TODO: counts by_state and by_queue for the dashboard
	http.Error(w, "not implemented", http.StatusNotImplemented)
}
