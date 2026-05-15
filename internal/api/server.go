package api

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/cronwatch/internal/job"
)

// Server wraps an HTTP server for the cronwatch API.
type Server struct {
	httpServer *http.Server
}

// NewServer creates and configures the HTTP server.
func NewServer(addr string, store *job.Store) *Server {
	mux := http.NewServeMux()
	h := NewHandler(store)
	RegisterRoutes(mux, h)

	return &Server{
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      mux,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
			IdleTimeout:  30 * time.Second,
		},
	}
}

// Start begins listening and serving HTTP requests.
func (s *Server) Start() error {
	log.Printf("cronwatch API listening on %s", s.httpServer.Addr)
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server with the provided context.
func (s *Server) Shutdown(ctx context.Context) error {
	log.Println("shutting down API server")
	return s.httpServer.Shutdown(ctx)
}
