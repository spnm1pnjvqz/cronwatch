package job

import (
	"fmt"
	"sync"
)

// Store is an in-memory repository for Jobs and their Runs.
type Store struct {
	mu   sync.RWMutex
	jobs map[string]*Job
	runs map[string][]*Run // keyed by job ID
}

// NewStore initialises an empty Store.
func NewStore() *Store {
	return &Store{
		jobs: make(map[string]*Job),
		runs: make(map[string][]*Run),
	}
}

// AddJob registers a new job. Returns an error if the ID already exists.
func (s *Store) AddJob(j *Job) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.jobs[j.ID]; exists {
		return fmt.Errorf("job %q already exists", j.ID)
	}
	s.jobs[j.ID] = j
	return nil
}

// GetJob retrieves a job by ID.
func (s *Store) GetJob(id string) (*Job, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	j, ok := s.jobs[id]
	if !ok {
		return nil, fmt.Errorf("job %q not found", id)
	}
	return j, nil
}

// ListJobs returns all registered jobs.
func (s *Store) ListJobs() []*Job {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*Job, 0, len(s.jobs))
	for _, j := range s.jobs {
		list = append(list, j)
	}
	return list
}

// RecordRun appends a run for the given job.
func (s *Store) RecordRun(r *Run) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.jobs[r.JobID]; !ok {
		return fmt.Errorf("job %q not found", r.JobID)
	}
	s.runs[r.JobID] = append(s.runs[r.JobID], r)
	return nil
}

// GetRuns returns all recorded runs for a job.
func (s *Store) GetRuns(jobID string) ([]*Run, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if _, ok := s.jobs[jobID]; !ok {
		return nil, fmt.Errorf("job %q not found", jobID)
	}
	return s.runs[jobID], nil
}
