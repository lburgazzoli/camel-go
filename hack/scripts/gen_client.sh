#!/bin/sh

if [ $# -ne 1 ]; then
    echo "project root is expected"
fi

PROJECT_ROOT="$1"
TMP_DIR=$( mktemp -d -t camel-go-client-gen-XXXXXXXX )

mkdir -p "${TMP_DIR}/client"
mkdir -p "${PROJECT_ROOT}/pkg/client"

echo "tmp dir: $TMP_DIR"

echo "applyconfiguration-gen"
"${PROJECT_ROOT}"/bin/applyconfiguration-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/applyconfiguration" \
  --output-pkg=github.com/lburgazzoli/camel-go/pkg/client/applyconfiguration \
  github.com/lburgazzoli/camel-go/api/camel/v2alpha1

echo "client-gen"
"${PROJECT_ROOT}"/bin/client-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/clientset" \
  --input-base=github.com/lburgazzoli/camel-go/api \
  --input=camel/v2alpha1 \
  --fake-clientset=false \
  --clientset-name "versioned"  \
  --apply-configuration-package=github.com/lburgazzoli/camel-go/pkg/client/applyconfiguration \
  --output-pkg=github.com/lburgazzoli/camel-go/pkg/client/clientset

echo "lister-gen"
"${PROJECT_ROOT}"/bin/lister-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/listers" \
  --output-pkg=github.com/lburgazzoli/camel-go/pkg/client/listers \
  github.com/lburgazzoli/camel-go/api/camel/v2alpha1

echo "informer-gen"
"${PROJECT_ROOT}"/bin/informer-gen \
  --go-header-file="${PROJECT_ROOT}/hack/boilerplate.go.txt" \
  --output-dir="${TMP_DIR}/client/informers" \
  --versioned-clientset-package=github.com/lburgazzoli/camel-go/pkg/client/clientset/versioned \
  --listers-package=github.com/lburgazzoli/camel-go/pkg/client/listers \
  --output-pkg=github.com/lburgazzoli/camel-go/pkg/client/informers \
  github.com/lburgazzoli/camel-go/api/camel/v2alpha1


## This should not be needed but for some reasons, the applyconfiguration-gen tool
## sets a wrong APIVersion for the Dapr type (operator/v2alpha1 instead of the one with
## the domain operator.dapr.io/v2alpha1).
##
## See: https://github.com/kubernetes/code-generator/issues/150
sed -i \
  's/WithAPIVersion(\"camel\/v2alpha1\")/WithAPIVersion(\"camel.apache.org\/v2alpha1\")/g' \
  "${TMP_DIR}"/client/applyconfiguration/camel/v2alpha1/integration.go

cp -r \
  "${TMP_DIR}"/client/* \
  "${PROJECT_ROOT}"/pkg/client

