#!/usr/bin/env bash
# 下载首页同款门头图到 uploads/shops/seed/
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
DIR="$ROOT/uploads/shops/seed"
mkdir -p "$DIR"
cd "$DIR"
curl -fsSL -o 1054.jpg "https://picsum.photos/id/1054/640/360"
curl -fsSL -o 1080.jpg "https://picsum.photos/id/1080/640/360"
curl -fsSL -o 1066.jpg "https://picsum.photos/id/1066/640/360"
echo "saved to $DIR"
ls -la
