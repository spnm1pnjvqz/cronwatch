package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVersionHandler_ReturnsOK(t *testing.T) {
	Version = "1.2.3"
	Commit = "abc1234"
	BuildTime = "2024-01-01T00:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler()(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}
}

func TestVersionHandler_ReturnsCorrectFields(t *testing.T) {
	Version = "0.9.0"
	Commit = "deadbeef"
	BuildTime = "2024-06-15T12:00:00Z"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler()(rr, req)

	var info BuildInfo
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Version != "0.9.0" {
		t.Errorf("expected version 0.9.0, got %s", info.Version)
	}
	if info.Commit != "deadbeef" {
		t.Errorf("expected commit deadbeef, got %s", info.Commit)
	}
	if info.BuildTime != "2024-06-15T12:00:00Z" {
		t.Errorf("expected build_time 2024-06-15T12:00:00Z, got %s", info.BuildTime)
	}
	if info.GoVersion == "" {
		t.Error("expected non-empty go_version")
	}
}

func TestVersionHandler_DefaultValues(t *testing.T) {
	Version = "dev"
	Commit = "none"
	BuildTime = "unknown"

	req := httptest.NewRequest(http.MethodGet, "/version", nil)
	rr := httptest.NewRecorder()

	VersionHandler()(rr, req)

	var info BuildInfo
	if err := json.NewDecoder(rr.Body).Decode(&info); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if info.Version != "dev" {
		t.Errorf("expected version dev, got %s", info.Version)
	}
}
