package alert

import (
	"errors"
	"testing"
	"time"
)

// mockNotifier records sent alerts and optionally returns an error.
type mockNotifier struct {
	sentAlerts  []Alert
	errToReturn error
}

func (m *mockNotifier) Send(a Alert) error {
	m.sentAlerts = append(m.sentAlerts, a)
	return m.errToReturn
}

func TestMultiNotifier_Send_AllSucceed(t *testing.T) {
	n1 := &mockNotifier{}
	n2 := &mockNotifier{}

	multi := NewMultiNotifier(n1, n2)
	a := NewDriftAlert("sync-job", 2*time.Minute, 4*time.Minute)

	if err := multi.Send(a); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(n1.sentAlerts) != 1 {
		t.Errorf("expected n1 to receive 1 alert, got %d", len(n1.sentAlerts))
	}
	if len(n2.sentAlerts) != 1 {
		t.Errorf("expected n2 to receive 1 alert, got %d", len(n2.sentAlerts))
	}
}

func TestMultiNotifier_Send_PartialFailure(t *testing.T) {
	sentinel := errors.New("delivery failed")
	n1 := &mockNotifier{errToReturn: sentinel}
	n2 := &mockNotifier{}

	multi := NewMultiNotifier(n1, n2)
	a := NewMissedAlert("report-job", time.Now().Add(-10*time.Minute))

	err := multi.Send(a)
	if !errors.Is(err, sentinel) {
		t.Errorf("expected sentinel error, got %v", err)
	}
	// n2 should still have received the alert despite n1 failing
	if len(n2.sentAlerts) != 1 {
		t.Errorf("expected n2 to receive 1 alert even after n1 error, got %d", len(n2.sentAlerts))
	}
}

func TestMultiNotifier_Send_Empty(t *testing.T) {
	multi := NewMultiNotifier()
	a := NewDriftAlert("empty-job", time.Minute, 2*time.Minute)

	if err := multi.Send(a); err != nil {
		t.Errorf("expected no error with zero notifiers, got %v", err)
	}
}

func TestMultiNotifier_Send_AllFail(t *testing.T) {
	sentinel1 := errors.New("notifier 1 failed")
	sentinel2 := errors.New("notifier 2 failed")
	n1 := &mockNotifier{errToReturn: sentinel1}
	n2 := &mockNotifier{errToReturn: sentinel2}

	multi := NewMultiNotifier(n1, n2)
	a := NewDriftAlert("batch-job", time.Minute, 3*time.Minute)

	err := multi.Send(a)
	if err == nil {
		t.Fatal("expected an error when all notifiers fail, got nil")
	}
	if !errors.Is(err, sentinel1) && !errors.Is(err, sentinel2) {
		t.Errorf("expected combined error to wrap at least one sentinel, got %v", err)
	}
}
