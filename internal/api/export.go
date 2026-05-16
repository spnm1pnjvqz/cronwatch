package api

import (
	"encoding/csv"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/user/cronwatch/internal/job"
)

type exportRow struct {
	JobName   string `json:"job_name"`
	RunAt     string `json:"run_at"`
	Duration  string `json:"duration_ms"`
	Drifted   string `json:"drifted"`
	DriftMs   string `json:"drift_ms"`
}

func (h *Handler) exportJobHistory(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		http.Error(w, "missing job name", http.StatusBadRequest)
		return
	}

	j, ok := h.store.Get(name)
	if !ok {
		http.Error(w, "job not found", http.StatusNotFound)
		return
	}

	runs := h.store.GetHistory(j.Name)
	format := r.URL.Query().Get("format")

	switch format {
	case "csv":
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", `attachment; filename="history.csv"`)
		cw := csv.NewWriter(w)
		_ = cw.Write([]string{"job_name", "run_at", "duration_ms", "drifted", "drift_ms"})
		for _, run := range runs {
			driftedStr := "false"
			driftMsStr := "0"
			if run.Drifted {
				driftedStr = "true"
				driftMsStr = strconv.FormatInt(run.DriftAmount.Milliseconds(), 10)
			}
			_ = cw.Write([]string{
				j.Name,
				run.RunAt.Format(time.RFC3339),
				strconv.FormatInt(run.Duration.Milliseconds(), 10),
				driftedStr,
				driftMsStr,
			})
		}
		cw.Flush()
	default:
		w.Header().Set("Content-Type", "application/json")
		rows := make([]exportRow, 0, len(runs))
		for _, run := range runs {
			driftedStr := "false"
			driftMsStr := "0"
			if run.Drifted {
				driftedStr = "true"
				driftMsStr = strconv.FormatInt(run.DriftAmount.Milliseconds(), 10)
			}
			rows = append(rows, exportRow{
				JobName:  j.Name,
				RunAt:    run.RunAt.Format(time.RFC3339),
				Duration: strconv.FormatInt(run.Duration.Milliseconds(), 10),
				Drifted:  driftedStr,
				DriftMs:  driftMsStr,
			})
		}
		_ = json.NewEncoder(w).Encode(rows)
	}
}

// ensure job.Run is referenced
var _ = (*job.Run)(nil)
