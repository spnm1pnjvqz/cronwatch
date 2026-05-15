package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithTags(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()
	s.Add(&job.Job{Name: "job-prod-1", Schedule: "* * * * *", Tags: job.NewTagSet([]string{"prod", "critical"})})
	s.Add(&job.Job{Name: "job-prod-2", Schedule: "*/5 * * * *", Tags: job.NewTagSet([]string{"prod"})})
	s.Add(&job.Job{Name: "job-staging", Schedule: "0 * * * *", Tags: job.NewTagSet([]string{"staging"})})
	return s
}

func TestGetJobsByTag_ReturnsMatchingJobs(t *testing.T) {
	s := setupStoreWithTags(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs?tag=prod", nil)
	rr := httptest.NewRecorder()
	h.GetJobsByTag(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(result))
	}
}

func TestGetJobsByTag_NoMatches(t *testing.T) {
	s := setupStoreWithTags(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs?tag=unknown", nil)
	rr := httptest.NewRecorder()
	h.GetJobsByTag(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(result))
	}
}

func TestGetJobsByTag_MissingParam(t *testing.T) {
	s := setupStoreWithTags(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	rr := httptest.NewRecorder()
	h.GetJobsByTag(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
