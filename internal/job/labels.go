package job

import (
	"sort"
	"strings"
)

// LabelSet holds key-value metadata labels for a job.
type LabelSet map[string]string

// NewLabelSet creates a LabelSet from a raw map, normalizing keys to lowercase
// and trimming whitespace from both keys and values.
func NewLabelSet(raw map[string]string) LabelSet {
	ls := make(LabelSet, len(raw))
	for k, v := range raw {
		key := strings.ToLower(strings.TrimSpace(k))
		val := strings.TrimSpace(v)
		if key != "" {
			ls[key] = val
		}
	}
	return ls
}

// Get returns the value for a label key (case-insensitive).
func (ls LabelSet) Get(key string) (string, bool) {
	v, ok := ls[strings.ToLower(strings.TrimSpace(key))]
	return v, ok
}

// Keys returns a sorted slice of all label keys.
func (ls LabelSet) Keys() []string {
	keys := make([]string, 0, len(ls))
	for k := range ls {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Matches reports whether the LabelSet contains all key-value pairs in the
// given selector map (case-insensitive keys, exact values).
func (ls LabelSet) Matches(selector map[string]string) bool {
	for k, v := range selector {
		got, ok := ls.Get(k)
		if !ok || got != v {
			return false
		}
	}
	return true
}

// FilterByLabels returns jobs whose LabelSet matches all key-value pairs in
// the selector.
func FilterByLabels(jobs []*Job, selector map[string]string) []*Job {
	if len(selector) == 0 {
		return jobs
	}
	out := make([]*Job, 0)
	for _, j := range jobs {
		if j.Labels.Matches(selector) {
			out = append(out, j)
		}
	}
	return out
}
