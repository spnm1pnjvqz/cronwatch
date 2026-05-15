package job

import (
	"fmt"
	"time"
)

// Job represents a monitored cron job.
type Job struct {
	Name           string
	Schedule       string
	Interval       time.Duration
	DriftThreshold time.Duration
	LastRun        time.Time
	LastDuration   time.Duration
	AvgDuration    time.Duration
}

// Run records the result of a single job execution.
type Run struct {
	Job      *Job
	Start    time.Time
	Duration time.Duration
	Drifted  bool
	Drift    time.Duration
}

// NewRun creates a Run for the given job and start/end times.
// It calculates whether execution time drifted beyond the job's threshold.
func NewRun(j *Job, start, end time.Time) (*Run, error) {
	if j == nil {
		return nil, fmt.Errorf("job must not be nil")
	}
	if end.Before(start) {
		return nil, fmt.Errorf("end time %v is before start time %v", end, start)
	}
	duration := end.Sub(start)
	r := &Run{
		Job:      j,
		Start:    start,
		Duration: duration,
	}
	if j.DriftThreshold > 0 && j.AvgDuration > 0 {
		diff := duration - j.AvgDuration
		if diff < 0 {
			diff = -diff
		}
		if diff > j.DriftThreshold {
			r.Drifted = true
			r.Drift = diff
		}
	}
	return r, nil
}

// MissedRun describes a job that did not execute within its expected interval.
type MissedRun struct {
	Job         *Job
	ExpectedAt  time.Time
	DetectedAt  time.Time
}

// NewMissedRun creates a MissedRun detected at detectedAt.
func NewMissedRun(j *Job, detectedAt time.Time) *MissedRun {
	return &MissedRun{
		Job:        j,
		ExpectedAt: j.LastRun.Add(j.Interval),
		DetectedAt: detectedAt,
	}
}

// IsMissed reports whether the job has missed its scheduled run by now.
func (j *Job) IsMissed(now time.Time) bool {
	if j.Interval <= 0 {
		return false
	}
	return now.After(j.LastRun.Add(j.Interval))
}

// HasDrift reports whether the last recorded duration exceeded the threshold.
func (j *Job) HasDrift() bool {
	if j.DriftThreshold <= 0 || j.AvgDuration <= 0 {
		return false
	}
	diff := j.LastDuration - j.AvgDuration
	if diff < 0 {
		diff = -diff
	}
	return diff > j.DriftThreshold
}
