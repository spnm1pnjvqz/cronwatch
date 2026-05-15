package api

import (
	"net/http"
)

// RegisterRoutes attaches all HTTP routes to the given mux.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/healthz", Healthz)
	mux.HandleFunc("/jobs", h.ListJobs)
	mux.HandleFunc("/jobs/", h.GetJob)
	mux.HandleFunc("/metrics", h.GetMetrics)
}

// Healthz is a simple liveness probe endpoint.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
