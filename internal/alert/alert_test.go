package alert_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cronwatch/internal/alert"
)

func TestNewDriftAlert(t *testing.T) {
	a := alert.NewDriftAlert("backup", 5*time.Minute, 8*time.Minute)

	if a.JobName != "backup" {
		t.Errorf("expected job name 'backup', got %q", a.JobName)
	}
	if a.Type != alert.AlertTypeDrift {
		t.Errorf("expected type drift, got %q", a.Type)
	}
	if a.Message == "" {
		t.Error("expected non-empty message")
	}
	if a.OccurredAt.IsZero() {
		t.Error("expected OccurredAt to be set")
	}
}

func TestNewMissedAlert(t *testing.T) {
	last := time.Now().Add(-2 * time.Hour)
	a := alert.NewMissedAlert("sync", last)

	if a.Type != alert.AlertTypeMissed {
		t.Errorf("expected type missed, got %q", a.Type)
	}
	if a.JobName != "sync" {
		t.Errorf("expected job name 'sync', got %q", a.JobName)
	}
}

func TestWebhookNotifier_Send_Success(t *testing.T) {
	var received map[string]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := alert.NewWebhookNotifier(server.URL)
	a := alert.NewDriftAlert("cleanup", 1*time.Minute, 3*time.Minute)

	if err := n.Send(a); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received["job_name"] != "cleanup" {
		t.Errorf("expected job_name 'cleanup', got %q", received["job_name"])
	}
	if received["type"] != "drift" {
		t.Errorf("expected type 'drift', got %q", received["type"])
	}
}

func TestWebhookNotifier_Send_NonOKStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := alert.NewWebhookNotifier(server.URL)
	a := alert.NewMissedAlert("report", time.Now().Add(-1*time.Hour))

	if err := n.Send(a); err == nil {
		t.Error("expected error for non-2xx response, got nil")
	}
}
