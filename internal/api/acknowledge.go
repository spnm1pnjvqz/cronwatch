package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/cronwatch/internal/job"
)

type acknowledgeRequest struct {
	Reason string `json:"reason"`
}

type acknowledgeResponse struct {
	JobName     string    `json:"job_name"`
	Acknowledged bool     `json:"acknowledged"`
	Reason      string    `json:"reason"`
	AckedAt     time.Time `json:"acked_at"`
}

func (h *Handler) acknowledgeJob(w http.ResponseWriter, r *http.Request) {
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

	var req acknowledgeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		req.Reason = ""
	}

	j.Acknowledged = true
	j.AckReason = req.Reason
	j.AckedAt = time.Now().UTC()
	h.store.Update(j)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acknowledgeResponse{
		JobName:      j.Name,
		Acknowledged: true,
		Reason:       j.AckReason,
		AckedAt:      j.AckedAt,
	})
}

func (h *Handler) unacknowledgeJob(w http.ResponseWriter, r *http.Request) {
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

	j.Acknowledged = false
	j.AckReason = ""
	j.AckedAt = time.Time{}
	h.store.Update(j)

	w.WriteHeader(http.StatusNoContent)
}

// ensure job.Job fields are referenced to avoid import issues
var _ = job.Job{}
