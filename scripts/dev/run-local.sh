#!/usr/bin/env bash
set -euo pipefail

API_CONFIG="${API_CONFIG:-configs/api-server.yaml}"
SCHEDULER_CONFIG="${SCHEDULER_CONFIG:-configs/scheduler.yaml}"

cleanup() {
  jobs -p | xargs -r kill >/dev/null 2>&1 || true
}
trap cleanup EXIT

echo "[1/2] starting api-server..."
go run ./cmd/api-server -config "${API_CONFIG}" &
API_PID=$!

sleep 2

echo "[2/2] starting scheduler..."
go run ./cmd/scheduler -config "${SCHEDULER_CONFIG}" &
SCHEDULER_PID=$!

echo "api-server pid=${API_PID}"
echo "scheduler  pid=${SCHEDULER_PID}"
echo "Press Ctrl+C to stop."

wait