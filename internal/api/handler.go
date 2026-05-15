package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatch/internal/job"
)

// Handler holds dependencies for HTTP handlers.
type Handler struct {
	store *job.Store
}

// NewHandler creates a new Handler with the given store.
func NewHandler(store *job.Store) *Handler {
	return &Handler{store: store}
}

// ListJobs returns all tracked jobs as JSON.
func (h *Handler) ListJobs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	jobs := h.store.All()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

// GetJob returns a single job by name.
func (h *Handler) GetJob(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing name query parameter", http.StatusBadRequest)
		return
	}
	j, ok := h.store.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(j); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
