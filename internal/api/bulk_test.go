package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithBulkJobs(t *testing.T, names ...string) *job.Store {
	t.Helper()
	s := job.NewStore()
	for _, name := range names {
		j := job.Job{
			ID:       name,
			Name:     name,
			Schedule: "@hourly",
		}
		if err := s.Add(j); err != nil {
			t.Fatalf("failed to add job %s: %v", name, err)
		}
	}
	return s
}

func TestBulkPauseJobs_Success(t *testing.T) {
	s := setupStoreWithBulkJobs(t, "job-a", "job-b", "job-c")
	h := NewHandler(s, nil)

	body, _ := json.Marshal(map[string][]string{"ids": {"job-a", "job-b"}})
	req := httptest.NewRequest(http.MethodPost, "/jobs/bulk/pause", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.BulkPauseJobs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["paused"] == nil {
		t.Error("expected 'paused' field in response")
	}
}

func TestBulkResumeJobs_Success(t *testing.T) {
	s := setupStoreWithBulkJobs(t, "job-x", "job-y")
	h := NewHandler(s, nil)

	// Pause first
	_ = s.Pause("job-x")
	_ = s.Pause("job-y")

	body, _ := json.Marshal(map[string][]string{"ids": {"job-x", "job-y"}})
	req := httptest.NewRequest(http.MethodPost, "/jobs/bulk/resume", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.BulkResumeJobs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if resp["resumed"] == nil {
		t.Error("expected 'resumed' field in response")
	}
}

func TestBulkPauseJobs_EmptyBody(t *testing.T) {
	s := setupStoreWithBulkJobs(t, "job-1")
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/bulk/pause", bytes.NewReader([]byte(`{}`)))
	rec := httptest.NewRecorder()

	h.BulkPauseJobs(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty ids, got %d", rec.Code)
	}
}

func TestBulkPauseJobs_PartialNotFound(t *testing.T) {
	s := setupStoreWithBulkJobs(t, "real-job")
	h := NewHandler(s, nil)

	body, _ := json.Marshal(map[string][]string{"ids": {"real-job", "ghost-job"}})
	req := httptest.NewRequest(http.MethodPost, "/jobs/bulk/pause", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	h.BulkPauseJobs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 for partial match, got %d", rec.Code)
	}

	var resp map[string]interface{}
	_ = json.NewDecoder(rec.Body).Decode(&resp)

	if resp["errors"] == nil {
		t.Error("expected 'errors' field for not-found jobs")
	}
}
