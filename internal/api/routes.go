package api

import "net/http"

// RegisterRoutes attaches all API routes to the given mux.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /healthz", Healthz)
	mux.HandleFunc("GET /api/v1/jobs", h.ListJobs)
	mux.HandleFunc("GET /api/v1/jobs/{name}", h.GetJob)
	mux.HandleFunc("GET /api/v1/jobs/{name}/history", h.GetJobHistory)
	mux.HandleFunc("GET /api/v1/metrics", h.GetMetrics)
}

// Healthz is a simple liveness probe endpoint.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
