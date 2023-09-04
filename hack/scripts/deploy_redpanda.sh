#!/bin/sh

set -e

if [ $# -ne 1 ]; then
    echo "project root is expected"
fi

PROJECT_ROOT="$1"

helm upgrade --install redpanda redpanda \
  --repo https://charts.redpanda.com \
  --namespace redpanda-system \
  --create-namespace \
  --valkue "${PROJECT_ROOT}/etc/examples/dapr/config/values_redpanda.yaml" \
  --wait
