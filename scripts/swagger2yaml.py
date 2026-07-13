#!/usr/bin/env python3
"""将 Swagger 2.0 JSON 转为 OpenAPI 3 YAML。"""
from __future__ import annotations

import json
import sys
from pathlib import Path


def ref_name(ref: str) -> str:
    return ref.split("/")[-1]


def convert_schema(defn: dict, definitions: dict) -> dict:
    if not defn:
        return {}
    if "$ref" in defn:
        return {"$ref": f"#/components/schemas/{ref_name(defn['$ref'])}"}
    if "allOf" in defn:
        return {"allOf": [convert_schema(item, definitions) for item in defn["allOf"]]}
    out: dict = {}
    t = defn.get("type")
    if t:
        out["type"] = t
    for key in ("format", "example", "description", "enum"):
        if key in defn:
            out[key] = defn[key]
    if t == "array" and "items" in defn:
        out["items"] = convert_schema(defn["items"], definitions)
    if t == "object" or "properties" in defn:
        out["type"] = "object"
        props = {}
        for k, v in defn.get("properties", {}).items():
            props[k] = convert_schema(v, definitions)
        out["properties"] = props
        if "required" in defn:
            out["required"] = defn["required"]
    return out


def convert_parameters(params: list, definitions: dict) -> tuple[list, dict | None]:
    out_params = []
    request_body = None
    for p in params:
        if p.get("in") == "body":
            request_body = {
                "required": p.get("required", True),
                "description": p.get("description", ""),
                "content": {
                    "application/json": {
                        "schema": convert_schema(p.get("schema", {}), definitions)
                    }
                },
            }
            continue
        param = {
            "name": p["name"],
            "in": p["in"],
            "required": p.get("required", False),
            "description": p.get("description", ""),
        }
        if "schema" in p:
            param["schema"] = convert_schema(p["schema"], definitions)
        elif "type" in p:
            param["schema"] = {"type": p["type"]}
        out_params.append(param)
    return out_params, request_body


def convert(spec: dict) -> dict:
    definitions = spec.get("definitions", {})
    components_schemas = {
        name: convert_schema(defn, definitions) for name, defn in definitions.items()
    }

    paths: dict = {}
    for path, methods in spec.get("paths", {}).items():
        paths[path] = {}
        for method, op in methods.items():
            if method.startswith("x-"):
                continue
            operation = {
                "summary": op.get("summary", ""),
                "description": op.get("description", ""),
                "tags": op.get("tags", []),
                "responses": {},
            }
            if op.get("security"):
                operation["security"] = [{"BearerAuth": []}]
            for code, resp in op.get("responses", {}).items():
                operation["responses"][code] = {
                    "description": resp.get("description", "response"),
                }
                if "schema" in resp:
                    operation["responses"][code]["content"] = {
                        "application/json": {
                            "schema": convert_schema(resp["schema"], definitions)
                        }
                    }
            if "parameters" in op:
                params, body = convert_parameters(op["parameters"], definitions)
                if params:
                    operation["parameters"] = params
                if body:
                    operation["requestBody"] = body
            paths[path][method] = operation

    host = spec.get("host", "localhost")
    schemes = spec.get("schemes", ["http"])
    base = spec.get("basePath", "/")
    url = f"{schemes[0]}://{host}{base.rstrip('/')}"

    return {
        "openapi": "3.0.3",
        "info": spec.get("info", {"title": "API", "version": "1.0.0"}),
        "servers": [{"url": url, "description": "APISIX 网关 / 本地"}],
        "tags": spec.get("tags", []),
        "paths": paths,
        "components": {
            "schemas": components_schemas,
            "securitySchemes": {
                "BearerAuth": {
                    "type": "http",
                    "scheme": "bearer",
                    "bearerFormat": "JWT",
                    "description": "JWT Token，Header: Authorization: Bearer {token}",
                }
            },
        },
    }


def dump_yaml(data: dict) -> str:
    try:
        import yaml  # type: ignore

        return yaml.dump(data, allow_unicode=True, sort_keys=False, default_flow_style=False)
    except ImportError:
        return json.dumps(data, ensure_ascii=False, indent=2)


def main() -> int:
    if len(sys.argv) != 3:
        print("usage: swagger2yaml.py <in.json> <out.yaml>", file=sys.stderr)
        return 1
    src = Path(sys.argv[1])
    dst = Path(sys.argv[2])
    spec = json.loads(src.read_text(encoding="utf-8"))
    dst.parent.mkdir(parents=True, exist_ok=True)
    dst.write_text(dump_yaml(convert(spec)), encoding="utf-8")
    print(f"openapi3 -> {dst}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
