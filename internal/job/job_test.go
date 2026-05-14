package job_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/job"
)

func baseJob() *job.Job {
	return &job.Job{
		ID:               "job-1",
		Name:             "nightly-backup",
		Schedule:         "0 2 * * *",
		ExpectedDuration: 10 * time.Minute,
		DriftThreshold:   2 * time.Minute,
		GracePeriod:      5 * time.Minute,
		CreatedAt:        time.Now(),
	}
}

func TestNewRun_Success(t *testing.T) {
	j := baseJob()
	start := time.Now()
	end := start.Add(10 * time.Minute)

	r := job.NewRun(j, start, end)

	if r.Status != job.StatusSuccess {
		t.Errorf("expected success, got %s", r.Status)
	}
	if r.Duration != 10*time.Minute {
		t.Errorf("expected 10m duration, got %s", r.Duration)
	}
}

func TestNewRun_Drift(t *testing.T) {
	j := baseJob()
	start := time.Now()
	// 13 minutes — 3 minutes over expected, exceeds 2m threshold
	end := start.Add(13 * time.Minute)

	r := job.NewRun(j, start, end)

	if r.Status != job.StatusDrift {
		t.Errorf("expected drift, got %s", r.Status)
	}
	if r.Message == "" {
		t.Error("expected drift message to be set")
	}
}

func TestNewRun_WithinDriftThreshold(t *testing.T) {
	j := baseJob()
	start := time.Now()
	// 11 minutes — 1 minute over, within 2m threshold
	end := start.Add(11 * time.Minute)

	r := job.NewRun(j, start, end)

	if r.Status != job.StatusSuccess {
		t.Errorf("expected success within threshold, got %s", r.Status)
	}
}

func TestNewRun_NoDriftThreshold(t *testing.T) {
	j := baseJob()
	j.DriftThreshold = 0
	start := time.Now()
	end := start.Add(60 * time.Minute)

	r := job.NewRun(j, start, end)

	if r.Status != job.StatusSuccess {
		t.Errorf("expected success when no threshold set, got %s", r.Status)
	}
}
