#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${NAMESPACE:-gpu-system}"
RELEASE_NAME="${RELEASE_NAME:-gpu-scheduler-platform}"

helm uninstall "${RELEASE_NAME}" -n "${NAMESPACE}" || true