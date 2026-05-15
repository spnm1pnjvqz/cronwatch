package api

import (
	"encoding/json"
	"net/http"

	"github.com/user/cronwatch/internal/config"
)

// alertConfigResponse is the JSON representation of alert configuration.
type alertConfigResponse struct {
	WebhookEnabled bool   `json:"webhook_enabled"`
	WebhookURL     string `json:"webhook_url,omitempty"`
	EmailEnabled   bool   `json:"email_enabled"`
	EmailFrom      string `json:"email_from,omitempty"`
	EmailTo        string `json:"email_to,omitempty"`
	SMTPHost       string `json:"smtp_host,omitempty"`
	SMTPPort       int    `json:"smtp_port,omitempty"`
}

// GetAlertConfig returns the current alert configuration.
func (h *Handler) GetAlertConfig(w http.ResponseWriter, r *http.Request) {
	cfg := h.cfg
	if cfg == nil {
		http.Error(w, "config not available", http.StatusInternalServerError)
		return
	}

	resp := alertConfigResponse{
		WebhookEnabled: cfg.Alerts.Webhook.URL != "",
		WebhookURL:     cfg.Alerts.Webhook.URL,
		EmailEnabled:   cfg.Alerts.Email.From != "",
		EmailFrom:      cfg.Alerts.Email.From,
		EmailTo:        cfg.Alerts.Email.To,
		SMTPHost:       cfg.Alerts.Email.SMTPHost,
		SMTPPort:       cfg.Alerts.Email.SMTPPort,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}

// Handler needs access to config; extend it here via a helper.
// cfg is stored on the Handler struct (added below as a patch).
var _ = (*config.Config)(nil) // ensure import is used
