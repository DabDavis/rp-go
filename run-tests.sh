#!/usr/bin/env bash
set -euo pipefail

if ! command -v xvfb-run >/dev/null 2>&1; then
  cat <<'MSG' >&2
xvfb-run is required to execute the Go test suite in a headless environment.
Install the xvfb package (e.g. `sudo apt-get install xvfb`) and re-run this script.
MSG
  exit 1
fi

cmd=("go" "test" "./...")
if [[ $# -gt 0 ]]; then
  cmd=("go" "test" "$@")
fi

exec xvfb-run -a "${cmd[@]}"
