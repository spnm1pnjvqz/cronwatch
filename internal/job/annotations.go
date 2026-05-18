package job

import "strings"

// AnnotationSet holds arbitrary key-value metadata attached to a job.
// Keys are normalised to lowercase; values are stored as-is.
type AnnotationSet struct {
	data map[string]string
}

// NewAnnotationSet constructs an AnnotationSet from a raw map.
// Blank keys are silently dropped; keys are lowercased.
func NewAnnotationSet(raw map[string]string) AnnotationSet {
	data := make(map[string]string, len(raw))
	for k, v := range raw {
		k = strings.TrimSpace(strings.ToLower(k))
		if k == "" {
			continue
		}
		data[k] = v
	}
	return AnnotationSet{data: data}
}

// Get returns the value for a key (case-insensitive) and whether it exists.
func (a AnnotationSet) Get(key string) (string, bool) {
	v, ok := a.data[strings.ToLower(key)]
	return v, ok
}

// Set stores or overwrites a key-value pair.
func (a *AnnotationSet) Set(key, value string) {
	key = strings.TrimSpace(strings.ToLower(key))
	if key == "" {
		return
	}
	if a.data == nil {
		a.data = make(map[string]string)
	}
	a.data[key] = value
}

// Delete removes a key from the set.
func (a *AnnotationSet) Delete(key string) {
	delete(a.data, strings.ToLower(key))
}

// Map returns a shallow copy of the underlying map.
func (a AnnotationSet) Map() map[string]string {
	out := make(map[string]string, len(a.data))
	for k, v := range a.data {
		out[k] = v
	}
	return out
}

// Len returns the number of annotations stored.
func (a AnnotationSet) Len() int {
	return len(a.data)
}
