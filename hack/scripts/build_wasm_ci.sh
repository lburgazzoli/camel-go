#!/bin/sh

if [ $# -ne 3 ]; then
    echo "project root, in and out are expected"
fi

PROJECT_ROOT="$1"
IN="$2"
OUT="$3"

tinygo build \
	-target=wasi \
	-scheduler=none \
	-o "${OUT}" \
	"${IN}"