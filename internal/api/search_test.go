package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cronwatch/internal/job"
)

func setupStoreWithSearchData(t *testing.T) *job.Store {
	t.Helper()
	s := job.NewStore()
	s.Add(&job.Job{Name: "backup-db", Schedule: "@daily",
		Tags:   job.NewTagSet([]string{"database", "backup"}),
		Labels: job.NewLabelSet(map[string]string{"env": "prod"}),
	})
	s.Add(&job.Job{Name: "send-report", Schedule: "@weekly",
		Tags:   job.NewTagSet([]string{"email", "report"}),
		Labels: job.NewLabelSet(map[string]string{"team": "analytics"}),
	})
	s.Add(&job.Job{Name: "cleanup-tmp", Schedule: "@hourly",
		Tags:   job.NewTagSet([]string{"maintenance"}),
		Labels: job.NewLabelSet(map[string]string{"env": "staging"}),
	})
	return s
}

// decodeSearchResponse decodes the JSON body from a search response recorder
// and fails the test if decoding fails.
func decodeSearchResponse(t *testing.T, rr *httptest.ResponseRecorder) map[string]interface{} {
	t.Helper()
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return resp
}

func TestSearchJobs_MatchByName(t *testing.T) {
	h := NewHandler(setupStoreWithSearchData(t), nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/search?q=backup", nil)
	rr := httptest.NewRecorder()
	h.SearchJobs(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	resp := decodeSearchResponse(t, rr)
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected 1 result, got %v", resp["count"])
	}
}

func TestSearchJobs_MatchByTag(t *testing.T) {
	h := NewHandler(setupStoreWithSearchData(t), nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/search?q=report", nil)
	rr := httptest.NewRecorder()
	h.SearchJobs(rr, req)

	resp := decodeSearchResponse(t, rr)
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected 1 result for tag match, got %v", resp["count"])
	}
}

func TestSearchJobs_MatchByLabelValue(t *testing.T) {
	h := NewHandler(setupStoreWithSearchData(t), nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/search?q=prod", nil)
	rr := httptest.NewRecorder()
	h.SearchJobs(rr, req)

	resp := decodeSearchResponse(t, rr)
	if int(resp["count"].(float64)) != 1 {
		t.Errorf("expected 1 result for label value match, got %v", resp["count"])
	}
}

func TestSearchJobs_NoMatches(t *testing.T) {
	h := NewHandler(setupStoreWithSearchData(t), nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/search?q=zzznomatch", nil)
	rr := httptest.NewRecorder()
	h.SearchJobs(rr, req)

	resp := decodeSearchResponse(t, rr)
	if int(resp["count"].(float64)) != 0 {
		t.Errorf("expected 0 results, got %v", resp["count"])
	}
}

func TestSearchJobs_MissingParam(t *testing.T) {
	h := NewHandler(setupStoreWithSearchData(t), nil)
	req := httptest.NewRequest(http.MethodGet, "/jobs/search", nil)
	rr := httptest.NewRecorder()
	h.SearchJobs(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}
