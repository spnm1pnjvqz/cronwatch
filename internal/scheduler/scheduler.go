package scheduler

import (
	"log"
	"sync"
	"time"

	"github.com/user/cronwatch/internal/alert"
	"github.com/user/cronwatch/internal/job"
)

// Scheduler periodically checks job stores for missed runs and drift.
type Scheduler struct {
	store     *job.Store
	notifier  alert.Notifier
	interval  time.Duration
	stopCh    chan struct{}
	wg        sync.WaitGroup
}

// NewScheduler creates a Scheduler that checks the given store on the
// provided interval and dispatches alerts via notifier.
func NewScheduler(store *job.Store, notifier alert.Notifier, interval time.Duration) *Scheduler {
	return &Scheduler{
		store:    store,
		notifier: notifier,
		interval: interval,
		stopCh:   make(chan struct{}),
	}
}

// Start begins the background checking loop.
func (s *Scheduler) Start() {
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				s.check()
			case <-s.stopCh:
				return
			}
		}
	}()
}

// Stop signals the background loop to exit and waits for it to finish.
func (s *Scheduler) Stop() {
	close(s.stopCh)
	s.wg.Wait()
}

// check evaluates all jobs in the store and sends appropriate alerts.
func (s *Scheduler) check() {
	now := time.Now()
	jobs := s.store.All()
	for _, j := range jobs {
		s.checkJob(j, now)
	}
}

// checkJob evaluates a single job and sends a missed or drift alert as needed.
func (s *Scheduler) checkJob(j *job.Job, now time.Time) {
	if j.IsMissed(now) {
		a := alert.NewMissedAlert(j)
		if err := s.notifier.Send(a); err != nil {
			log.Printf("scheduler: failed to send missed alert for job %q: %v", j.Name, err)
		}
		return
	}
	if j.HasDrift() {
		a := alert.NewDriftAlert(j)
		if err := s.notifier.Send(a); err != nil {
			log.Printf("scheduler: failed to send drift alert for job %q: %v", j.Name, err)
		}
	}
}
