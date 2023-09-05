#!/bin/sh

set -e

helm upgrade --install redpanda redpanda \
  --repo https://charts.redpanda.com \
  --namespace redpanda-system \
  --create-namespace \
  --set external.enabled=false \
  --set statefulset.initContainers.setDataDirOwnership.enabled=true \
  --set statefulset.replicas=1 \
  --set tls.enabled=false \
  --set console.enabled=false

kubectl wait \
  --namespace redpanda-system \
  --for=condition=ready pod \
  --selector=app.kubernetes.io/component=redpanda-statefulset \
  --timeout=120s

kubectl -n redpanda-system exec \
  -ti redpanda-0 \
  -c redpanda -- \
    rpk topic create --brokers redpanda-0.redpanda.redpanda-system.svc.cluster.local.:9093 near

kubectl -n redpanda-system exec \
  -ti redpanda-0 \
  -c redpanda -- \
    rpk topic create --brokers redpanda-0.redpanda.redpanda-system.svc.cluster.local.:9093 far