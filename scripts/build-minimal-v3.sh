#!/bin/bash
# Prometheus V3 Build Script
# Builds minimal, deep-minimal, and UPX-compressed variants

set -euo pipefail

VERSION=${1:-"v3.0.0-minimal"}
ARCH=${2:-"amd64"}

echo "=== Prometheus Minimal Core V3 Build Script ==="
echo "Version: $VERSION"
echo "Architecture: $ARCH"
echo ""

# Build flags
BASE_TAGS="minimal,remove_all_sd"
LDFLAGS="-s -w -X github.com/prometheus/common/version.Version=$VERSION"
TRIMPATH="-trimpath"

# Create output directory
mkdir -p dist

echo "=== Stage 1: Building V2 Minimal (Baseline) ==="
CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build \
  -tags "$BASE_TAGS" \
  -ldflags="$LDFLAGS" \
  $TRIMPATH \
  -o dist/prometheus-v2-minimal \
  ./cmd/prometheus

V2_SIZE=$(stat -c%s dist/prometheus-v2-minimal 2>/dev/null || stat -f%z dist/prometheus-v2-minimal)
V2_SIZE_MB=$((V2_SIZE / 1048576))
echo "✓ V2 Minimal built: ${V2_SIZE_MB} MB"
echo ""

echo "=== Stage 2: Building V3 Deep Minimal (Static + Optimized) ==="
# V3 uses stricter build settings but same tags (histogram/exemplar removal proved too complex)
CGO_ENABLED=0 GOOS=linux GOARCH=$ARCH go build \
  -tags "$BASE_TAGS" \
  -ldflags="$LDFLAGS" \
  $TRIMPATH \
  -buildmode=pie \
  -o dist/prometheus-v3-deep \
  ./cmd/prometheus

V3_SIZE=$(stat -c%s dist/prometheus-v3-deep 2>/dev/null || stat -f%z dist/prometheus-v3-deep)
V3_SIZE_MB=$((V3_SIZE / 1048576))
echo "✓ V3 Deep built: ${V3_SIZE_MB} MB"
echo ""

# Check if UPX is available
if command -v upx &> /dev/null; then
  echo "=== Stage 3: UPX Compression ==="

  cp dist/prometheus-v3-deep dist/prometheus-v3-upx
  upx --best dist/prometheus-v3-upx 2>&1 | grep -E "compressed|packed" || true

  UPX_SIZE=$(stat -c%s dist/prometheus-v3-upx 2>/dev/null || stat -f%z dist/prometheus-v3-upx)
  UPX_SIZE_MB=$((UPX_SIZE / 1048576))
  echo "✓ UPX compressed: ${UPX_SIZE_MB} MB"
  echo ""
else
  echo "⚠ UPX not found, skipping compression"
  echo "  Install: apt-get install upx (Linux) or brew install upx (macOS)"
  echo ""
fi

echo "=== Build Summary ==="
echo "V2 Minimal:    ${V2_SIZE_MB} MB  (baseline)"
echo "V3 Deep:       ${V3_SIZE_MB} MB  (static, optimized)"
if command -v upx &> /dev/null; then
  REDUCTION=$(( (V2_SIZE - UPX_SIZE) * 100 / V2_SIZE ))
  echo "V3 + UPX:      ${UPX_SIZE_MB} MB  (-${REDUCTION}% from V2)"
fi
echo ""

echo "=== Output Files ==="
ls -lh dist/prometheus-* 2>/dev/null || ls -l dist/prometheus-*
echo ""

echo "=== Verification ==="
echo "Testing V3 binary..."
if dist/prometheus-v3-deep --version > /dev/null 2>&1; then
  echo "✓ V3 binary runs successfully"
  dist/prometheus-v3-deep --version
else
  echo "✗ V3 binary failed to run"
  exit 1
fi

if command -v upx &> /dev/null && [ -f dist/prometheus-v3-upx ]; then
  echo ""
  echo "Testing UPX compressed binary..."
  if dist/prometheus-v3-upx --version > /dev/null 2>&1; then
    echo "✓ UPX binary runs successfully"
  else
    echo "✗ UPX binary failed to run"
  fi
fi

echo ""
echo "=== Build Complete ==="
