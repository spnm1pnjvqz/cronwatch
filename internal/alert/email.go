package alert

import (
	"fmt"
	"net/smtp"
	"strings"
)

// EmailConfig holds SMTP configuration for sending alert emails.
type EmailConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	To       []string
}

// EmailNotifier sends alerts via email using SMTP.
type EmailNotifier struct {
	cfg  EmailConfig
	auth smtp.Auth
}

// NewEmailNotifier creates a new EmailNotifier with the given configuration.
func NewEmailNotifier(cfg EmailConfig) *EmailNotifier {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)
	return &EmailNotifier{cfg: cfg, auth: auth}
}

// Send delivers the alert as an email message.
func (e *EmailNotifier) Send(a Alert) error {
	subject := fmt.Sprintf("[cronwatch] %s: %s", a.Type, a.JobName)
	body := fmt.Sprintf(
		"Job: %s\nType: %s\nMessage: %s\nTime: %s",
		a.JobName, a.Type, a.Message, a.OccurredAt.Format("2006-01-02 15:04:05 UTC"),
	)

	msg := strings.Join([]string{
		"From: " + e.cfg.From,
		"To: " + strings.Join(e.cfg.To, ", "),
		"Subject: " + subject,
		"",
		body,
	}, "\r\n")

	addr := fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port)
	return smtp.SendMail(addr, e.auth, e.cfg.From, e.cfg.To, []byte(msg))
}
