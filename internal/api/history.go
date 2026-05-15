package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatch/internal/job"
)

// RunHistoryResponse represents a paginated list of job runs.
type RunHistoryResponse struct {
	JobName string    `json:"job_name"`
	Runs    []job.Run `json:"runs"`
	Total   int       `json:"total"`
}

// GetJobHistory returns the run history for a specific job.
func (h *Handler) GetJobHistory(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}

	runs, err := h.store.GetHistory(name)
	if err != nil {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	resp := RunHistoryResponse{
		JobName: name,
		Runs:    runs,
		Total:   len(runs),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
