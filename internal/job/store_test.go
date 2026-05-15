package job_test

import (
	"testing"
	"time"

	"github.com/cronwatch/internal/job"
)

func TestStore_AddAndGet(t *testing.T) {
	s := job.NewStore()
	j := &job.Job{Name: "backup", Schedule: "@daily"}
	s.Add(j)

	got, err := s.Get("backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "backup" {
		t.Errorf("expected backup, got %s", got.Name)
	}
}

func TestStore_Get_NotFound(t *testing.T) {
	s := job.NewStore()
	_, err := s.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing job")
	}
}

func TestStore_List(t *testing.T) {
	s := job.NewStore()
	s.Add(&job.Job{Name: "a"})
	s.Add(&job.Job{Name: "b"})

	list := s.List()
	if len(list) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(list))
	}
}

func TestStore_RecordRun_UpdatesLastRun(t *testing.T) {
	s := job.NewStore()
	s.Add(&job.Job{Name: "sync"})

	now := time.Now()
	r := job.Run{StartedAt: now, Duration: time.Second, Status: "ok"}
	if err := s.RecordRun("sync", r); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	j, _ := s.Get("sync")
	if j.LastRun == nil {
		t.Fatal("expected LastRun to be set")
	}
	if !j.LastRun.StartedAt.Equal(now) {
		t.Errorf("unexpected StartedAt: %v", j.LastRun.StartedAt)
	}
}

func TestStore_GetHistory(t *testing.T) {
	s := job.NewStore()
	s.Add(&job.Job{Name: "report"})

	for i := 0; i < 3; i++ {
		s.RecordRun("report", job.Run{StartedAt: time.Now(), Status: "ok"})
	}

	runs, err := s.GetHistory("report")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(runs) != 3 {
		t.Errorf("expected 3 runs, got %d", len(runs))
	}
}

func TestStore_GetHistory_NotFound(t *testing.T) {
	s := job.NewStore()
	_, err := s.GetHistory("ghost")
	if err == nil {
		t.Fatal("expected error for missing job")
	}
}
