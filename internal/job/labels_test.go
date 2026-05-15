package job

import (
	"testing"
)

func TestNewLabelSet_NormalizesKeys(t *testing.T) {
	ls := NewLabelSet(map[string]string{
		"  Env  ": "production",
		"TEAM":    "platform",
	})
	if v, ok := ls.Get("env"); !ok || v != "production" {
		t.Errorf("expected env=production, got %q ok=%v", v, ok)
	}
	if v, ok := ls.Get("team"); !ok || v != "platform" {
		t.Errorf("expected team=platform, got %q ok=%v", v, ok)
	}
}

func TestNewLabelSet_IgnoresBlankKeys(t *testing.T) {
	ls := NewLabelSet(map[string]string{
		"   ": "ignored",
		"env": "staging",
	})
	if len(ls) != 1 {
		t.Errorf("expected 1 label, got %d", len(ls))
	}
}

func TestLabelSet_Get_CaseInsensitive(t *testing.T) {
	ls := NewLabelSet(map[string]string{"region": "us-east-1"})
	v, ok := ls.Get("REGION")
	if !ok || v != "us-east-1" {
		t.Errorf("expected us-east-1, got %q ok=%v", v, ok)
	}
}

func TestLabelSet_Keys_Sorted(t *testing.T) {
	ls := NewLabelSet(map[string]string{"z": "1", "a": "2", "m": "3"})
	keys := ls.Keys()
	expected := []string{"a", "m", "z"}
	for i, k := range expected {
		if keys[i] != k {
			t.Errorf("expected keys[%d]=%q, got %q", i, k, keys[i])
		}
	}
}

func TestLabelSet_Matches_AllPresent(t *testing.T) {
	ls := NewLabelSet(map[string]string{"env": "prod", "team": "ops"})
	if !ls.Matches(map[string]string{"env": "prod"}) {
		t.Error("expected match")
	}
}

func TestLabelSet_Matches_MissingKey(t *testing.T) {
	ls := NewLabelSet(map[string]string{"env": "prod"})
	if ls.Matches(map[string]string{"team": "ops"}) {
		t.Error("expected no match")
	}
}

func TestFilterByLabels_ReturnsMatching(t *testing.T) {
	jobs := []*Job{
		{Name: "a", Labels: NewLabelSet(map[string]string{"env": "prod"})},
		{Name: "b", Labels: NewLabelSet(map[string]string{"env": "staging"})},
		{Name: "c", Labels: NewLabelSet(map[string]string{"env": "prod", "team": "ops"})},
	}
	result := FilterByLabels(jobs, map[string]string{"env": "prod"})
	if len(result) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(result))
	}
}

func TestFilterByLabels_EmptySelector_ReturnsAll(t *testing.T) {
	jobs := []*Job{
		{Name: "a"},
		{Name: "b"},
	}
	result := FilterByLabels(jobs, map[string]string{})
	if len(result) != 2 {
		t.Errorf("expected 2 jobs, got %d", len(result))
	}
}
