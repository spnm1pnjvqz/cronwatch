package alert

import (
	"fmt"
	"time"
)

// AlertType represents the kind of alert being sent.
type AlertType string

const (
	AlertTypeDrift  AlertType = "drift"
	AlertTypeMissed AlertType = "missed"
)

// Alert holds the data for a single alert event.
type Alert struct {
	JobName   string
	Type      AlertType
	Message   string
	OccurredAt time.Time
}

// Notifier is the interface implemented by alert backends (webhook, email, etc.).
type Notifier interface {
	Send(a Alert) error
}

// NewDriftAlert creates an alert for a job that has drifted beyond its threshold.
func NewDriftAlert(jobName string, expected, actual time.Duration) Alert {
	return Alert{
		JobName:    jobName,
		Type:       AlertTypeDrift,
		Message:    fmt.Sprintf("job %q drifted: expected ~%s, got %s", jobName, expected.Round(time.Second), actual.Round(time.Second)),
		OccurredAt: time.Now().UTC(),
	}
}

// NewMissedAlert creates an alert for a job that did not run within its schedule.
func NewMissedAlert(jobName string, lastSeen time.Time) Alert {
	return Alert{
		JobName:    jobName,
		Type:       AlertTypeMissed,
		Message:    fmt.Sprintf("job %q missed its scheduled run (last seen: %s)", jobName, lastSeen.UTC().Format(time.RFC3339)),
		OccurredAt: time.Now().UTC(),
	}
}
