package job

import (
	"testing"
)

func TestNewAnnotationSet_NormalizesKeys(t *testing.T) {
	a := NewAnnotationSet(map[string]string{
		"Owner":   "alice",
		"  Team ": "platform",
	})
	if v, ok := a.Get("owner"); !ok || v != "alice" {
		t.Errorf("expected owner=alice, got %q ok=%v", v, ok)
	}
	if v, ok := a.Get("team"); !ok || v != "platform" {
		t.Errorf("expected team=platform, got %q ok=%v", v, ok)
	}
}

func TestNewAnnotationSet_IgnoresBlankKeys(t *testing.T) {
	a := NewAnnotationSet(map[string]string{
		"":    "ignored",
		"  ": "also ignored",
		"ok": "kept",
	})
	if a.Len() != 1 {
		t.Errorf("expected 1 annotation, got %d", a.Len())
	}
}

func TestAnnotationSet_Get_CaseInsensitive(t *testing.T) {
	a := NewAnnotationSet(map[string]string{"Region": "us-east-1"})
	v, ok := a.Get("REGION")
	if !ok || v != "us-east-1" {
		t.Errorf("expected us-east-1, got %q ok=%v", v, ok)
	}
}

func TestAnnotationSet_Set_And_Delete(t *testing.T) {
	var a AnnotationSet
	a.Set("env", "production")
	if v, ok := a.Get("env"); !ok || v != "production" {
		t.Errorf("expected env=production after Set")
	}
	a.Delete("env")
	if _, ok := a.Get("env"); ok {
		t.Error("expected env to be deleted")
	}
}

func TestAnnotationSet_Map_ReturnsCopy(t *testing.T) {
	a := NewAnnotationSet(map[string]string{"k": "v"})
	m := a.Map()
	m["k"] = "mutated"
	if v, _ := a.Get("k"); v != "v" {
		t.Error("Map() should return a copy, not a reference")
	}
}

func TestAnnotationSet_Set_BlankKeyIgnored(t *testing.T) {
	var a AnnotationSet
	a.Set("", "noop")
	if a.Len() != 0 {
		t.Errorf("expected 0 annotations after blank key Set, got %d", a.Len())
	}
}
