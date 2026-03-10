#!/usr/bin/env bash
set -euo pipefail

OUT_DIR="${OUT_DIR:-bin}"
mkdir -p "${OUT_DIR}"

echo "Building binaries into ${OUT_DIR} ..."

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/api-server" ./cmd/api-server
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/scheduler" ./cmd/scheduler
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/controller" ./cmd/controller
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/webhook" ./cmd/webhook
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "${OUT_DIR}/agent" ./cmd/agent

echo "Build completed."