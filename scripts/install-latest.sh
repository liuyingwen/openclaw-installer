#!/usr/bin/env bash
set -euo pipefail

OPENCLAW_INSTALLER_REPO="${OPENCLAW_INSTALLER_REPO:-liuyingwen/openclaw-installer}"
OPENCLAW_INSTALLER_TEMP_DIR=""

cleanup_temp_dir() {
  if [[ -n "${OPENCLAW_INSTALLER_TEMP_DIR:-}" ]]; then
    rm -rf -- "${OPENCLAW_INSTALLER_TEMP_DIR}"
  fi
}

resolve_asset_name() {
  local kernel="$1"
  local machine="$2"
  local goos
  local goarch

  case "${kernel}" in
    Darwin)
      goos="darwin"
      ;;
    Linux)
      goos="linux"
      ;;
    *)
      echo "unsupported operating system: ${kernel}" >&2
      return 1
      ;;
  esac

  case "${machine}" in
    arm64)
      if [[ "${goos}" != "darwin" ]]; then
        echo "unsupported architecture for ${goos}: ${machine}" >&2
        return 1
      fi
      goarch="arm64"
      ;;
    x86_64 | amd64)
      goarch="amd64"
      ;;
    *)
      echo "unsupported architecture: ${machine}" >&2
      return 1
      ;;
  esac

  printf 'openclaw-installer-%s-%s\n' "${goos}" "${goarch}"
}

latest_release_url() {
  local asset_name="$1"
  printf 'https://github.com/%s/releases/latest/download/%s\n' "${OPENCLAW_INSTALLER_REPO}" "${asset_name}"
}

download_file() {
  local url="$1"
  local output_path="$2"

  if command -v curl >/dev/null 2>&1; then
    curl -fsSL -o "${output_path}" "${url}"
    return 0
  fi

  if command -v wget >/dev/null 2>&1; then
    wget -qO "${output_path}" "${url}"
    return 0
  fi

  echo "missing downloader: install curl or wget first" >&2
  return 1
}

main() {
  local asset_name
  local url
  local binary_path

  asset_name="$(resolve_asset_name "$(uname -s)" "$(uname -m)")"
  url="$(latest_release_url "${asset_name}")"
  OPENCLAW_INSTALLER_TEMP_DIR="$(mktemp -d)"
  binary_path="${OPENCLAW_INSTALLER_TEMP_DIR}/openclaw-installer"

  trap cleanup_temp_dir EXIT

  echo "Downloading ${asset_name}..." >&2
  download_file "${url}" "${binary_path}"
  chmod +x "${binary_path}"

  if [[ "$#" -eq 0 ]]; then
    set -- install --yes
  fi

  "${binary_path}" "$@"
}

if [[ "${OPENCLAW_INSTALLER_INSTALL_LATEST_SOURCED:-0}" != "1" ]]; then
  main "$@"
fi
