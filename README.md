# cronwatch

Monitors cron job execution times and alerts on drift or missed runs via webhook or email.

## Installation

```bash
go install github.com/yourusername/cronwatch@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronwatch.git && cd cronwatch && go build ./...
```

## Usage

Define your jobs in a `cronwatch.yaml` config file:

```yaml
jobs:
  - name: daily-backup
    schedule: "0 2 * * *"
    tolerance: 5m
    alert:
      webhook: "https://hooks.example.com/notify"
      email: "ops@example.com"

  - name: hourly-sync
    schedule: "0 * * * *"
    tolerance: 2m
```

Start the watcher:

```bash
cronwatch --config cronwatch.yaml
```

Ping cronwatch from your cron job to record execution:

```bash
# In your crontab
0 2 * * * /usr/local/bin/backup.sh && curl -s https://localhost:8080/ping/daily-backup
```

cronwatch will alert you if a job misses its scheduled window or if execution times drift beyond the configured tolerance.

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--config` | `cronwatch.yaml` | Path to config file |
| `--port` | `8080` | HTTP listener port |
| `--log-level` | `info` | Log verbosity |

## License

MIT © 2024 cronwatch contributors