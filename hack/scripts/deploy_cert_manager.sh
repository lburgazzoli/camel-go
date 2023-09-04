#!/bin/sh

set -e

VERSION="v1.12.4"


kubectl apply -f "https://github.com/cert-manager/cert-manager/releases/download/${VERSION}/cert-manager.crds.yaml"


helm repo add --force-update jetstack https://charts.jetstack.io
helm repo update

helm install cert-manager jetstack/cert-manager \
  --namespace cert-manager \
  --create-namespace \
  --version "${VERSION}"