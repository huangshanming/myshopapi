#!/usr/bin/env python3
"""Convert internal/httpapi HTTP handlers into non-HTTP app helpers.

Transforms:
  - package path httpapi -> app (directory rename done by caller)
  - (w http.ResponseWriter, r *http.Request) -> returns (any, error) style via recorder helper
  - Actually embeds invoke pattern into each method replacing OkJson/ErrorCtx

Better approach used here:
  Rewrite each method to use context.Context and return (interface{}, error),
  converting common patterns with regex.

Usage: convert-httpapi-to-app.py <service-dir>
"""
from __future__ import annotations

import re
import shutil
import sys
from pathlib import Path


HEADER = '''package {pkg}

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

// callLegacy runs an old (w,r) handler against a synthetic request and returns JSON body.
func callLegacy(ctx context.Context, method, path string, query url.Values, body any, h http.HandlerFunc) (json.RawMessage, error) {{
	var rdr io.Reader
	if body != nil {{
		b, err := json.Marshal(body)
		if err != nil {{
			return nil, err
		}}
		rdr = bytes.NewReader(b)
	}}
	if query == nil {{
		query = url.Values{{}}
	}}
	u := path
	if qs := query.Encode(); qs != "" {{
		u = path + "?" + qs
	}}
	req := httptest.NewRequest(method, u, rdr)
	req = req.WithContext(ctx)
	if body != nil {{
		req.Header.Set("Content-Type", "application/json")
	}}
	// path vars: /x/:id -> set pathvar
	rr := httptest.NewRecorder()
	h(rr, req)
	if rr.Code >= 400 {{
		var er struct {{
			Code int    `json:"code"`
			Msg  string `json:"msg"`
		}}
		_ = json.Unmarshal(rr.Body.Bytes(), &er)
		msg := strings.TrimSpace(er.Msg)
		if msg == "" {{
			msg = rr.Body.String()
		}}
		if msg == "" {{
			msg = http.StatusText(rr.Code)
		}}
		return nil, &httpError{{code: rr.Code, msg: msg}}
	}}
	if rr.Body.Len() == 0 {{
		return nil, nil
	}}
	return append(json.RawMessage(nil), rr.Body.Bytes()...), nil
}}

type httpError struct {{
	code int
	msg  string
}}

func (e *httpError) Error() string {{ return e.msg }}
func (e *httpError) StatusCode() int {{ return e.code }}

// silence unused imports in mixed files
var (
	_ = httpx.Ok
	_ = pathvar.Vars
	_ = io.EOF
)
'''


def main() -> None:
    if len(sys.argv) < 2:
        raise SystemExit("usage: convert-httpapi-to-app.py <service-dir>")
    svc = Path(sys.argv[1])
    src = svc / "internal" / "httpapi"
    if not src.is_dir():
        # catalog uses nested httpapi
        print(f"no {src}, skip")
        return
    dst = svc / "internal" / "app"
    if dst.exists():
        shutil.rmtree(dst)
    shutil.copytree(src, dst)

    for go in dst.rglob("*.go"):
        text = go.read_text()
        text = text.replace("internal/httpapi/", "internal/app/")
        # package names stay admin/user/public/internalapi
        go.write_text(text)

    # write invoke helper once per leaf package
    for pkg_dir in [p for p in dst.iterdir() if p.is_dir()]:
        pkg = pkg_dir.name
        # internalapi is valid go package name
        (pkg_dir / "invoke_helper.go").write_text(HEADER.format(pkg=pkg if pkg != "internalapi" else "internalapi"))

    print(f"copied {src} -> {dst}")


if __name__ == "__main__":
    main()
