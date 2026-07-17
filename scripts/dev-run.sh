#!/usr/bin/env bash
# 兼容旧入口：转发到 scripts/dev.sh（air 热更新）
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
exec bash "$ROOT/scripts/dev.sh" "${1:-all}"
