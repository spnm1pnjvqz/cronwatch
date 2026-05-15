package job

import (
	"sort"
	"testing"
)

func TestNewTagSet_NormalizesAndDeduplicates(t *testing.T) {
	ts := NewTagSet([]string{"Prod", "prod", " PROD ", "staging"})
	if !ts.Contains("prod") {
		t.Error("expected 'prod' to be present")
	}
	if !ts.Contains("staging") {
		t.Error("expected 'staging' to be present")
	}
	slice := ts.Slice()
	if len(slice) != 2 {
		t.Errorf("expected 2 unique tags, got %d", len(slice))
	}
}

func TestNewTagSet_IgnoresBlanks(t *testing.T) {
	ts := NewTagSet([]string{"", "  ", "valid"})
	slice := ts.Slice()
	if len(slice) != 1 {
		t.Errorf("expected 1 tag, got %d", len(slice))
	}
}

func TestTagSet_Contains_CaseInsensitive(t *testing.T) {
	ts := NewTagSet([]string{"Critical"})
	if !ts.Contains("CRITICAL") {
		t.Error("Contains should be case-insensitive")
	}
}

func TestTagSet_Slice_ReturnsCopy(t *testing.T) {
	ts := NewTagSet([]string{"a", "b", "c"})
	s := ts.Slice()
	sort.Strings(s)
	if len(s) != 3 {
		t.Errorf("expected 3 items, got %d", len(s))
	}
}

func TestFilterByTag_ReturnsMatchingJobs(t *testing.T) {
	jobs := []*Job{
		{Name: "job-a", Tags: NewTagSet([]string{"prod", "critical"})},
		{Name: "job-b", Tags: NewTagSet([]string{"staging"})},
		{Name: "job-c", Tags: NewTagSet([]string{"prod"})},
	}
	result := FilterByTag(jobs, "prod")
	if len(result) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(result))
	}
}

func TestFilterByTag_NoMatches(t *testing.T) {
	jobs := []*Job{
		{Name: "job-a", Tags: NewTagSet([]string{"staging"})},
	}
	result := FilterByTag(jobs, "prod")
	if len(result) != 0 {
		t.Errorf("expected 0 jobs, got %d", len(result))
	}
}
