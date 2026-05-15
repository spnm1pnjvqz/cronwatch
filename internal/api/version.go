package api

import (
	"encoding/json"
	"net/http"
	"runtime"
)

// BuildInfo holds version and build metadata for the service.
type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	BuildTime string `json:"build_time"`
	GoVersion string `json:"go_version"`
}

// These variables are intended to be set via -ldflags at build time.
var (
	Version   = "dev"
	Commit    = "none"
	BuildTime = "unknown"
)

// VersionHandler returns a handler that responds with build info as JSON.
func VersionHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info := BuildInfo{
			Version:   Version,
			Commit:    Commit,
			BuildTime: BuildTime,
			GoVersion: runtime.Version(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(info); err != nil {
			http.Error(w, "failed to encode version info", http.StatusInternalServerError)
		}
	}
}
