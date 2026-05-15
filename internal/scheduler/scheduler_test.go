package scheduler_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/alert"
	"github.com/user/cronwatch/internal/job"
	"github.com/user/cronwatch/internal/scheduler"
)

// mockNotifier records every alert it receives.
type mockNotifier struct {
	mu     sync.Mutex
	alerts []alert.Alert
	err    error
}

func (m *mockNotifier) Send(a alert.Alert) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.alerts = append(m.alerts, a)
	return m.err
}

func (m *mockNotifier) count() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.alerts)
}

func TestScheduler_DetectsMissedRun(t *testing.T) {
	store := job.NewStore()
	j := job.Job{
		Name:     "backup",
		Schedule: "@daily",
		// LastRun two days ago — definitely missed
		LastRun:  time.Now().Add(-48 * time.Hour),
		Interval: 24 * time.Hour,
	}
	store.Add(j)

	notifier := &mockNotifier{}
	sched := scheduler.NewScheduler(store, notifier, 50*time.Millisecond)
	sched.Start()
	time.Sleep(120 * time.Millisecond)
	sched.Stop()

	if notifier.count() == 0 {
		t.Error("expected at least one missed-run alert, got none")
	}
}

func TestScheduler_NoAlertsForHealthyJob(t *testing.T) {
	store := job.NewStore()
	j := job.Job{
		Name:     "heartbeat",
		Schedule: "@hourly",
		LastRun:  time.Now().Add(-30 * time.Minute),
		Interval: time.Hour,
	}
	store.Add(j)

	notifier := &mockNotifier{}
	sched := scheduler.NewScheduler(store, notifier, 50*time.Millisecond)
	sched.Start()
	time.Sleep(120 * time.Millisecond)
	sched.Stop()

	if notifier.count() != 0 {
		t.Errorf("expected no alerts for healthy job, got %d", notifier.count())
	}
}

func TestScheduler_NotifierErrorDoesNotPanic(t *testing.T) {
	store := job.NewStore()
	j := job.Job{
		Name:     "flaky",
		Schedule: "@daily",
		LastRun:  time.Now().Add(-48 * time.Hour),
		Interval: 24 * time.Hour,
	}
	store.Add(j)

	notifier := &mockNotifier{err: errors.New("connection refused")}
	sched := scheduler.NewScheduler(store, notifier, 50*time.Millisecond)
	sched.Start()
	time.Sleep(120 * time.Millisecond)
	sched.Stop() // must not panic
}
