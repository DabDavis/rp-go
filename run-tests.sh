#!/usr/bin/env bash
set -euo pipefail

args=("$@")
if [[ ${#args[@]} -eq 0 ]]; then
  args=("./...")
fi

go test -tags headless "${args[@]}"
