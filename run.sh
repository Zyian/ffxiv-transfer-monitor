#!/bin/bash

if ! [ -x "$(command -v go)" ]; then
  echo 'Go is not installed, please install with your package manager' >&2
  exit 1
fi

go run monitor.go
