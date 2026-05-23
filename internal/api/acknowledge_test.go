package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithAck(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()
	s.Add(job.Job{Name: "backup", Schedule: "@daily"})
	return s
}

func TestAcknowledgeJob_Success(t *testing.T) {
	s := setupStoreWithAck(t)
	h := NewHandler(s, nil)

	body := bytes.NewBufferString(`{"reason":"known flap"}`)
	req := httptest.NewRequest(http.MethodPost, "/jobs/backup/acknowledge", body)
	req.SetPathValue("name", "backup")
	rec := httptest.NewRecorder()

	h.acknowledgeJob(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp acknowledgeResponse
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if !resp.Acknowledged {
		t.Error("expected acknowledged=true")
	}
	if resp.Reason != "known flap" {
		t.Errorf("expected reason 'known flap', got %q", resp.Reason)
	}
	if resp.AckedAt.IsZero() {
		t.Error("expected non-zero AckedAt")
	}
}

func TestAcknowledgeJob_NotFound(t *testing.T) {
	s := setupStoreWithAck(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs/ghost/acknowledge", nil)
	req.SetPathValue("name", "ghost")
	rec := httptest.NewRecorder()

	h.acknowledgeJob(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestAcknowledgeJob_MissingParam(t *testing.T) {
	s := setupStoreWithAck(t)
	h := NewHandler(s, nil)

	req := httptest.NewRequest(http.MethodPost, "/jobs//acknowledge", nil)
	rec := httptest.NewRecorder()

	h.acknowledgeJob(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

func TestUnacknowledgeJob_Success(t *testing.T) {
	s := setupStoreWithAck(t)
	j, _ := s.Get("backup")
	j.Acknowledged = true
	j.AckReason = "old reason"
	j.AckedAt = time.Now()
	s.Update(j)

	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodDelete, "/jobs/backup/acknowledge", nil)
	req.SetPathValue("name", "backup")
	rec := httptest.NewRecorder()

	h.unacknowledgeJob(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	updated, _ := s.Get("backup")
	if updated.Acknowledged {
		t.Error("expected acknowledged=false after unack")
	}
}
