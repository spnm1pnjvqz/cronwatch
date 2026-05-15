package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/cronwatch/internal/job"
)

// GetJobsByLabel returns all jobs that have a matching label key=value pair.
// Query params: key, value
func (h *Handler) GetJobsByLabel(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	value := r.URL.Query().Get("value")

	if key == "" || value == "" {
		http.Error(w, "missing required query params: key and value", http.StatusBadRequest)
		return
	}

	all := h.store.List()
	matched := job.FilterByLabels(all, map[string]string{key: value})

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(matched); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
