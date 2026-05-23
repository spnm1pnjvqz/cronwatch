package job

import (
	"fmt"
	"sync"
)

// Store holds jobs and their run history in memory.
type Store struct {
	mu      sync.RWMutex
	jobs    map[string]*Job
	history map[string][]Run
	paused  map[string]bool
}

func NewStore() *Store {
	return &Store{
		jobs:    make(map[string]*Job),
		history: make(map[string][]Run),
		paused:  make(map[string]bool),
	}
}

func (s *Store) Add(j Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.Name] = &j
}

func (s *Store) Get(name string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job %q not found", name)
	}
	return j, nil
}

func (s *Store) List() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		out = append(out, j)
	}
	return out
}

func (s *Store) RecordRun(name string, r Run) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	j, ok := s.jobs[name]
	if !ok {
		return fmt.Errorf("job %q not found", name)
	}
	j.LastRun = &r
	s.history[name] = append(s.history[name], r)
	return nil
}

func (s *Store) GetHistory(name string) ([]Run, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.jobs[name]; !ok {
		return nil, fmt.Errorf("job %q not found", name)
	}
	return s.history[name], nil
}

// SetPaused sets the paused state for a job. Returns error if job not found.
func (s *Store) SetPaused(name string, paused bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; !ok {
		return fmt.Errorf("job %q not found", name)
	}
	s.paused[name] = paused
	s.jobs[name].Paused = paused
	return nil
}

// IsPaused reports whether the named job is currently paused.
func (s *Store) IsPaused(name string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.paused[name]
}

// Delete removes a job and its run history from the store.
// Returns an error if the job does not exist.
func (s *Store) Delete(name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[name]; !ok {
		return fmt.Errorf("job %q not found", name)
	}
	delete(s.jobs, name)
	delete(s.history, name)
	delete(s.paused, name)
	return nil
}
