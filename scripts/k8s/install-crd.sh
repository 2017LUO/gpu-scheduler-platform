#!/usr/bin/env bash
set -euo pipefail

kubectl apply -f api/crd/gpujobs.yaml
kubectl apply -f api/crd/gpuqueues.yaml
kubectl apply -f api/crd/gpuquotas.yaml
kubectl apply -f api/crd/gpupolicies.yaml
kubectl apply -f api/crd/gpuclustersnapshots.yaml

echo "CRDs installed."