#!/usr/bin/env bash
set -euo pipefail

cmd=("go" "test" "./...")
if [[ $# -gt 0 ]]; then
  cmd=("go" "test" "$@")
fi

if command -v xvfb-run >/dev/null 2>&1; then
  exec xvfb-run -a "${cmd[@]}"
fi

if ! command -v Xvfb >/dev/null 2>&1; then
  cat <<'MSG' >&2
Neither `xvfb-run` nor a raw `Xvfb` binary are available to provide a virtual framebuffer.
Install the `xvfb` package (e.g. `sudo apt-get install xvfb`) and re-run this script.
MSG
  exit 1
fi

XVFB_DISPLAY="${XVFB_DISPLAY:-:99}"
XVFB_GEOMETRY="${XVFB_GEOMETRY:-1024x768x24}"

cleanup() {
  if [[ -n "${XVFB_PID:-}" ]]; then
    kill "${XVFB_PID}" 2>/dev/null || true
    wait "${XVFB_PID}" 2>/dev/null || true
  fi
}
trap cleanup EXIT

Xvfb "${XVFB_DISPLAY}" -screen 0 "${XVFB_GEOMETRY}" &
XVFB_PID=$!
export DISPLAY="${XVFB_DISPLAY}"

"${cmd[@]}"
