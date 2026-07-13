#!/usr/bin/env python3
"""合并三微服务 swagger.json 为一份网关视角文档。"""
from __future__ import annotations

import json
import sys
from pathlib import Path
from copy import deepcopy


def load(path: Path) -> dict:
    with path.open(encoding="utf-8") as f:
        return json.load(f)


def merge_specs(paths: list[Path]) -> dict:
    merged: dict = {}
    for i, p in enumerate(paths):
        spec = load(p)
        if i == 0:
            merged = deepcopy(spec)
            merged.setdefault("paths", {})
            merged.setdefault("definitions", {})
            merged.setdefault("tags", [])
            continue
        merged["paths"].update(spec.get("paths", {}))
        merged["definitions"].update(spec.get("definitions", {}))
        existing = {t.get("name") for t in merged.get("tags", [])}
        for tag in spec.get("tags", []):
            if tag.get("name") not in existing:
                merged.setdefault("tags", []).append(tag)
        if "securityDefinitions" in spec:
            merged.setdefault("securityDefinitions", {}).update(spec["securityDefinitions"])

    merged["info"] = {
        "title": "mymall 商城 API",
        "description": (
            "宠物商城 REST API（user + catalog + order 合并文档）。\n\n"
            "统一响应 `{code, msg, data}`，code=200 表示业务成功。\n"
            "经 APISIX 网关: http://localhost:9080\n"
            "本地直连: user :8881, catalog :8882, order :8883"
        ),
        "version": "1.0.0",
        "contact": {"name": "mymall"},
    }
    merged["host"] = "localhost:9080"
    merged["basePath"] = "/"
    merged["schemes"] = ["http"]
    merged["tags"] = [
        {"name": "用户", "description": "注册、登录、个人资料"},
        {"name": "商品", "description": "商品列表与详情"},
        {"name": "分类", "description": "商品分类"},
        {"name": "订单", "description": "下单、查单、取消"},
    ]
    return merged


def main() -> int:
    if len(sys.argv) < 3:
        print("usage: merge_openapi.py <out.json> <spec1.json> [spec2.json ...]", file=sys.stderr)
        return 1
    out = Path(sys.argv[1])
    inputs = [Path(p) for p in sys.argv[2:]]
    for p in inputs:
        if not p.exists():
            print(f"missing: {p}", file=sys.stderr)
            return 1
    merged = merge_specs(inputs)
    out.parent.mkdir(parents=True, exist_ok=True)
    with out.open("w", encoding="utf-8") as f:
        json.dump(merged, f, ensure_ascii=False, indent=2)
    print(f"merged -> {out} ({len(merged.get('paths', {}))} paths)")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
