#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
DIST_DIR="${ROOT_DIR}/dist"

mkdir -p "${DIST_DIR}"

build() {
  local goos="$1"
  local goarch="$2"
  local output="$3"

  echo "Building ${output}"
  env GOCACHE="${TMPDIR:-/tmp}/openclaw-gocache" \
    GOMODCACHE="${TMPDIR:-/tmp}/openclaw-gomodcache" \
    GOOS="${goos}" GOARCH="${goarch}" \
    go build -o "${DIST_DIR}/${output}" ./cmd/openclaw-installer
}

cd "${ROOT_DIR}"

build darwin arm64 openclaw-installer-darwin-arm64
build darwin amd64 openclaw-installer-darwin-amd64
build linux amd64 openclaw-installer-linux-amd64
build windows amd64 openclaw-installer-windows-amd64.exe

echo "Artifacts written to ${DIST_DIR}"
