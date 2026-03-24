# Prometheus Minimal Core Build

This document describes how to build and use minimal versions of Prometheus optimized for resource-constrained environments.

## Overview

The Prometheus minimal builds reduce binary size through progressive optimization levels:

| Version | Size | Reduction | Method |
|---------|------|-----------|--------|
| **Full** | 153 MB | — | Standard build |
| **V2 Minimal** | 43 MB | 72% | Remove cloud SD + remote storage |
| **V3 + UPX** | **15 MB** | **90%** | V2 + static compilation + UPX compression |

## What's Included

All minimal builds retain core Prometheus functionality:
- **Scrape**: Metrics collection engine with `static_configs` and `file_sd`
- **TSDB**: Local time-series database with block storage and WAL
- **PromQL**: Full query language support via HTTP API
- **Rules**: Alert rule evaluation engine
- **Notifier**: Alert notification via webhook
- **HTTP API**: Query endpoints (`/api/v1/query`, `/api/v1/query_range`, etc.)

## What's Excluded

- **Service Discovery**: AWS, Azure, GCP, Kubernetes, Consul, Eureka, and all other cloud SDs (keeps file_sd and static_configs)
- **Remote Storage**: No remote_write or remote_read support
- **OTLP**: No OpenTelemetry ingestion
- **Web UI**: No embedded React UI (API-only, external Grafana recommended)

## Building

### Quick Start: V3 Build Script

The easiest way to build all variants:

```bash
# Build V2, V3, and V3+UPX
./scripts/build-minimal-v3.sh

# Output in dist/:
# - prometheus-v2-minimal  (43 MB)
# - prometheus-v3-deep     (46 MB, static)
# - prometheus-v3-upx      (15 MB, compressed)
```

### V2 Minimal Build (43 MB)

```bash
go build \
  -tags "minimal,remove_all_sd" \
  -ldflags="-s -w" \
  -trimpath \
  -o prometheus-minimal \
  ./cmd/prometheus
```

### V3 Deep Build (46 MB, Static)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -tags "minimal,remove_all_sd" \
  -ldflags="-s -w" \
  -trimpath \
  -buildmode=pie \
  -o prometheus-v3-deep \
  ./cmd/prometheus
```

### V3 + UPX Compression (15 MB)

Requires [UPX](https://upx.github.io/) to be installed:

```bash
# Install UPX
apt-get install upx          # Debian/Ubuntu
brew install upx             # macOS
pacman -S upx                # Arch Linux

# Build V3 then compress
CGO_ENABLED=0 go build -tags "minimal,remove_all_sd" -ldflags="-s -w" -trimpath -o prometheus ./cmd/prometheus
upx --best -o prometheus-compressed prometheus
```

**UPX Trade-offs:**
- ✅ 65% size reduction (43MB → 15MB)
- ✅ Functionally identical
- ✅ No runtime performance impact (after decompression)
- ⚠️ Startup time +200-500ms (decompression overhead)
- ⚠️ Slightly higher memory usage during startup
- ⚠️ Some security scanners may flag UPX-compressed binaries

**When to use UPX:**
- ✅ Disk/bandwidth constrained environments (IoT, edge, air-gapped)
- ✅ Container images where size matters
- ✅ Long-running processes (startup delay amortized)
- ❌ Frequently restarted services (startup overhead accumulates)
- ❌ Environments with strict security scanning

### Build Tags

- `minimal`: Disables remote storage (remote read/write) and OTLP
- `remove_all_sd`: Disables all cloud service discovery (keeps file_sd and static_configs)

### Cross-Platform Builds

```bash
# ARM64 (Raspberry Pi, edge gateways)
GOOS=linux GOARCH=arm64 ./scripts/build-minimal-v3.sh

# ARMv7 (IoT devices)
GOOS=linux GOARCH=arm GOARM=7 ./scripts/build-minimal-v3.sh

# MIPS (routers, industrial)
GOOS=linux GOARCH=mipsle GOMIPS=softfloat ./scripts/build-minimal-v3.sh
```

ARM/MIPS binaries are typically 10-20% smaller than amd64.

## Docker Images

### V2 Standard Minimal (distroless)

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build \
    -tags "minimal,remove_all_sd" \
    -ldflags="-s -w" \
    -trimpath \
    -o /prometheus \
    ./cmd/prometheus

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /prometheus /bin/prometheus
COPY documentation/examples/prometheus-minimal.yml /etc/prometheus/prometheus.yml
EXPOSE 9090
VOLUME ["/prometheus"]
ENTRYPOINT ["/bin/prometheus"]
CMD ["--config.file=/etc/prometheus/prometheus.yml", \
     "--storage.tsdb.path=/prometheus"]
```

**Image size:** ~48 MB (5MB base + 43MB binary)

### V3 + UPX (scratch)

```dockerfile
FROM golang:1.26-alpine AS builder
WORKDIR /src
COPY . .
RUN apk add --no-cache upx && \
    CGO_ENABLED=0 go build \
    -tags "minimal,remove_all_sd" \
    -ldflags="-s -w" \
    -trimpath \
    -o /prometheus-uncompressed \
    ./cmd/prometheus && \
    upx --best -o /prometheus /prometheus-uncompressed

FROM scratch
COPY --from=builder /prometheus /prometheus
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
EXPOSE 9090
ENTRYPOINT ["/prometheus"]
CMD ["--config.file=/etc/prometheus/prometheus.yml", \
     "--storage.tsdb.path=/prometheus"]
```

**Image size:** ~15 MB (0MB base + 15MB compressed binary)

## Size Comparison

| Build Type | Binary Size | Docker Image | Reduction |
|------------|-------------|--------------|-----------|
| Full (standard) | 153 MB | ~175 MB | — |
| V2 Minimal | 43 MB | ~48 MB | 72% |
| V3 Static | 46 MB | ~46 MB | 70% |
| V3 + UPX | **15 MB** | **~15 MB** | **90%** |

## Performance Characteristics

### V2 Minimal
- No performance impact vs. full build
- Same memory footprint
- Identical query performance

### V3 Static (CGO_ENABLED=0)
- Negligible performance difference (<1%)
- Pure Go crypto/net implementations
- Fully reproducible builds

### V3 + UPX
- Startup time: +200-500ms (one-time decompression)
- Runtime performance: identical after startup
- Memory: +5-10MB during startup, then normal

## Use Cases

### V2 Minimal
- General purpose minimal deployment
- Development environments
- When startup time is critical
- CI/CD pipelines

### V3 Static
- Containers (distroless, scratch)
- Reproducible builds
- Air-gapped environments
- Security-sensitive deployments

### V3 + UPX
- IoT and edge devices
- Bandwidth-constrained deployments
- Embedded systems
- Firmware with size limits
- Long-running services

## Not Suitable For

Avoid minimal builds if you need:
- Kubernetes or cloud provider service discovery
- Remote storage integration (Thanos, Cortex, Mimir)
- Long-term storage with downsampling
- Multi-cluster federation
- OTLP ingestion
- Built-in Web UI

Use full Prometheus or pair with external tools (Thanos, Grafana).

## Configuration

### Minimal Configuration Example

```yaml
# prometheus.yml
global:
  scrape_interval: 15s
  evaluation_interval: 15s

rule_files:
  - "/etc/prometheus/rules/*.yml"

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

  # Self-monitoring
  - job_name: "prometheus"
    static_configs:
      - targets: ["localhost:9090"]
```

See `documentation/examples/prometheus-minimal.yml` for more examples.

### File-Based Service Discovery

Target file example (`/etc/prometheus/targets/app.json`):

```json
[
  {
    "targets": ["10.0.0.1:9100", "10.0.0.2:9100"],
    "labels": {
      "env": "production",
      "role": "web"
    }
  }
]
```

### Resource Constraints for Edge Deployments

```bash
./prometheus \
  --config.file=/etc/prometheus/prometheus.yml \
  --storage.tsdb.path=/data \
  --storage.tsdb.retention.time=6h \
  --storage.tsdb.retention.size=500MB \
  --storage.tsdb.wal-compression \
  --storage.tsdb.min-block-duration=30m \
  --query.max-concurrency=5 \
  --query.timeout=30s
```

| Parameter | Edge Value | Description |
|-----------|------------|-------------|
| `retention.time` | 6h-24h | Short-term local data only |
| `retention.size` | 200MB-1GB | Disk space limit |
| `wal-compression` | enabled | Reduce disk writes |
| `query.max-concurrency` | 5-10 | Limit concurrent queries |
| `query.timeout` | 30s | Prevent slow queries |

### Memory Estimation

| Active Series | Estimated Memory |
|---------------|------------------|
| 1,000 | ~50 MB |
| 10,000 | ~150 MB |
| 50,000 | ~500 MB |

Each active series uses ~4-8 KB (index + WAL buffer).

## API Limitations

The following endpoints return `501 Not Implemented` in minimal builds:
- `/api/v1/admin/tsdb/*` (remote read handlers)
- `/api/v1/write` (remote write)
- `/api/v1/otlp/*` (OTLP ingestion)

All query APIs work normally:
- ✅ `/api/v1/query`
- ✅ `/api/v1/query_range`
- ✅ `/api/v1/series`
- ✅ `/api/v1/labels`
- ✅ `/api/v1/label/:name/values`
- ✅ `/api/v1/alerts`
- ✅ `/api/v1/rules`
- ✅ `/api/v1/targets`
- ✅ `/-/healthy`, `/-/ready`

## Verification

### Functional Tests

```bash
# Start Prometheus
./prometheus-v3-upx --config.file=test.yml &
PID=$!
sleep 3

# Health check
curl -f http://localhost:9090/-/healthy

# Query test
curl -f "http://localhost:9090/api/v1/query?query=up"

# Targets check
curl -f http://localhost:9090/api/v1/targets

# Cleanup
kill $PID
```

### Size Verification

```bash
# Check no cloud SDK symbols remain
go tool nm prometheus-minimal | grep -E "aws-sdk|azure-sdk|k8s.io/client-go" && echo "FAIL" || echo "PASS"

# Verify UPX compression
upx -t prometheus-v3-upx && echo "Valid UPX binary" || echo "Not compressed"
```

## Upgrading from Full Build

1. **Remove unsupported config sections:**
   - Delete `remote_write` and `remote_read` blocks
   - Replace cloud SD with `static_configs` or `file_sd_configs`

2. **Update integrations:**
   - Use external Grafana for visualization
   - Configure long-term storage separately (if needed)

3. **Deploy minimal binary:**
   ```bash
   systemctl stop prometheus
   cp prometheus-v3-upx /usr/local/bin/prometheus
   systemctl start prometheus
   ```

4. **Verify:**
   - Check targets: `curl http://localhost:9090/api/v1/targets`
   - Check rules: `curl http://localhost:9090/api/v1/rules`

## Troubleshooting

### UPX Binary Won't Start

**Symptom:** Segmentation fault or "Exec format error"

**Solutions:**
- Try `upx --best` instead of `--ultra-brute`
- Check kernel version compatibility
- Fall back to uncompressed V3 static build
- Some container runtimes need `--cap-add=SYS_PTRACE`

### High Memory Usage

**Symptom:** Memory usage higher than expected

**Causes & Solutions:**
- Too many active series → Reduce scrape targets or retention
- WAL not compacting → Check `--storage.tsdb.min-block-duration`
- Query concurrency → Lower `--query.max-concurrency`

### Slow Startup with UPX

**Expected:** UPX adds 200-500ms decompression time

**Mitigation:**
- Use uncompressed build if startup time critical
- For containers, use init containers to decompress once

## Contributing

When modifying code affecting minimal builds:

1. **Add build tags to new files:**
   ```go
   //go:build !minimal

   package mypackage
   ```

2. **Create minimal stubs when needed:**
   ```go
   //go:build minimal

   package mypackage

   func ExpensiveFeature() error {
       return fmt.Errorf("not supported in minimal build")
   }
   ```

3. **Test both builds:**
   ```bash
   go build ./cmd/prometheus                              # Full
   go build -tags minimal,remove_all_sd ./cmd/prometheus # Minimal
   ```

4. **Verify size regression:**
   ```bash
   ./scripts/build-minimal-v3.sh
   # Check that V3+UPX stays ≤20 MB
   ```

## References

- [UPX - Ultimate Packer for eXecutables](https://upx.github.io/)
- [Go Build Constraints](https://pkg.go.dev/cmd/go#hdr-Build_constraints)
- [Prometheus Configuration](https://prometheus.io/docs/prometheus/latest/configuration/configuration/)
- [Can I have a smaller Prometheus?](https://wejick.wordpress.com/2022/01/29/can-i-have-a-smaller-prometheus/) - Community experience

## Version History

- **V1**: Research and planning
- **V2**: Minimal build with SD/remote storage removed (43 MB, 72% reduction)
- **V3**: Static compilation + UPX compression (15 MB, 90% reduction) ✅ Current
