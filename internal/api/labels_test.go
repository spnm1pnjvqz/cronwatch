package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithLabels(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()

	j1 := job.Job{
		ID:       "job-1",
		Name:     "Backup Job",
		Schedule: "0 2 * * *",
		Labels:   job.NewLabelSet(map[string]string{"env": "prod", "team": "infra"}),
	}
	j2 := job.Job{
		ID:       "job-2",
		Name:     "Report Job",
		Schedule: "0 6 * * *",
		Labels:   job.NewLabelSet(map[string]string{"env": "prod", "team": "data"}),
	}
	j3 := job.Job{
		ID:       "job-3",
		Name:     "Dev Job",
		Schedule: "*/5 * * * *",
		Labels:   job.NewLabelSet(map[string]string{"env": "dev", "team": "infra"}),
	}

	_ = s.Add(j1)
	_ = s.Add(j2)
	_ = s.Add(j3)

	now := time.Now()
	_ = s.RecordRun("job-1", job.Run{StartedAt: now, Duration: 2 * time.Second})
	return s
}

func TestGetJobsByLabel_ReturnsMatchingJobs(t *testing.T) {
	s := setupStoreWithLabels(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/label?key=env&value=prod", nil)
	w := httptest.NewRecorder()
	h.GetJobsByLabel(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var jobs []job.Job
	if err := json.NewDecoder(w.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(jobs) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(jobs))
	}
}

func TestGetJobsByLabel_NoMatches(t *testing.T) {
	s := setupStoreWithLabels(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs/label?key=env&value=staging", nil)
	w := httptest.NewRecorder()
	h.GetJobsByLabel(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var jobs []job.Job
	if err := json.NewDecoder(w.Body).Decode(&jobs); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if len(jobs) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(jobs))
	}
}

func TestGetJobsByLabel_MissingParams(t *testing.T) {
	s := setupStoreWithLabels(t)
	h := NewHandler(s, nil)

	tests := []struct {
		name string
		url  string
	}{
		{"missing key", "/jobs/label?value=prod"},
		{"missing value", "/jobs/label?key=env"},
		{"missing both", "/jobs/label"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.url, nil)
			w := httptest.NewRecorder()
			h.GetJobsByLabel(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}
