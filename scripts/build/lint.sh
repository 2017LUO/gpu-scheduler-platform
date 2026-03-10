#!/usr/bin/env bash
set -euo pipefail

go fmt ./...
go vet ./...

if command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run ./...
else
  echo "golangci-lint not found, skip."
fi