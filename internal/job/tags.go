package job

import "strings"

// Tag represents a label attached to a job for filtering and grouping.
type Tag string

// TagSet holds a deduplicated collection of tags.
type TagSet map[Tag]struct{}

// NewTagSet creates a TagSet from a slice of strings.
func NewTagSet(tags []string) TagSet {
	ts := make(TagSet, len(tags))
	for _, t := range tags {
		normalized := Tag(strings.TrimSpace(strings.ToLower(t)))
		if normalized != "" {
			ts[normalized] = struct{}{}
		}
	}
	return ts
}

// Contains reports whether the TagSet includes the given tag.
func (ts TagSet) Contains(tag string) bool {
	_, ok := ts[Tag(strings.TrimSpace(strings.ToLower(tag)))]
	return ok
}

// Slice returns the tags as a sorted slice of strings.
func (ts TagSet) Slice() []string {
	out := make([]string, 0, len(ts))
	for t := range ts {
		out = append(out, string(t))
	}
	return out
}

// FilterByTag returns only the jobs whose TagSet contains the given tag.
func FilterByTag(jobs []*Job, tag string) []*Job {
	var result []*Job
	for _, j := range jobs {
		if j.Tags.Contains(tag) {
			result = append(result, j)
		}
	}
	return result
}
