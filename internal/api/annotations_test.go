package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithAnnotations() *job.Store {
	s := job.NewStore()
	j := job.Job{
		Name:        "annotated-job",
		Schedule:    "@hourly",
		Annotations: job.NewAnnotationSet(map[string]string{"owner": "alice", "env": "prod"}),
	}
	s.Add(j)
	return s
}

func TestGetJobAnnotations_Found(t *testing.T) {
	s := setupStoreWithAnnotations()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/annotated-job/annotations", nil)
	req.SetPathValue("name", "annotated-job")
	rec := httptest.NewRecorder()
	h.getJobAnnotations(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	var result map[string]string
	json.NewDecoder(rec.Body).Decode(&result)
	if result["owner"] != "alice" {
		t.Errorf("expected owner=alice, got %q", result["owner"])
	}
}

func TestGetJobAnnotations_NotFound(t *testing.T) {
	s := job.NewStore()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/ghost/annotations", nil)
	req.SetPathValue("name", "ghost")
	rec := httptest.NewRecorder()
	h.getJobAnnotations(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rec.Code)
	}
}

func TestSetJobAnnotation_Success(t *testing.T) {
	s := setupStoreWithAnnotations()
	h := NewHandler(s, nil)
	body, _ := json.Marshal(map[string]string{"key": "team", "value": "platform"})
	req := httptest.NewRequest(http.MethodPost, "/jobs/annotated-job/annotations", bytes.NewReader(body))
	req.SetPathValue("name", "annotated-job")
	rec := httptest.NewRecorder()
	h.setJobAnnotation(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	j, _ := s.Get("annotated-job")
	if v, ok := j.Annotations.Get("team"); !ok || v != "platform" {
		t.Errorf("expected team=platform after set")
	}
}

func TestSetJobAnnotation_BadBody(t *testing.T) {
	s := setupStoreWithAnnotations()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodPost, "/jobs/annotated-job/annotations", bytes.NewBufferString("not-json"))
	req.SetPathValue("name", "annotated-job")
	rec := httptest.NewRecorder()
	h.setJobAnnotation(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rec.Code)
	}
}

func TestDeleteJobAnnotation_Success(t *testing.T) {
	s := setupStoreWithAnnotations()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodDelete, "/jobs/annotated-job/annotations/env", nil)
	req.SetPathValue("name", "annotated-job")
	req.SetPathValue("key", "env")
	rec := httptest.NewRecorder()
	h.deleteJobAnnotation(rec, req)
	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rec.Code)
	}
	j, _ := s.Get("annotated-job")
	if _, ok := j.Annotations.Get("env"); ok {
		t.Error("expected env annotation to be deleted")
	}
}
