package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// JobConfig holds configuration for a single monitored cron job.
type JobConfig struct {
	Name           string        `yaml:"name"`
	Schedule       string        `yaml:"schedule"`
	DriftThreshold time.Duration `yaml:"drift_threshold"`
}

// WebhookConfig holds configuration for webhook notifications.
type WebhookConfig struct {
	URL string `yaml:"url"`
}

// EmailConfig holds configuration for email notifications.
type EmailConfig struct {
	SMTPHost string `yaml:"smtp_host"`
	SMTPPort int    `yaml:"smtp_port"`
	From     string `yaml:"from"`
	To       string `yaml:"to"`
}

// NotifierConfig holds all notifier configurations.
type NotifierConfig struct {
	Webhook *WebhookConfig `yaml:"webhook,omitempty"`
	Email   *EmailConfig   `yaml:"email,omitempty"`
}

// Config is the top-level application configuration.
type Config struct {
	Jobs      []JobConfig    `yaml:"jobs"`
	Notifiers NotifierConfig `yaml:"notifiers"`
}

// Load reads and parses a YAML configuration file from the given path.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading config file %q: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config file %q: %w", path, err)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

// validate checks that the configuration is semantically valid.
func (c *Config) validate() error {
	if len(c.Jobs) == 0 {
		return fmt.Errorf("at least one job must be configured")
	}
	for i, j := range c.Jobs {
		if j.Name == "" {
			return fmt.Errorf("job[%d]: name is required", i)
		}
		if j.Schedule == "" {
			return fmt.Errorf("job[%d] %q: schedule is required", i, j.Name)
		}
	}
	return nil
}
