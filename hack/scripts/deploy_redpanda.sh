#!/bin/sh

set -e

helm upgrade --install redpanda redpanda \
  --repo https://charts.redpanda.com \
  --namespace redpanda-system \
  --create-namespace \
  --set nameOverride="redpanda" \
  --set console.enabled=false \
  --set connectors.enabled=false \
  --set external.enabled=false \
  --set monitoring.enabled=false \
  --set statefulset.replicas=1 \
  --wait
