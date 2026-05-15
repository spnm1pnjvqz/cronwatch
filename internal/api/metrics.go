package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/user/cronwatch/internal/job"
)

// MetricsSummary holds aggregated statistics across all tracked jobs.
type MetricsSummary struct {
	TotalJobs    int            `json:"total_jobs"`
	Healthy      int            `json:"healthy"`
	Drifted      int            `json:"drifted"`
	Missed       int            `json:"missed"`
	GeneratedAt  time.Time      `json:"generated_at"`
	JobSummaries []JobSummary   `json:"job_summaries"`
}

// JobSummary holds per-job statistics.
type JobSummary struct {
	Name       string     `json:"name"`
	LastRun    *time.Time `json:"last_run,omitempty"`
	RunCount   int        `json:"run_count"`
	DriftCount int        `json:"drift_count"`
	Missed     bool       `json:"missed"`
}

// GetMetrics returns an aggregated summary of all job states.
func (h *Handler) GetMetrics(w http.ResponseWriter, r *http.Request) {
	jobs := h.store.List()

	summary := MetricsSummary{
		TotalJobs:   len(jobs),
		GeneratedAt: time.Now().UTC(),
	}

	for _, j := range jobs {
		js := buildJobSummary(j)
		if j.Missed {
			summary.Missed++
		} else if js.DriftCount > 0 {
			summary.Drifted++
		} else {
			summary.Healthy++
		}
		summary.JobSummaries = append(summary.JobSummaries, js)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(summary)
}

func buildJobSummary(j *job.Job) JobSummary {
	js := JobSummary{
		Name:   j.Name,
		Missed: j.Missed,
	}
	for _, run := range j.Runs {
		js.RunCount++
		if run.Drifted {
			js.DriftCount++
		}
		if js.LastRun == nil || run.ExecutedAt.After(*js.LastRun) {
			t := run.ExecutedAt
			js.LastRun = &t
		}
	}
	return js
}
