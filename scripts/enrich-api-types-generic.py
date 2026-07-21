#!/usr/bin/env python3
"""Generic .api type enrichment: give every route a Req/Resp."""
from __future__ import annotations

import re
import sys
from pathlib import Path

COMMON = r'''
type PageReq {
	Page     int `form:"page,default=1"`
	PageSize int `form:"page_size,default=10"`
}
type IdPathReq {
	Id uint64 `path:"id"`
}
type FilePathReq {
	File string `path:"file"`
}
type CodePathReq {
	Code string `path:"code"`
}
type EmptyResp {}
type PageListResp {
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}
type AnyResp {
	Data interface{} `json:"data,optional"`
}
type JSONBody {
	// placeholder — body decoded inside legacy app via httpinvoke raw JSON
}
'''


def main() -> None:
    path = Path(sys.argv[1])
    text = path.read_text()
    if "type PageReq {" in text:
        print("already enriched", path)
        return

    m = re.match(r"(syntax\s*=\s*\"v1\"\s*\n)", text)
    if not m:
        raise SystemExit("no syntax")
    text = text[: m.end()] + "\n" + COMMON + "\n" + text[m.end() :]

    def repl(mm: re.Match) -> str:
        indent, handler, method, pth = mm.group(1), mm.group(2), mm.group(3), mm.group(4)
        method_l = method.lower()
        req = None
        resp = "AnyResp"
        if ":id" in pth or "/:id" in pth:
            if method_l in ("get", "delete"):
                req = "IdPathReq"
            else:
                req = "IdPathReq"  # path id; body still passed via invoke if needed
        if ":file" in pth:
            req = "FilePathReq"
            resp = None
        if ":code" in pth and method_l == "post":
            req = "CodePathReq"
        if method_l == "get" and "page" in handler.lower() or (
            method_l == "get" and any(x in handler.lower() for x in ("list", "logs", "orders", "products", "users"))
        ):
            if req is None:
                req = "PageReq"
            resp = "PageListResp"
        if method_l in ("post", "put", "patch") and req is None:
            req = "JSONBody"
        if handler in ("Healthz", "Readyz"):
            req, resp = None, "EmptyResp"
        if handler == "Metrics":
            req, resp = None, None
        line = f"{indent}{method} {pth}"
        if req:
            line += f" ({req})"
        if resp:
            line += f" returns ({resp})"
        return f"{indent}@handler {handler}\n{line}"

    text2, n = re.subn(
        r"^([ \t]*)@handler\s+(\w+)\s*\n[ \t]*(get|post|put|delete|patch)\s+(\S+)",
        repl,
        text,
        flags=re.M | re.I,
    )
    path.write_text(text2)
    print(f"enriched {path}: {n} routes")


if __name__ == "__main__":
    main()
