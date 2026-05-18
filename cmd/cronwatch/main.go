// main is the entry point for the cronwatch daemon.
// It loads configuration, initializes the job store, wires up notifiers,
// starts the scheduler, and serves the HTTP API.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/user/cronwatch/internal/alert"
	"github.com/user/cronwatch/internal/api"
	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/job"
	"github.com/user/cronwatch/internal/scheduler"
)

func main() {
	cfgPath := flag.String("config", "cronwatch.yaml", "path to configuration file")
	flag.Parse()

	// Load configuration from disk.
	cfg, err := config.Load(*cfgPath)
	if err != nil {
		log.Fatalf("cronwatch: failed to load config: %v", err)
	}

	// Build the in-memory job store and register all configured jobs.
	store := job.NewStore()
	for _, jcfg := range cfg.Jobs {
		j := job.Job{
			ID:             jcfg.Name,
			Name:           jcfg.Name,
			Schedule:       jcfg.Schedule,
			DriftThreshold: time.Duration(jcfg.DriftThreshold),
		}
		if err := store.Add(j); err != nil {
			log.Fatalf("cronwatch: failed to register job %q: %v", jcfg.Name, err)
		}
	}
	log.Printf("cronwatch: registered %d job(s)", len(cfg.Jobs))

	// Build the notifier chain from configured alert targets.
	var notifiers []alert.Notifier
	if cfg.Alerts.Webhook.URL != "" {
		notifiers = append(notifiers, alert.NewWebhookNotifier(cfg.Alerts.Webhook.URL))
		log.Printf("cronwatch: webhook notifier enabled (%s)", cfg.Alerts.Webhook.URL)
	}
	if cfg.Alerts.Email.To != "" {
		en, err := alert.NewEmailNotifier(
			cfg.Alerts.Email.SMTPHost,
			cfg.Alerts.Email.SMTPPort,
			cfg.Alerts.Email.From,
			cfg.Alerts.Email.To,
		)
		if err != nil {
			log.Fatalf("cronwatch: failed to create email notifier: %v", err)
		}
		notifiers = append(notifiers, en)
		log.Printf("cronwatch: email notifier enabled (to: %s)", cfg.Alerts.Email.To)
	}
	notifier := alert.NewMultiNotifier(notifiers...)

	// Start the scheduler.
	sched := scheduler.NewScheduler(store, notifier)
	schedCtx, schedCancel := context.WithCancel(context.Background())
	defer schedCancel()
	go sched.Run(schedCtx)
	log.Println("cronwatch: scheduler started")

	// Start the HTTP API server.
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := api.NewServer(addr, store, cfg)
	httpServer := &http.Server{
		Addr:         addr,
		Handler:      srv,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		log.Printf("cronwatch: API listening on %s", addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("cronwatch: HTTP server error: %v", err)
		}
	}()

	// Wait for termination signal, then gracefully shut down.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("cronwatch: shutting down...")

	schedCancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("cronwatch: graceful shutdown error: %v", err)
	}
	log.Println("cronwatch: stopped")
}
