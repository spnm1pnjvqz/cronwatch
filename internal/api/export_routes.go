package api

import "net/http"

// RegisterExportRoutes attaches the history export endpoint to the given mux.
// The export endpoint supports optional ?format=csv query parameter;
// defaults to JSON when the parameter is absent or unrecognised.
func RegisterExportRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("GET /jobs/{name}/export", h.exportJobHistory)
}
