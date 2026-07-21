#!/usr/bin/env python3
"""Generate invoke-based logic for a service after FORCE_REGEN.

Reads mapping file services/<svc>/.layering-map.json if present:
  { "HandlerName": { "import": "...", "call": "pkg.NewX(l.svcCtx).Method", "http": "GET", "path": "/..." } }

Or discovers from leftover comments. For automatic mode, scans internal/app for
Handler methods and matches by handler name == method name.
"""
from __future__ import annotations

import json
import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]


def load_map(svc: str) -> dict:
    p = ROOT / "services" / svc / ".layering-map.json"
    if p.exists():
        return json.loads(p.read_text())
    return {}


def discover_app_handlers(svc_dir: Path) -> dict[str, tuple[str, str, str]]:
    """methodName -> (import_path, ctor_expr_template, package_alias)"""
    out = {}
    app_roots = [svc_dir / "internal" / "app"]
    for nest in ("product", "content", "shopops", "notify"):
        app_roots.append(svc_dir / "internal" / nest / "app")

    for app in app_roots:
        if not app.is_dir():
            continue
        for go in app.rglob("*.go"):
            if go.name.startswith("invoke"):
                continue
            text = go.read_text()
            pkg = re.search(r"^package (\w+)", text, re.M)
            if not pkg:
                continue
            pkg_name = pkg.group(1)
            # relative import
            rel = go.parent.relative_to(svc_dir)
            imp = f"mymall/services/{svc_dir.name}/{rel.as_posix()}"
            alias = {
                "user": "huser",
                "admin": "hadmin",
                "public": "hpublic",
                "merchant": "hmerchant",
                "internalapi": "hinternal",
            }.get(pkg_name, f"h{pkg_name}")

            # NewXxxHandler
            news = re.findall(r"func New(\w+Handler)\(", text)
            methods = re.findall(
                r"func \(h \*(\w+)\) (\w+)\(w http\.ResponseWriter, r \*http\.Request\)",
                text,
            )
            # also plain funcs
            funcs = re.findall(
                r"^func (\w+)\(w http\.ResponseWriter, r \*http\.Request\)",
                text,
                re.M,
            )
            for recv, method in methods:
                ctor = None
                for n in news:
                    if n.replace("Handler", "") in recv or recv.startswith(n.replace("Handler", "")):
                        ctor = f"{alias}.New{n}(l.svcCtx).{method}"
                        break
                if ctor is None and news:
                    # pick first New that returns matching type
                    for n in news:
                        if n.replace("Handler", "") == recv.replace("Handler", ""):
                            ctor = f"{alias}.New{n}(l.svcCtx).{method}"
                            break
                if ctor is None:
                    for n in news:
                        ctor = f"{alias}.New{n}(l.svcCtx).{method}"
                        break
                if ctor:
                    out[method] = (imp, ctor, alias)

            for fn in funcs:
                out[fn] = (imp, f"{alias}.{fn}", alias)
    return out


def parse_routes(api_path: Path) -> list[tuple[str, str, str]]:
    """handler, method, path"""
    text = api_path.read_text()
    routes = []
    for m in re.finditer(
        r"@handler\s+(\w+)\s*\n\s*(get|post|put|delete|patch)\s+(\S+)",
        text,
        re.I,
    ):
        path = m.group(3)
        path = re.sub(r"\s*\(.*$", "", path)  # strip (Req) returns
        path = re.sub(r"\s*returns.*$", "", path)
        routes.append((m.group(1), m.group(2).upper(), path.strip()))
    return routes


def write_logic_file(path: Path, svc: str, handler: str, http: str, route_path: str, discovered: dict, explicit: dict):
    text = path.read_text()
    pkg = re.search(r"^package (\w+)", text, re.M).group(1)
    typ = re.search(r"type (\w+) struct", text).group(1)
    fn = re.search(rf"func \(l \*{typ}\) (\w+)\(", text).group(1)
    req_m = re.search(r"req \*types\.(\w+)", text)
    resp_m = re.search(r"resp \*types\.(\w+)", text)
    req_type = req_m.group(1) if req_m else None
    resp_type = resp_m.group(1) if resp_m else None

    # specials
    if fn in ("Healthz", "Metrics", "Readyz"):
        return write_health_special(path, pkg, typ, fn, svc, resp_type)

    meta = explicit.get(fn) or explicit.get(handler)
    if meta:
        ctor = meta["call"]
        imp = meta["import"]
        alias = meta.get("alias", "happ")
        http = meta.get("http", http)
        route_path = meta.get("path", route_path)
    else:
        # match by method name
        info = discovered.get(fn)
        if not info:
            # try without prefix
            for k, v in discovered.items():
                if k.lower() == fn.lower() or fn.endswith(k) or k.endswith(fn):
                    info = v
                    break
        if not info:
            print("SKIP no app method for", fn, path)
            # leave todo stub that returns error
            body = 'return nil, fmt.Errorf("not implemented: %s")' if resp_type else 'return fmt.Errorf("not implemented")'
            write_stub(path, pkg, typ, fn, req_type, resp_type, body)
            return
        imp, ctor, alias = info

    # build path vars from req
    path_vars = "nil"
    if req_type and ("{Id}" in route_path.replace(":id", "{Id}") or ":id" in route_path):
        path_vars = 'map[string]string{"id": fmt.Sprintf("%d", req.Id)}'
        route_path_clean = route_path.split("(")[0].strip()
    else:
        route_path_clean = route_path.split("(")[0].strip()
    if ":file" in route_path_clean and req_type:
        path_vars = 'map[string]string{"file": req.File}'
    if ":code" in route_path_clean and req_type:
        path_vars = 'map[string]string{"code": req.Code}'

    query = "nil"
    body = "nil"
    if req_type == "PageReq":
        query = 'url.Values{"page": {fmt.Sprintf("%d", req.Page)}, "page_size": {fmt.Sprintf("%d", req.PageSize)}}'
    elif req_type == "JSONBody":
        body = "req"
    elif req_type and req_type not in ("IdPathReq", "FilePathReq", "CodePathReq", "PageReq", "EmptyResp"):
        body = "req"
    elif req_type == "IdPathReq" and http in ("POST", "PUT", "PATCH"):
        body = "req"

    if resp_type is None:
        call = f'''_, err := httpinvoke.Run(ctx, "{http}", "{route_path_clean}", {path_vars}, {query}, {body}, {ctor})
	if err != nil {{
		return err
	}}
	return nil'''
    elif resp_type == "EmptyResp":
        call = f'''_, err := httpinvoke.Run(ctx, "{http}", "{route_path_clean}", {path_vars}, {query}, {body}, {ctor})
	if err != nil {{
		return nil, err
	}}
	return &types.EmptyResp{{}}, nil'''
    elif resp_type == "PageListResp":
        call = f'''raw, err := httpinvoke.Run(ctx, "{http}", "{route_path_clean}", {path_vars}, {query}, {body}, {ctor})
	if err != nil {{
		return nil, err
	}}
	var out types.PageListResp
	if err := httpinvoke.Decode(raw, &out); err != nil {{
		var list interface{{}}
		if err2 := httpinvoke.Decode(raw, &list); err2 == nil {{
			return &types.PageListResp{{List: list}}, nil
		}}
		return nil, err
	}}
	return &out, nil'''
    else:
        call = f'''raw, err := httpinvoke.Run(ctx, "{http}", "{route_path_clean}", {path_vars}, {query}, {body}, {ctor})
	if err != nil {{
		return nil, err
	}}
	var data interface{{}}
	if err := httpinvoke.Decode(raw, &data); err != nil {{
		return nil, err
	}}
	return &types.{resp_type}{{Data: data}}, nil'''

    alias_imp = f'\t{alias} "{imp}"' if not ctor.startswith(alias + ".") and f"{alias}." not in ctor else f'\t{alias} "{imp}"'
    # always import alias from ctor
    am = re.match(r"(\w+)\.", ctor)
    if am:
        alias = am.group(1)
        alias_imp = f'\t{alias} "{imp}"'

    req_param = f", req *types.{req_type}" if req_type else ""
    if resp_type:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) (resp *types.{resp_type}, err error)"
    else:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) error"

    out = f'''package {pkg}

import (
	"context"
	"fmt"
	"net/url"

	"mymall/pkg/httpinvoke"
	"mymall/services/{svc}/internal/svc"
	"mymall/services/{svc}/internal/types"
{alias_imp}

	"github.com/zeromicro/go-zero/core/logx"
)

type {typ} struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{typ}(svcCtx *svc.ServiceContext) *{typ} {{
	return &{typ}{{
		Logger: logx.WithContext(context.Background()),
		svcCtx: svcCtx,
	}}
}}

{sig} {{
	_ = fmt.Sprintf
	_ = url.Values{{}}
{call}
}}
'''
    path.write_text(out)


def write_stub(path, pkg, typ, fn, req_type, resp_type, body):
    req_param = f", req *types.{req_type}" if req_type else ""
    if resp_type:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) (resp *types.{resp_type}, err error)"
    else:
        sig = f"func (l *{typ}) {fn}(ctx context.Context{req_param}) error"
    svc = path.parts[path.parts.index("services") + 1]
    out = f'''package {pkg}

import (
	"context"
	"fmt"

	"mymall/services/{svc}/internal/svc"
	"mymall/services/{svc}/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type {typ} struct {{
	logx.Logger
	svcCtx *svc.ServiceContext
}}

func New{typ}(svcCtx *svc.ServiceContext) *{typ} {{
	return &{typ}{{Logger: logx.WithContext(context.Background()), svcCtx: svcCtx}}
}}

{sig} {{
	_ = ctx
	{(" _ = req" if req_type else "")}
	{body}
}}
'''
    path.write_text(out)


def write_health_special(path, pkg, typ, fn, svc, resp_type):
    if fn == "Healthz":
        body = "return &types.EmptyResp{}, nil"
    elif fn == "Readyz":
        body = "return &types.EmptyResp{}, nil"
    else:
        body = "return nil"
        resp_type = None
    write_stub(path, pkg, typ, fn, None, resp_type, body)


def main():
    svc = sys.argv[1]
    svc_dir = ROOT / "services" / svc
    api = next(svc_dir.joinpath("api").glob("*.api"))
    routes = {h: (m, p) for h, m, p in parse_routes(api)}
    discovered = discover_app_handlers(svc_dir)
    explicit = load_map(svc)
    print(f"discovered {len(discovered)} app methods, {len(routes)} routes")

    for logic in (svc_dir / "internal" / "logic").rglob("*_logic.go"):
        text = logic.read_text()
        fn_m = re.search(r"func \(l \*\w+\) (\w+)\(", text)
        if not fn_m:
            continue
        fn = fn_m.group(1)
        http, path = routes.get(fn, ("GET", "/"))
        write_logic_file(logic, svc, fn, http, path, discovered, explicit)
        print("gen", logic.relative_to(ROOT))


if __name__ == "__main__":
    main()
