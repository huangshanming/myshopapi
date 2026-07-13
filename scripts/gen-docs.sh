#!/usr/bin/env bash
# 从 Go handler 注释生成 OpenAPI 文档
# 输出: docs/openapi/mymall.yaml（合并三微服务，网关视角）
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
OUT_DIR="$ROOT/docs/openapi"
TMP_DIR="$ROOT/docs/.tmp-swag"
export PATH="$PATH:$(go env GOPATH)/bin"

green() { printf '\033[32m%s\033[0m\n' "$*"; }
yellow() { printf '\033[33m%s\033[0m\n' "$*"; }

if ! command -v swag >/dev/null 2>&1; then
  yellow "安装 swag..."
  go install github.com/swaggo/swag/cmd/swag@v1.16.4
fi

rm -rf "$TMP_DIR"
mkdir -p "$TMP_DIR" "$OUT_DIR"

run_swag() {
  local svc="$1"
  yellow "==> swag init: $svc"
  cd "$ROOT/services/$svc"
  swag init \
    -g cmd/main.go \
    -o docs \
    --dir . \
    --parseDependency \
    --parseInternal \
    --outputTypes json
  cp docs/swagger.json "$TMP_DIR/$svc.json"
}

run_swag user-service
run_swag catalog-service
run_swag order-service

yellow "==> 合并三份 swagger.json"
python3 "$ROOT/scripts/merge_openapi.py" \
  "$TMP_DIR/mymall.swagger.json" \
  "$TMP_DIR/user-service.json" \
  "$TMP_DIR/catalog-service.json" \
  "$TMP_DIR/order-service.json"

yellow "==> 转换为 OpenAPI 3 YAML"
python3 "$ROOT/scripts/swagger2yaml.py" \
  "$TMP_DIR/mymall.swagger.json" \
  "$OUT_DIR/mymall.yaml"

cp "$TMP_DIR/mymall.swagger.json" "$OUT_DIR/mymall.swagger.json"

green "文档已生成:"
echo "  $OUT_DIR/mymall.yaml          ← 复制给 AI / 导入 Apifox"
echo "  $OUT_DIR/mymall.swagger.json  ← Swagger 2.0 备份"
echo ""
echo "浏览文档:"
echo "  bash scripts/serve-docs.sh    # 推荐，http://localhost:9099/scalar/index.html"
