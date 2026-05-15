package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/cronwatch/internal/job"
)

// GetJobsByTag handles GET /jobs?tag=<value> and filters jobs by tag.
func (h *Handler) GetJobsByTag(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	if tag == "" {
		http.Error(w, "missing tag query parameter", http.StatusBadRequest)
		return
	}

	all := h.store.List()
	matched := job.FilterByTag(all, tag)

	type jobSummary struct {
		Name     string   `json:"name"`
		Schedule string   `json:"schedule"`
		Tags     []string `json:"tags"`
	}

	results := make([]jobSummary, 0, len(matched))
	for _, j := range matched {
		results = append(results, jobSummary{
			Name:     j.Name,
			Schedule: j.Schedule,
			Tags:     j.Tags.Slice(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(results)
}
