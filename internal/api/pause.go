package api

import (
	"encoding/json"
	"net/http"
)

type pauseResponse struct {
	JobName string `json:"job_name"`
	Paused  bool   `json:"paused"`
	Message string `json:"message"`
}

func (h *Handler) PauseJob(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}

	if err := h.store.SetPaused(name, true); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pauseResponse{
		JobName: name,
		Paused:  true,
		Message: "job paused",
	})
}

func (h *Handler) ResumeJob(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}

	if err := h.store.SetPaused(name, false); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pauseResponse{
		JobName: name,
		Paused:  false,
		Message: "job resumed",
	})
}
