#!/usr/bin/env bash
set -euo pipefail

# Usage:
# ./scripts/build.sh -> build for current platform, output: bin/<name>
# GOOS=linux GOARCH=amd64 ./scripts/build.sh -> cross-build

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(cd "$script_dir/.." && pwd)"
OUTDIR="${OUTDIR:-$repo_root/bin}"
NAME="${NAME:-stockportfolio-tui}"
LDFLAGS="${LDFLAGS:-}"

mkdir -p "$OUTDIR"

GOOS=${GOOS:-$(go env GOOS)}
GOARCH=${GOARCH:-$(go env GOARCH)}

echo "Building $NAME for ${GOOS}/${GOARCH}..."
# static binary where possible. We need CGO for some parts (e.g., sqlite), so we can't fully disable it.
export CGO_ENABLED=1

out="$OUTDIR/${NAME}-${GOOS}-${GOARCH}"

cd "$repo_root/src"
go build -trimpath -ldflags "$LDFLAGS" -o "$out" ./cmd/tui

# make executable and report size
chmod +x "$out"
echo "Built: $out ($(stat -c '%s' "$out") bytes)"