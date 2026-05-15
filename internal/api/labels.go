package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/user/cronwatch/internal/job"
)

// GetJobsByLabels handles GET /jobs/labels?env=prod&team=ops
// It returns all jobs whose labels match every query parameter provided.
func (h *Handler) GetJobsByLabels(w http.ResponseWriter, r *http.Request) {
	selector := make(map[string]string)
	for key, vals := range r.URL.Query() {
		if len(vals) > 0 {
			selector[strings.ToLower(strings.TrimSpace(key))] = strings.TrimSpace(vals[0])
		}
	}

	if len(selector) == 0 {
		http.Error(w, "at least one label selector query param required", http.StatusBadRequest)
		return
	}

	all := h.store.List()
	matched := job.FilterByLabels(all, selector)

	type jobSummary struct {
		Name   string            `json:"name"`
		Labels map[string]string `json:"labels"`
	}

	result := make([]jobSummary, 0, len(matched))
	for _, j := range matched {
		labels := make(map[string]string)
		for _, k := range j.Labels.Keys() {
			v, _ := j.Labels.Get(k)
			labels[k] = v
		}
		result = append(result, jobSummary{Name: j.Name, Labels: labels})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
