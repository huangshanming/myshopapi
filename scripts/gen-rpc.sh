#!/usr/bin/env bash
# Generate pb + goctl zrpc client/server stubs for user/catalog/merchant.
# Same process still boots rest + zrpc via cmd/main.go (no separate rpc binary).
#
# Usage: ./scripts/gen-rpc.sh [user|catalog|merchant|all]
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
export PATH="$PATH:$(go env GOPATH)/bin"
HOME_TPL="${GOCTL_HOME:-$ROOT/deploy/goctl-template}"
TARGET="${1:-all}"

if ! command -v goctl >/dev/null 2>&1; then
  echo "goctl not found; install: go install github.com/zeromicro/go-zero/tools/goctl@v1.8.5" >&2
  exit 1
fi

# goctl derives a broken qualifier (e.g. v1_userv1) from path .../user/v1; rewrite to pb package name.
fix_pb_alias() {
  local dir="$1"
  local wrong="$2"
  local right="$3"
  local import_path="$4"
  while IFS= read -r -d '' f; do
    # ensure aliased import
    if grep -q "\"${import_path}\"" "$f" && ! grep -q "${right} \"${import_path}\"" "$f"; then
      perl -i -pe "s|\"${import_path}\"|${right} \"${import_path}\"|g" "$f"
    fi
    perl -i -pe "s/${wrong}/${right}/g" "$f"
  done < <(find "$dir" -type f -name '*.go' -print0)
}

gen_one() {
  local name="$1"          # user | catalog | merchant
  local proto_rel="user/v1/${name}.proto"
  case "$name" in
    user) proto_rel="user/v1/user.proto" ;;
    catalog) proto_rel="catalog/v1/catalog.proto" ;;
    merchant) proto_rel="merchant/v1/merchant.proto" ;;
  esac
  local svc_dir="$ROOT/services/${name}-service"
  local staging_rel=".tmp/rpc-${name}"
  local staging="$ROOT/$staging_rel"
  local pb_pkg="${name}v1"
  local wrong_alias="v1_${pb_pkg}"
  local import_path="mymall/api/gen/${name}/v1"
  local client_pkg="${name}service"
  # goctl names client dir from service: UserService -> userservice
  case "$name" in
    user) client_pkg="userservice" ;;
    catalog) client_pkg="catalogservice" ;;
    merchant) client_pkg="merchantservice" ;;
  esac

  if [[ ! -f "$ROOT/api/proto/$proto_rel" ]]; then
    echo "missing proto: api/proto/$proto_rel" >&2
    exit 1
  fi

  echo "==> goctl rpc: $name"
  rm -rf "$staging"
  mkdir -p "$staging"

  (
    cd "$ROOT"
    # Invoke from api/proto so paths stay consistent (goctl absolutizes args).
    cd api/proto
    goctl rpc protoc "$proto_rel" \
      --proto_path=. \
      --go_out="$ROOT/api/gen" \
      --go-grpc_out="$ROOT/api/gen" \
      --go_opt=paths=source_relative \
      --go-grpc_opt=paths=source_relative \
      --zrpc_out="$ROOT/$staging_rel" \
      --style go_zero \
      --home "$HOME_TPL"
  )

  fix_pb_alias "$staging" "$wrong_alias" "$pb_pkg" "$import_path"

  # pb already written under api/gen/<name>/v1 via go_package + source_relative

  # Shared goctl clients → api/rpcclient/<pkg>
  local client_src
  client_src="$(find "$staging" -type d -name "$client_pkg" | head -1)"
  if [[ -z "$client_src" ]]; then
    echo "client dir $client_pkg not found under $staging" >&2
    exit 1
  fi
  local client_dst="$ROOT/api/rpcclient/$client_pkg"
  mkdir -p "$client_dst"
  rm -f "$client_dst"/*.go
  cp "$client_src"/*.go "$client_dst/"
  # rewrite any residual staging imports (should already be api/gen)
  find "$client_dst" -name '*.go' -print0 | xargs -0 perl -i -pe "s|mymall/\\.tmp/rpc-${name}/[^\"]+|${import_path}|g"

  # Server stub → services/<svc>/internal/server/*_service_server.go (DO NOT EDIT header kept)
  local server_src
  server_src="$(find "$staging/internal/server" -name '*_server.go' | head -1)"
  mkdir -p "$svc_dir/internal/server"
  local server_base
  server_base="$(basename "$server_src")"
  local server_dst="$svc_dir/internal/server/$server_base"
  cp "$server_src" "$server_dst"
  perl -i -pe "s|mymall/\\.tmp/rpc-${name}/internal/logic|mymall/services/${name}-service/internal/rpclogic|g" "$server_dst"
  perl -i -pe "s|mymall/\\.tmp/rpc-${name}/internal/svc|mymall/services/${name}-service/internal/svc|g" "$server_dst"
  # package logic → rpclogic
  perl -i -pe 's|"mymall/services/'"${name}"'-service/internal/rpclogic"|rpclogic "mymall/services/'"${name}"'-service/internal/rpclogic"|g' "$server_dst"
  perl -i -pe 's/\blogic\./rpclogic./g' "$server_dst"

  # RPC logic stubs → internal/rpclogic (do not overwrite existing filled files unless FORCE_RPC_LOGIC=1)
  mkdir -p "$svc_dir/internal/rpclogic"
  for f in "$staging"/internal/logic/*.go; do
    [[ -f "$f" ]] || continue
    local base dst
    base="$(basename "$f")"
    dst="$svc_dir/internal/rpclogic/$base"
    if [[ -f "$dst" && "${FORCE_RPC_LOGIC:-0}" != "1" ]]; then
      echo "    keep $dst"
      continue
    fi
    cp "$f" "$dst"
    # package logic → rpclogic; fix imports to service svc + pb
    perl -i -pe 's/^package logic$/package rpclogic/' "$dst"
    perl -i -pe "s|mymall/\\.tmp/rpc-${name}/internal/svc|mymall/services/${name}-service/internal/svc|g" "$dst"
    perl -i -pe "s|mymall/\\.tmp/rpc-${name}/[^\"]+|${import_path}|g" "$dst"
    fix_pb_alias "$(dirname "$dst")" "$wrong_alias" "$pb_pkg" "$import_path"
    echo "    wrote $dst"
  done

  # Drop goctl main/etc/config — HTTP service owns those
  echo "    client -> api/rpcclient/$client_pkg"
  echo "    server -> $server_dst"
}

case "$TARGET" in
  user|catalog|merchant) gen_one "$TARGET" ;;
  all)
    gen_one user
    gen_one catalog
    gen_one merchant
    ;;
  *)
    echo "usage: $0 [user|catalog|merchant|all]" >&2
    exit 2
    ;;
esac

echo "goctl rpc done."
