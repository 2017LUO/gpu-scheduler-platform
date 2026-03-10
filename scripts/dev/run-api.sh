#!/usr/bin/env bash
set -euo pipefail

CONFIG="${1:-configs/api-server.yaml}"

echo "Starting api-server with config: ${CONFIG}"
go run ./cmd/api-server -config "${CONFIG}"