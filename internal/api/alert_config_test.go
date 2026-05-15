package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/job"
)

func setupHandlerWithConfig(cfg *config.Config) *Handler {
	store := job.NewStore()
	h := NewHandler(store)
	h.cfg = cfg
	return h
}

func TestGetAlertConfig_WebhookAndEmail(t *testing.T) {
	cfg := &config.Config{
		Alerts: config.Alerts{
			Webhook: config.Webhook{URL: "https://hooks.example.com/notify"},
			Email: config.Email{
				From:     "cron@example.com",
				To:       "ops@example.com",
				SMTPHost: "smtp.example.com",
				SMTPPort: 587,
			},
		},
	}
	h := setupHandlerWithConfig(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts/config", nil)
	w := httptest.NewRecorder()
	h.GetAlertConfig(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp alertConfigResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if !resp.WebhookEnabled {
		t.Error("expected webhook_enabled to be true")
	}
	if resp.WebhookURL != "https://hooks.example.com/notify" {
		t.Errorf("unexpected webhook_url: %s", resp.WebhookURL)
	}
	if !resp.EmailEnabled {
		t.Error("expected email_enabled to be true")
	}
	if resp.SMTPPort != 587 {
		t.Errorf("expected smtp_port 587, got %d", resp.SMTPPort)
	}
}

func TestGetAlertConfig_NoConfig(t *testing.T) {
	store := job.NewStore()
	h := NewHandler(store)
	// h.cfg is nil by default

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts/config", nil)
	w := httptest.NewRecorder()
	h.GetAlertConfig(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestGetAlertConfig_EmptyAlerts(t *testing.T) {
	cfg := &config.Config{}
	h := setupHandlerWithConfig(cfg)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/alerts/config", nil)
	w := httptest.NewRecorder()
	h.GetAlertConfig(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp alertConfigResponse
	_ = json.NewDecoder(w.Body).Decode(&resp)

	if resp.WebhookEnabled {
		t.Error("expected webhook_enabled to be false")
	}
	if resp.EmailEnabled {
		t.Error("expected email_enabled to be false")
	}
}
