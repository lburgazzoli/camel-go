#!/bin/sh

set -e

DAPR_VERSION="1.11.0"

helm upgrade --install dapr dapr \
    --repo https://dapr.github.io/helm-charts \
		--version="${DAPR_VERSION}" \
		--create-namespace \
		--namespace dapr-system \
		--wait