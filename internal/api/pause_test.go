package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

func TestPauseJob_Success(t *testing.T) {
	store := setupStore()
	store.Add(job.Job{Name: "backup", Schedule: "@daily"})
	h := NewHandler(store, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/backup/pause", nil)
	req.SetPathValue("name", "backup")
	w := httptest.NewRecorder()

	h.PauseJob(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp pauseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !resp.Paused {
		t.Error("expected paused=true")
	}
	if resp.JobName != "backup" {
		t.Errorf("expected job_name=backup, got %s", resp.JobName)
	}
}

func TestPauseJob_NotFound(t *testing.T) {
	store := setupStore()
	h := NewHandler(store, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/ghost/pause", nil)
	req.SetPathValue("name", "ghost")
	w := httptest.NewRecorder()

	h.PauseJob(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestResumeJob_Success(t *testing.T) {
	store := setupStore()
	store.Add(job.Job{Name: "backup", Schedule: "@daily"})
	store.SetPaused("backup", true)
	h := NewHandler(store, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/backup/resume", nil)
	req.SetPathValue("name", "backup")
	w := httptest.NewRecorder()

	h.ResumeJob(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp pauseResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Paused {
		t.Error("expected paused=false")
	}
}

func TestPauseJob_MissingParam(t *testing.T) {
	store := setupStore()
	h := NewHandler(store, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs//pause", nil)
	w := httptest.NewRecorder()

	h.PauseJob(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
