# Prometheus Minimal Core Build

This document describes how to build and use a minimal version of Prometheus optimized for resource-constrained environments.

## Overview

The minimal Prometheus build reduces binary size by **~72%** (from 153MB to 43MB) by removing:
- All cloud and Kubernetes service discovery mechanisms
- Remote read/write capabilities
- OTLP ingestion endpoints

## What's Included

The minimal build retains core Prometheus functionality:
- **Scrape**: Metrics collection engine with `static_configs` and `file_sd`
- **TSDB**: Local time-series database with block storage and WAL
- **PromQL**: Full query language support via HTTP API
- **Rules**: Alert rule evaluation engine
- **Notifier**: Alert notification via webhook
- **HTTP API**: Query endpoints (`/api/v1/query`, `/api/v1/query_range`, etc.)

## What's Excluded

- **Service Discovery**: AWS, Azure, GCP, Kubernetes, Consul, Eureka, and all other cloud SDs
- **Remote Storage**: No remote_write or remote_read support
- **OTLP**: No OpenTelemetry ingestion

## Building

### Using go build

```bash
go build \
  -tags "minimal,remove_all_sd" \
  -ldflags="-s -w" \
  -trimpath \
  -o prometheus-minimal \
  ./cmd/prometheus
```

### Build Tags

- `minimal`: Disables remote storage (remote read/write)
- `remove_all_sd`: Disables all cloud service discovery (keeps file_sd and static_configs)

### Using Docker

```bash
docker build -f Dockerfile.minimal -t prometheus-minimal:latest .
```

Expected image size: **~25-35 MB** (including distroless base)

## Running

The minimal build is a drop-in replacement for standard Prometheus when you don't need cloud SD or remote storage:

```bash
./prometheus-minimal --config.file=prometheus-minimal.yml
```

### Configuration Example

```yaml
global:
  scrape_interval: 15s

scrape_configs:
  # Static targets
  - job_name: "node"
    static_configs:
      - targets: ["localhost:9100"]

  # File-based service discovery
  - job_name: "app"
    file_sd_configs:
      - files:
          - "/etc/prometheus/targets/*.json"
        refresh_interval: 30s
```

See `documentation/examples/prometheus-minimal.yml` for a complete example.

## Size Comparison

| Build Type | Binary Size | Description |
|------------|-------------|-------------|
| Full       | 153 MB      | All features enabled |
| Minimal    | 43 MB       | With `-tags "minimal,remove_all_sd"` |
| Reduction  | **72%**     | Size savings |

## Use Cases

The minimal build is ideal for:

- **Edge Devices**: Resource-constrained edge nodes
- **Development**: Fast local development environments
- **Embedded Systems**: Appliances and IoT devices
- **Containers**: Reduced container image sizes
- **CI/CD**: Temporary monitoring in pipelines
- **Single-Node Deployments**: Simple standalone monitoring

## Not Suitable For

Avoid the minimal build if you need:
- Kubernetes or cloud provider service discovery
- Remote storage integration (Thanos, Cortex, etc.)
- Long-term storage with downsampling
- Multi-cluster federation
- OTLP ingestion

## File-Based Service Discovery

Since cloud SD is unavailable, use file_sd for dynamic targets:

### Example target file (`/etc/prometheus/targets/app.json`):

```json
[
  {
    "targets": ["10.0.1.10:9100", "10.0.1.11:9100"],
    "labels": {
      "job": "node",
      "env": "production"
    }
  }
]
```

The file can be updated dynamically, and Prometheus will reload targets based on the `refresh_interval`.

## API Limitations

The following API endpoints return `501 Not Implemented` in minimal builds:
- Remote read endpoint
- Remote write endpoint
- OTLP write endpoint

All query APIs (`/api/v1/query`, `/api/v1/query_range`, etc.) work normally.

## Performance

The minimal build:
- Uses less memory (no remote write queues)
- Faster startup (fewer dependencies to initialize)
- Lower CPU usage (no cloud API calls)
- Identical query performance to full build

## Upgrading from Full Build

1. Remove `remote_write` and `remote_read` sections from config
2. Replace cloud SD with `static_configs` or `file_sd_configs`
3. Build/deploy minimal binary
4. Restart with updated configuration

## Integration with Alerting

Alert notifications work via webhook in minimal builds. Configure your alert rules normally:

```yaml
rule_files:
  - "alerts/*.yml"

alerting:
  alertmanagers:
    - static_configs:
        - targets: ["localhost:9093"]
```

## Contributing

When modifying code that affects the minimal build:
- Add `//go:build !minimal` to files that should be excluded
- Create corresponding `_minimal.go` stub files when needed
- Test both `go build` and `go build -tags minimal` succeed

## References

- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Prometheus Configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
- Inspired by: [Can I have a smaller Prometheus?](https://wejick.wordpress.com/2022/01/29/can-i-have-a-smaller-prometheus/)
