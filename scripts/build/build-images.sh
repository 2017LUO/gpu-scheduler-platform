#!/usr/bin/env bash
set -euo pipefail

IMAGE_REGISTRY="${IMAGE_REGISTRY:-local}"
IMAGE_TAG="${IMAGE_TAG:-dev}"

docker build -f Dockerfile.api -t "${IMAGE_REGISTRY}/gpu-api-server:${IMAGE_TAG}" .
docker build -f Dockerfile.scheduler -t "${IMAGE_REGISTRY}/gpu-scheduler:${IMAGE_TAG}" .
docker build -f Dockerfile.controller -t "${IMAGE_REGISTRY}/gpu-controller:${IMAGE_TAG}" .
docker build -f Dockerfile.webhook -t "${IMAGE_REGISTRY}/gpu-webhook:${IMAGE_TAG}" .
docker build -f Dockerfile.agent -t "${IMAGE_REGISTRY}/gpu-agent:${IMAGE_TAG}" .

echo "Images built:"
echo "  ${IMAGE_REGISTRY}/gpu-api-server:${IMAGE_TAG}"
echo "  ${IMAGE_REGISTRY}/gpu-scheduler:${IMAGE_TAG}"
echo "  ${IMAGE_REGISTRY}/gpu-controller:${IMAGE_TAG}"
echo "  ${IMAGE_REGISTRY}/gpu-webhook:${IMAGE_TAG}"
echo "  ${IMAGE_REGISTRY}/gpu-agent:${IMAGE_TAG}"