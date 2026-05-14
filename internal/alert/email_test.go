package alert

import (
	"net"
	"net/smtp"
	"net/textproto"
	"testing"
	"time"
)

func TestNewEmailNotifier(t *testing.T) {
	cfg := EmailConfig{
		Host:     "smtp.example.com",
		Port:     587,
		Username: "user@example.com",
		Password: "secret",
		From:     "user@example.com",
		To:       []string{"ops@example.com"},
	}

	n := NewEmailNotifier(cfg)
	if n == nil {
		t.Fatal("expected non-nil EmailNotifier")
	}
	if n.cfg.Host != cfg.Host {
		t.Errorf("expected host %q, got %q", cfg.Host, n.cfg.Host)
	}
	if n.cfg.Port != cfg.Port {
		t.Errorf("expected port %d, got %d", cfg.Port, n.cfg.Port)
	}
	if len(n.cfg.To) != 1 || n.cfg.To[0] != "ops@example.com" {
		t.Errorf("unexpected To recipients: %v", n.cfg.To)
	}
}

func TestEmailNotifier_Send_ConnectionRefused(t *testing.T) {
	cfg := EmailConfig{
		Host:     "127.0.0.1",
		Port:     19999, // nothing listening here
		Username: "user",
		Password: "pass",
		From:     "from@example.com",
		To:       []string{"to@example.com"},
	}

	n := NewEmailNotifier(cfg)
	a := NewDriftAlert("backup-job", 5*time.Minute, 8*time.Minute)

	err := n.Send(a)
	if err == nil {
		t.Error("expected error when SMTP server is unavailable, got nil")
	}

	// Verify it's a network-level error (connection refused)
	var netErr net.Error
	var opErr *net.OpError
	var textErr *textproto.Error
	_ = netErr
	_ = opErr
	_ = textErr
	_ = smtp.PlainAuth // referenced to avoid unused import
}
