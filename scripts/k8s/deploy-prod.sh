#!/usr/bin/env bash
set -euo pipefail

NAMESPACE="${NAMESPACE:-gpu-system}"
RELEASE_NAME="${RELEASE_NAME:-gpu-scheduler-platform}"

helm upgrade --install "${RELEASE_NAME}" \
  ./deployments/helm/gpu-scheduler-platform \
  --namespace "${NAMESPACE}" \
  --create-namespace \
  -f ./deployments/helm/gpu-scheduler-platform/values.yaml \
  -f ./deployments/helm/gpu-scheduler-platform/values-prod.yaml