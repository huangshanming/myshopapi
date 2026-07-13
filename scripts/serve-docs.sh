#!/usr/bin/env bash
# 本地启动文档服务（Scalar 必须通过 HTTP 访问，不能直接 open file://）
#
# 用法:
#   bash scripts/serve-docs.sh          # 启动 → http://localhost:9099/scalar/index.html
#   bash scripts/serve-docs.sh --stop   # 停止占用端口的旧服务
#   bash scripts/serve-docs.sh --scalar # Scalar CLI → http://localhost:9099/
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DOCS_DIR="$ROOT/docs"
SPEC="$DOCS_DIR/openapi/mymall.yaml"
PORT="${DOCS_PORT:-9099}"
URL="http://localhost:$PORT/scalar/index.html"

port_in_use() {
  lsof -i ":$PORT" -sTCP:LISTEN >/dev/null 2>&1
}

stop_server() {
  local pids
  pids=$(lsof -ti ":$PORT" -sTCP:LISTEN 2>/dev/null || true)
  if [[ -z "$pids" ]]; then
    echo "端口 $PORT 上没有运行中的文档服务"
    return 0
  fi
  echo "停止端口 $PORT 上的服务 (PID: $pids)..."
  kill $pids 2>/dev/null || true
  sleep 0.5
  echo "已停止"
}

if [[ "${1:-}" == "--stop" ]]; then
  stop_server
  exit 0
fi

if [[ ! -f "$SPEC" ]]; then
  echo "未找到 $SPEC，请先执行: bash scripts/gen-docs.sh"
  exit 1
fi

if port_in_use; then
  if curl -sf "$URL" >/dev/null 2>&1; then
    echo "文档服务已在运行，直接打开即可："
    echo "  $URL"
    echo ""
    echo "若要重启: bash scripts/serve-docs.sh --stop && bash scripts/serve-docs.sh"
    exit 0
  fi
  echo "端口 $PORT 已被其他程序占用，请："
  echo "  1. 换端口: DOCS_PORT=9100 bash scripts/serve-docs.sh"
  echo "  2. 或释放端口: bash scripts/serve-docs.sh --stop"
  exit 1
fi

MODE="${1:-python}"

if [[ "$MODE" == "--scalar" ]] && command -v npx >/dev/null 2>&1; then
  echo "Scalar CLI 文档服务（须保持此终端运行）"
  echo "  打开: http://localhost:$PORT/"
  echo "  按 Ctrl+C 停止"
  echo ""
  exec npx --yes @scalar/cli serve "$SPEC" --port "$PORT"
fi

if ! command -v python3 >/dev/null 2>&1; then
  echo "需要 Python3。或尝试: bash scripts/serve-docs.sh --scalar"
  exit 1
fi

echo "文档 HTTP 服务（须保持此终端运行，关终端即停）"
echo "  打开: $URL"
echo "  按 Ctrl+C 停止"
echo ""

cd "$DOCS_DIR"
exec python3 -m http.server "$PORT"
