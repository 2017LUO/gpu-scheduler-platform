#!/usr/bin/env bash
set -euo pipefail

CONFIG="${1:-configs/scheduler.yaml}"

echo "Starting scheduler with config: ${CONFIG}"
go run ./cmd/scheduler -config "${CONFIG}"