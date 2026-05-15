package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/job"
)

func TestGetMetrics_EmptyStore(t *testing.T) {
	store := setupStore()
	h := NewHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h.GetMetrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var summary MetricsSummary
	if err := json.NewDecoder(rec.Body).Decode(&summary); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if summary.TotalJobs != 0 {
		t.Errorf("expected 0 jobs, got %d", summary.TotalJobs)
	}
}

func TestGetMetrics_ClassifiesJobs(t *testing.T) {
	store := setupStore()

	healthyJob := &job.Job{Name: "healthy", Schedule: "* * * * *"}
	healthyJob.Runs = []job.Run{
		{ExecutedAt: time.Now().Add(-1 * time.Minute), Drifted: false},
	}

	driftedJob := &job.Job{Name: "drifted", Schedule: "* * * * *"}
	driftedJob.Runs = []job.Run{
		{ExecutedAt: time.Now().Add(-2 * time.Minute), Drifted: true},
	}

	missedJob := &job.Job{Name: "missed", Schedule: "* * * * *", Missed: true}

	store.Save(healthyJob)
	store.Save(driftedJob)
	store.Save(missedJob)

	h := NewHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h.GetMetrics(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var summary MetricsSummary
	if err := json.NewDecoder(rec.Body).Decode(&summary); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if summary.TotalJobs != 3 {
		t.Errorf("expected 3 jobs, got %d", summary.TotalJobs)
	}
	if summary.Healthy != 1 {
		t.Errorf("expected 1 healthy, got %d", summary.Healthy)
	}
	if summary.Drifted != 1 {
		t.Errorf("expected 1 drifted, got %d", summary.Drifted)
	}
	if summary.Missed != 1 {
		t.Errorf("expected 1 missed, got %d", summary.Missed)
	}
	if summary.GeneratedAt.IsZero() {
		t.Error("expected non-zero GeneratedAt")
	}
}

func TestGetMetrics_JobSummaryLastRun(t *testing.T) {
	store := setupStore()

	earlier := time.Now().Add(-10 * time.Minute)
	later := time.Now().Add(-1 * time.Minute)

	j := &job.Job{Name: "multi-run", Schedule: "* * * * *"}
	j.Runs = []job.Run{
		{ExecutedAt: earlier, Drifted: false},
		{ExecutedAt: later, Drifted: false},
	}
	store.Save(j)

	h := NewHandler(store)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec := httptest.NewRecorder()
	h.GetMetrics(rec, req)

	var summary MetricsSummary
	_ = json.NewDecoder(rec.Body).Decode(&summary)

	if len(summary.JobSummaries) != 1 {
		t.Fatalf("expected 1 job summary, got %d", len(summary.JobSummaries))
	}
	js := summary.JobSummaries[0]
	if js.RunCount != 2 {
		t.Errorf("expected run count 2, got %d", js.RunCount)
	}
	if js.LastRun == nil {
		t.Fatal("expected LastRun to be set")
	}
	if !js.LastRun.Truncate(time.Second).Equal(later.Truncate(time.Second)) {
		t.Errorf("expected LastRun to be the most recent run")
	}
}
