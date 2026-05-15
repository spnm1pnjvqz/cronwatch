package config

import (
	"os"
	"testing"
	"time"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "cronwatch-*.yaml")
	if err != nil {
		t.Fatalf("creating temp file: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad_ValidConfig(t *testing.T) {
	path := writeTempConfig(t, `
jobs:
  - name: backup
    schedule: "0 2 * * *"
    drift_threshold: 5m
notifiers:
  webhook:
    url: https://hooks.example.com/alert
`)

	cfg, err := Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(cfg.Jobs) != 1 {
		t.Fatalf("expected 1 job, got %d", len(cfg.Jobs))
	}
	if cfg.Jobs[0].Name != "backup" {
		t.Errorf("expected job name 'backup', got %q", cfg.Jobs[0].Name)
	}
	if cfg.Jobs[0].DriftThreshold != 5*time.Minute {
		t.Errorf("expected drift threshold 5m, got %v", cfg.Jobs[0].DriftThreshold)
	}
	if cfg.Notifiers.Webhook == nil {
		t.Fatal("expected webhook notifier to be set")
	}
	if cfg.Notifiers.Webhook.URL != "https://hooks.example.com/alert" {
		t.Errorf("unexpected webhook URL: %q", cfg.Notifiers.Webhook.URL)
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := Load("/nonexistent/path/config.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_NoJobs(t *testing.T) {
	path := writeTempConfig(t, `jobs: []\n`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for empty jobs list")
	}
}

func TestLoad_JobMissingName(t *testing.T) {
	path := writeTempConfig(t, `
jobs:
  - schedule: "* * * * *"
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing job name")
	}
}

func TestLoad_JobMissingSchedule(t *testing.T) {
	path := writeTempConfig(t, `
jobs:
  - name: myjob
`)
	_, err := Load(path)
	if err == nil {
		t.Fatal("expected validation error for missing schedule")
	}
}
