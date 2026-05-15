package job

import (
	"fmt"
	"sync"
)

// Store holds job definitions and their run history.
type Store struct {
	mu      sync.RWMutex
	jobs    map[string]*Job
	history map[string][]Run
}

// NewStore creates an empty Store.
func NewStore() *Store {
	return &Store{
		jobs:    make(map[string]*Job),
		history: make(map[string][]Run),
	}
}

// Add registers a job in the store.
func (s *Store) Add(j *Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs[j.Name] = j
}

// Get retrieves a job by name.
func (s *Store) Get(name string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[name]
	if !ok {
		return nil, fmt.Errorf("job %q not found", name)
	}
	return j, nil
}

// List returns all registered jobs.
func (s *Store) List() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		out = append(out, j)
	}
	return out
}

// RecordRun appends a run to the job's history and updates the job's LastRun.
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

// GetHistory returns all recorded runs for a job.
func (s *Store) GetHistory(name string) ([]Run, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.jobs[name]; !ok {
		return nil, fmt.Errorf("job %q not found", name)
	}
	runs := s.history[name]
	out := make([]Run, len(runs))
	copy(out, runs)
	return out, nil
}
