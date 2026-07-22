#!/usr/bin/env bash
# Deprecated alias: prefer ./scripts/gen-rpc.sh
# Still regenerates pb + goctl zrpc stubs for user/catalog/merchant.
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
exec "$ROOT/scripts/gen-rpc.sh" "${1:-all}"
