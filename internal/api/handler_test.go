package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/job"
)

func setupStore(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()
	s.Upsert(job.Job{
		Name:          "backup",
		Schedule:      "0 2 * * *",
		LastRun:       time.Now().Add(-1 * time.Hour),
		DriftThreshold: 5 * time.Minute,
	})
	return s
}

func TestListJobs_ReturnsJobs(t *testing.T) {
	s := setupStore(t)
	h := api.NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/jobs", nil)
	rec := httptest.NewRecorder()
	h.ListJobs(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var jobs []job.Job
	if err := json.NewDecoder(rec.Body).Decode(&jobs); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(jobs))
	}
}

func TestGetJob_Found(t *testing.T) {
	s := setupStore(t)
	h := api.NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/jobs/detail?name=backup", nil)
	rec := httptest.NewRecorder()
	h.GetJob(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var j job.Job
	if err := json.NewDecoder(rec.Body).Decode(&j); err != nil {
		t.Fatalf("failed to decode: %v", err)
	}
	if j.Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", j.Name)
	}
}

func TestGetJob_NotFound(t *testing.T) {
	s := setupStore(t)
	h := api.NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/jobs/detail?name=nonexistent", nil)
	rec := httptest.NewRecorder()
	h.GetJob(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestGetJob_MissingParam(t *testing.T) {
	s := setupStore(t)
	h := api.NewHandler(s)

	req := httptest.NewRequest(http.MethodGet, "/jobs/detail", nil)
	rec := httptest.NewRecorder()
	h.GetJob(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()
	api.Healthz(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
}
