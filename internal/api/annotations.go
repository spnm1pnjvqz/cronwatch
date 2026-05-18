package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/cronwatch/internal/job"
)

// getJobAnnotations returns the annotation map for a specific job.
func (h *Handler) getJobAnnotations(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}
	j, ok := h.store.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(j.Annotations.Map())
}

// setJobAnnotation adds or updates a single annotation on a job.
func (h *Handler) setJobAnnotation(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}
	j, ok := h.store.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	var body struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Key == "" {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	j.Annotations.Set(body.Key, body.Value)
	h.store.Update(j)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
}

// deleteJobAnnotation removes an annotation key from a job.
func (h *Handler) deleteJobAnnotation(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	key := r.PathValue("key")
	if name == "" || key == "" {
		http.Error(w, "missing job name or key", http.StatusBadRequest)
		return
	}
	j, ok := h.store.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}
	j.Annotations.Delete(key)
	h.store.Update(j)
	w.WriteHeader(http.StatusNoContent)
}

// filterByAnnotation is a helper used by search/filter routes.
func filterByAnnotation(jobs []job.Job, key, value string) []job.Job {
	var out []job.Job
	for _, j := range jobs {
		if v, ok := j.Annotations.Get(key); ok && v == value {
			out = append(out, j)
		}
	}
	return out
}
