#!/usr/bin/env bash

set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)

docker compose -f "$SCRIPT_DIR/compose.yaml" up --build --abort-on-container-exit

echo
echo "Logs are in: ${LOG_DIR:-$SCRIPT_DIR/out}"
