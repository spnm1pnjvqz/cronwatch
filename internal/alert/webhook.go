package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// WebhookNotifier sends alerts as JSON POST requests to a configured URL.
type WebhookNotifier struct {
	URL    string
	Client *http.Client
}

// NewWebhookNotifier returns a WebhookNotifier with a sensible default timeout.
func NewWebhookNotifier(url string) *WebhookNotifier {
	return &WebhookNotifier{
		URL: url,
		Client: &http.Client{Timeout: 10 * time.Second},
	}
}

type webhookPayload struct {
	JobName    string `json:"job_name"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	OccurredAt string `json:"occurred_at"`
}

// Send marshals the alert and POSTs it to the configured webhook URL.
func (w *WebhookNotifier) Send(a Alert) error {
	payload := webhookPayload{
		JobName:    a.JobName,
		Type:       string(a.Type),
		Message:    a.Message,
		OccurredAt: a.OccurredAt.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	resp, err := w.Client.Post(w.URL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("webhook: post request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("webhook: unexpected status %d from %s", resp.StatusCode, w.URL)
	}

	return nil
}
