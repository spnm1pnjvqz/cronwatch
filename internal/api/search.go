package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cronwatch/internal/job"
)

// searchResult wraps a job with a match score/reason for search results.
type searchResult struct {
	Job    *job.Job `json:"job"`
	Reason string   `json:"match_reason"`
}

// SearchJobs handles GET /jobs/search?q=<query>
// It performs a case-insensitive substring search across job name, tags, and labels.
func (h *Handler) SearchJobs(w http.ResponseWriter, r *http.Request) {
	query := strings.TrimSpace(r.URL.Query().Get("q"))
	if query == "" {
		http.Error(w, "missing query parameter: q", http.StatusBadRequest)
		return
	}

	lower := strings.ToLower(query)
	all := h.store.List()
	results := make([]searchResult, 0)

	for _, j := range all {
		reason := matchReason(j, lower)
		if reason != "" {
			results = append(results, searchResult{Job: j, Reason: reason})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"query":   query,
		"count":   len(results),
		"results": results,
	})
}

// matchReason returns a human-readable reason if the job matches the query,
// or an empty string if it does not match.
func matchReason(j *job.Job, lower string) string {
	if strings.Contains(strings.ToLower(j.Name), lower) {
		return "name"
	}
	for _, tag := range j.Tags.Slice() {
		if strings.Contains(strings.ToLower(tag), lower) {
			return "tag:" + tag
		}
	}
	for _, k := range j.Labels.Keys() {
		if strings.Contains(strings.ToLower(k), lower) {
			return "label_key:" + k
		}
		v, _ := j.Labels.Get(k)
		if strings.Contains(strings.ToLower(v), lower) {
			return "label_value:" + k
		}
	}
	return ""
}
