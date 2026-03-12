#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# shellcheck disable=SC1091
OPENCLAW_INSTALLER_INSTALL_LATEST_SOURCED=1 source "${ROOT_DIR}/scripts/install-latest.sh"

assert_eq() {
  local expected="$1"
  local actual="$2"

  if [[ "${expected}" != "${actual}" ]]; then
    echo "expected: ${expected}" >&2
    echo "actual:   ${actual}" >&2
    exit 1
  fi
}

assert_fails() {
  if "$@" >/dev/null 2>&1; then
    echo "expected command to fail: $*" >&2
    exit 1
  fi
}

assert_eq "openclaw-installer-darwin-arm64" "$(resolve_asset_name Darwin arm64)"
assert_eq "openclaw-installer-darwin-amd64" "$(resolve_asset_name Darwin x86_64)"
assert_eq "openclaw-installer-linux-amd64" "$(resolve_asset_name Linux x86_64)"
assert_eq \
  "https://github.com/liuyingwen/openclaw-installer/releases/latest/download/openclaw-installer-linux-amd64" \
  "$(latest_release_url openclaw-installer-linux-amd64)"

assert_fails resolve_asset_name Linux aarch64
assert_fails resolve_asset_name FreeBSD x86_64

echo "install-latest.sh tests passed"
