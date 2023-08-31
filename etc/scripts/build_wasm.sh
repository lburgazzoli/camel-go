#!/bin/sh

if [ $# -ne 3 ]; then
    echo "project root, in and out are expected"
fi

PROJECT_ROOT="$1"
IN="$2"
OUT="$3"

docker run \
		--rm \
		-v "${PROJECT_ROOT}":/src:Z \
		-w /src \
		tinygo/tinygo:"${TINYGO_VERSION}" \
		tinygo build \
			-target=wasi \
			-scheduler=none \
			-o "${OUT}" \
			"${IN}"