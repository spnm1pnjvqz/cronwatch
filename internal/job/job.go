package job

import (
	"fmt"
	"time"
)

// Job represents a monitored cron job.
type Job struct {
	Name           string        `json:"name"`
	Schedule       string        `json:"schedule"`
	DriftThreshold time.Duration `json:"drift_threshold,omitempty"`
	LastRun        *Run          `json:"last_run,omitempty"`
	Paused         bool          `json:"paused"`
}

// Run captures a single execution record.
type Run struct {
	StartedAt  time.Time     `json:"started_at"`
	FinishedAt time.Time     `json:"finished_at"`
	Duration   time.Duration `json:"duration"`
	Drift      time.Duration `json:"drift,omitempty"`
	Missed     bool          `json:"missed"`
	Error      string        `json:"error,omitempty"`
}

// NewRun creates a Run for a completed execution. Returns an error if drift
// exceeds the job's threshold (when one is set).
func NewRun(j Job, started, finished time.Time, expectedStart time.Time) (Run, error) {
	drift := started.Sub(expectedStart)
	if drift < 0 {
		drift = -drift
	}

	r := Run{
		StartedAt:  started,
		FinishedAt: finished,
		Duration:   finished.Sub(started),
		Drift:      drift,
	}

	if j.DriftThreshold > 0 && drift > j.DriftThreshold {
		return r, fmt.Errorf(
			"job %q drifted %s (threshold %s)",
			j.Name, drift, j.DriftThreshold,
		)
	}
	return r, nil
}

// NewMissedRun creates a Run record representing a missed execution.
func NewMissedRun(expectedStart time.Time) Run {
	return Run{
		StartedAt: expectedStart,
		Missed:    true,
		Error:     "run not detected",
	}
}
