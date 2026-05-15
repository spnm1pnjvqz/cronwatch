package api

import "net/http"

// RegisterRoutes attaches all API routes to the given mux.
func RegisterRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/jobs", h.ListJobs)
	mux.HandleFunc("/jobs/detail", h.GetJob)
	mux.HandleFunc("/healthz", Healthz)
}

// Healthz is a simple liveness probe endpoint.
func Healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
