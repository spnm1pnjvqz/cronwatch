package api

import "net/http"

// RegisterBulkRoutes registers bulk operation endpoints on the given mux.
func RegisterBulkRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("/jobs/bulk/pause", h.BulkPauseJobs)
	mux.HandleFunc("/jobs/bulk/resume", h.BulkResumeJobs)
}
