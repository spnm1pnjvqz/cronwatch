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

func TestGetJobHistory_Found(t *testing.T) {
	store := setupStore()
	store.Add(&job.Job{Name: "backup", Schedule: "@daily"})
	store.RecordRun("backup", job.Run{StartedAt: time.Now(), Duration: time.Second, Status: "ok"})
	store.RecordRun("backup", job.Run{StartedAt: time.Now(), Duration: 2 * time.Second, Status: "ok"})

	h := api.NewHandler(store)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/backup/history", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var resp api.RunHistoryResponse
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.JobName != "backup" {
		t.Errorf("expected job_name backup, got %s", resp.JobName)
	}
	if resp.Total != 2 {
		t.Errorf("expected 2 runs, got %d", resp.Total)
	}
}

func TestGetJobHistory_NotFound(t *testing.T) {
	store := setupStore()
	h := api.NewHandler(store)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/ghost/history", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", rw.Code)
	}
}

func TestGetJobHistory_EmptyHistory(t *testing.T) {
	store := setupStore()
	store.Add(&job.Job{Name: "idle", Schedule: "@hourly"})

	h := api.NewHandler(store)
	mux := http.NewServeMux()
	api.RegisterRoutes(mux, h)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/jobs/idle/history", nil)
	rw := httptest.NewRecorder()
	mux.ServeHTTP(rw, req)

	if rw.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rw.Code)
	}

	var resp api.RunHistoryResponse
	if err := json.NewDecoder(rw.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if resp.Total != 0 {
		t.Errorf("expected 0 runs, got %d", resp.Total)
	}
}
