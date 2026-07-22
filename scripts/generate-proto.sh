#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="$PATH:$(go env GOPATH)/bin"

mkdir -p "$ROOT/api/gen/user/v1" "$ROOT/api/gen/catalog/v1" "$ROOT/api/gen/merchant/v1"

protoc \
  --proto_path="$ROOT/api/proto" \
  --go_out="$ROOT/api/gen" --go_opt=paths=source_relative \
  --go-grpc_out="$ROOT/api/gen" --go-grpc_opt=paths=source_relative \
  "$ROOT/api/proto/user/v1/user.proto" \
  "$ROOT/api/proto/catalog/v1/catalog.proto" \
  "$ROOT/api/proto/merchant/v1/merchant.proto"

echo "proto 生成完成 -> api/gen/"
