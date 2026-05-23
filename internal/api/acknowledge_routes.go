package api

import "net/http"

// RegisterAcknowledgeRoutes wires acknowledge/unacknowledge endpoints
// onto the provided ServeMux.
//
//	POST   /jobs/{name}/acknowledge   — mark a job alert as acknowledged
//	DELETE /jobs/{name}/acknowledge   — clear acknowledgement
func RegisterAcknowledgeRoutes(mux *http.ServeMux, h *Handler) {
	mux.HandleFunc("POST /jobs/{name}/acknowledge", h.acknowledgeJob)
	mux.HandleFunc("DELETE /jobs/{name}/acknowledge", h.unacknowledgeJob)
}
