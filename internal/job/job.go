package job

import "time"

// Status represents the execution status of a cron job run.
type Status string

const (
	StatusSuccess Status = "success"
	StatusMissed  Status = "missed"
	StatusDrift   Status = "drift"
)

// Job defines a monitored cron job and its expected schedule.
type Job struct {
	ID              string        `json:"id"`
	Name            string        `json:"name"`
	Schedule        string        `json:"schedule"` // cron expression
	ExpectedDuration time.Duration `json:"expected_duration"`
	DriftThreshold  time.Duration `json:"drift_threshold"`
	GracePeriod     time.Duration `json:"grace_period"`
	CreatedAt       time.Time     `json:"created_at"`
}

// Run records a single execution of a cron job.
type Run struct {
	JobID     string        `json:"job_id"`
	StartedAt time.Time     `json:"started_at"`
	EndedAt   time.Time     `json:"ended_at"`
	Duration  time.Duration `json:"duration"`
	Status    Status        `json:"status"`
	Message   string        `json:"message,omitempty"`
}

// NewRun creates a Run and evaluates its status against the job's thresholds.
func NewRun(j *Job, startedAt, endedAt time.Time) *Run {
	duration := endedAt.Sub(startedAt)
	r := &Run{
		JobID:     j.ID,
		StartedAt: startedAt,
		EndedAt:   endedAt,
		Duration:  duration,
		Status:    StatusSuccess,
	}

	if j.ExpectedDuration > 0 && j.DriftThreshold > 0 {
		diff := duration - j.ExpectedDuration
		if diff < 0 {
			diff = -diff
		}
		if diff > j.DriftThreshold {
			r.Status = StatusDrift
			r.Message = "execution time drifted beyond threshold"
		}
	}

	return r
}

// IsMissed reports whether the run was recorded as missed.
func (r *Run) IsMissed() bool {
	return r.Status == StatusMissed
}

// NewMissedRun creates a Run with StatusMissed for a job that did not execute
// within its expected window. The expectedAt time is used as both start and end.
func NewMissedRun(j *Job, expectedAt time.Time) *Run {
	return &Run{
		JobID:     j.ID,
		StartedAt: expectedAt,
		EndedAt:   expectedAt,
		Duration:  0,
		Status:    StatusMissed,
		Message:   "job did not run within the expected window",
	}
}
