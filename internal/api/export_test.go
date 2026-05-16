package api

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/job"
)

func setupStoreWithHistory(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()
	j := job.Job{
		Name:     "backup",
		Schedule: "@daily",
	}
	s.Add(j)
	run1 := job.Run{
		RunAt:    time.Now().Add(-2 * time.Hour),
		Duration: 500 * time.Millisecond,
		Drifted:  false,
	}
	run2 := job.Run{
		RunAt:       time.Now().Add(-1 * time.Hour),
		Duration:    1200 * time.Millisecond,
		Drifted:     true,
		DriftAmount: 200 * time.Millisecond,
	}
	s.RecordRun("backup", run1)
	s.RecordRun("backup", run2)
	return s
}

func TestExportJobHistory_JSON(t *testing.T) {
	s := setupStoreWithHistory(t)
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/backup/export", nil)
	req.SetPathValue("name", "backup")
	w := httptest.NewRecorder()
	h.exportJobHistory(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var rows []map[string]string
	if err := json.NewDecoder(w.Body).Decode(&rows); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if len(rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(rows))
	}
	if rows[1]["drifted"] != "true" {
		t.Errorf("expected second row to be drifted")
	}
}

func TestExportJobHistory_CSV(t *testing.T) {
	s := setupStoreWithHistory(t)
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/backup/export?format=csv", nil)
	req.SetPathValue("name", "backup")
	w := httptest.NewRecorder()
	h.exportJobHistory(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "text/csv" {
		t.Errorf("expected text/csv, got %s", ct)
	}
	r := csv.NewReader(strings.NewReader(w.Body.String()))
	records, err := r.ReadAll()
	if err != nil {
		t.Fatalf("failed to parse CSV: %v", err)
	}
	// header + 2 data rows
	if len(records) != 3 {
		t.Errorf("expected 3 CSV records, got %d", len(records))
	}
}

func TestExportJobHistory_NotFound(t *testing.T) {
	s := job.NewStore()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/ghost/export", nil)
	req.SetPathValue("name", "ghost")
	w := httptest.NewRecorder()
	h.exportJobHistory(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestExportJobHistory_MissingParam(t *testing.T) {
	s := job.NewStore()
	h := NewHandler(s, nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs//export", nil)
	w := httptest.NewRecorder()
	h.exportJobHistory(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}
