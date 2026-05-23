package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/cronwatch/internal/job"
)

// bulkPauseRequest holds a list of job names to pause or resume.
type bulkPauseRequest struct {
	Jobs []string `json:"jobs"`
}

// bulkPauseResponse reports the outcome for each job in a bulk operation.
type bulkPauseResponse struct {
	Succeeded []string          `json:"succeeded"`
	Failed    map[string]string `json:"failed"`
}

// BulkPauseJobs pauses multiple jobs in a single request.
// POST /jobs/bulk/pause
// Body: {"jobs": ["job-a", "job-b"]}
func (h *Handler) BulkPauseJobs(w http.ResponseWriter, r *http.Request) {
	var req bulkPauseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Jobs) == 0 {
		http.Error(w, "invalid or empty jobs list", http.StatusBadRequest)
		return
	}

	resp := bulkPauseResponse{
		Succeeded: []string{},
		Failed:    map[string]string{},
	}

	for _, name := range req.Jobs {
		j, err := h.store.Get(name)
		if err != nil {
			resp.Failed[name] = "not found"
			continue
		}
		j.Paused = true
		if err := h.store.Update(j); err != nil {
			resp.Failed[name] = err.Error()
			continue
		}
		resp.Succeeded = append(resp.Succeeded, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BulkResumeJobs resumes multiple paused jobs in a single request.
// POST /jobs/bulk/resume
// Body: {"jobs": ["job-a", "job-b"]}
func (h *Handler) BulkResumeJobs(w http.ResponseWriter, r *http.Request) {
	var req bulkPauseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Jobs) == 0 {
		http.Error(w, "invalid or empty jobs list", http.StatusBadRequest)
		return
	}

	resp := bulkPauseResponse{
		Succeeded: []string{},
		Failed:    map[string]string{},
	}

	for _, name := range req.Jobs {
		j, err := h.store.Get(name)
		if err != nil {
			resp.Failed[name] = "not found"
			continue
		}
		j.Paused = false
		if err := h.store.Update(j); err != nil {
			resp.Failed[name] = err.Error()
			continue
		}
		resp.Succeeded = append(resp.Succeeded, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// BulkTagJobs adds a tag to multiple jobs in a single request.
// POST /jobs/bulk/tag
// Body: {"jobs": ["job-a", "job-b"], "tag": "critical"}
func (h *Handler) BulkTagJobs(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Jobs []string `json:"jobs"`
		Tag  string   `json:"tag"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Jobs) == 0 || req.Tag == "" {
		http.Error(w, "invalid request: jobs and tag are required", http.StatusBadRequest)
		return
	}

	resp := bulkPauseResponse{
		Succeeded: []string{},
		Failed:    map[string]string{},
	}

	for _, name := range req.Jobs {
		j, err := h.store.Get(name)
		if err != nil {
			resp.Failed[name] = "not found"
			continue
		}
		existing := j.Tags.Slice()
		j.Tags = job.NewTagSet(append(existing, req.Tag)...)
		if err := h.store.Update(j); err != nil {
			resp.Failed[name] = err.Error()
			continue
		}
		resp.Succeeded = append(resp.Succeeded, name)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
