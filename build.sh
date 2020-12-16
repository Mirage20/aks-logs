#!/usr/bin/env bash

GIT_REVISION=$(git rev-parse --short --verify HEAD)
TIME=$(date -u +%Y%m%d.%H%M%S)
VERSION=1.0.${GIT_REVISION}.${TIME}

build_artifacts () {
  local os=$1
  local arch=$2
  GO111MODULE=on GOOS=$os GOARCH=$arch go build -ldflags "-X main.versionString=${VERSION}" ./cmd/aks-logs/
  file aks-logs
  tar -czvf aks-logs-"$os"-x64.tar.gz aks-logs
  rm aks-logs
}

build_artifacts linux amd64
build_artifacts darwin amd64
