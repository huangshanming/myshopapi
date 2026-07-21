#!/usr/bin/env python3
"""Convert app HTTP handlers to (ctx, appinput.CallInput) (any, error)
and rewrite logic to call them instead of httpinvoke.
"""
from __future__ import annotations

import re
import sys
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]

APP_METHOD = re.compile(
    r"func \((?P<recv>[^)]+)\) (?P<name>\w+)\(w http\.ResponseWriter, r \*http\.Request\) "
)
PLAIN_METHOD = re.compile(
    r"^func (?P<name>\w+)\(w http\.ResponseWriter, r \*http\.Request\) ",
    re.M,
)
MULTIPART_MARKERS = ("ParseMultipartForm", "FormFile", "MultipartForm")


def skip_string(text: str, i: int) -> int:
    q = text[i]
    i += 1
    if q == "`":
        while i < len(text) and text[i] != "`":
            i += 1
        return i + 1 if i < len(text) else i
    while i < len(text):
        if text[i] == "\\":
            i += 2
            continue
        if text[i] == q:
            return i + 1
        i += 1
    return i


def find_matching(text: str, open_idx: int, open_ch: str = "(", close_ch: str = ")") -> int:
    assert text[open_idx] == open_ch
    depth = 0
    i = open_idx
    while i < len(text):
        c = text[i]
        if c in ('"', "'", "`"):
            i = skip_string(text, i)
            continue
        if c == open_ch:
            depth += 1
        elif c == close_ch:
            depth -= 1
            if depth == 0:
                return i
        i += 1
    raise RuntimeError(f"unbalanced {open_ch} at {open_idx}")


def find_brace_block(text: str, brace_idx: int) -> tuple[int, int]:
    end = find_matching(text, brace_idx, "{", "}")
    return brace_idx, end + 1


def replace_httpx_calls(body: str) -> str:
    """Replace httpx.Error*/OkJson* with return statements. Body includes outer braces."""
    s = body
    patterns = [
        # ErrorCtx(ctx_or_r, w, ERR) [;|\n] return
        ("ErrorCtx", True),
        ("Error", True),
        ("OkJsonCtx", False),
        ("OkJson", False),
    ]

    def find_calls(src: str, name: str) -> list[tuple[int, int, str]]:
        out = []
        needle = f"httpx.{name}("
        start = 0
        while True:
            i = src.find(needle, start)
            if i < 0:
                break
            open_paren = i + len(needle) - 1
            close = find_matching(src, open_paren)
            call = src[i : close + 1]
            out.append((i, close + 1, call))
            start = close + 1
        return out

    # Process from end so indices stay valid
    for name, is_err in patterns:
        calls = find_calls(s, name)
        for start, end, call in reversed(calls):
            # Extract last argument (error or value)
            open_paren = call.index("(")
            close_rel = len(call) - 1
            inner = call[open_paren + 1 : close_rel]
            # split top-level commas
            args = split_args(inner)
            if is_err:
                expr = args[-1].strip()
                repl = f"return nil, {expr}\n"
            else:
                expr = args[-1].strip()
                repl = f"return {expr}, nil\n"

            # Consume following whitespace + bare `return`
            j = end
            while j < len(s) and s[j] in " \t\r":
                j += 1
            if j < len(s) and s[j] == "\n":
                j += 1
            k = j
            while k < len(s) and s[k] in " \t":
                k += 1
            if s.startswith("return", k):
                line_end = s.find("\n", k)
                if line_end < 0:
                    line_end = len(s)
                line = s[k:line_end].strip()
                if line == "return" or line == "return;":
                    j = line_end + 1 if line_end < len(s) else line_end

            s = s[:start] + repl + s[j:]
    return s


def split_args(inner: str) -> list[str]:
    args = []
    depth = 0
    start = 0
    i = 0
    while i < len(inner):
        c = inner[i]
        if c in ('"', "'", "`"):
            i = skip_string(inner, i)
            continue
        if c in "({[":
            depth += 1
        elif c in ")}]":
            depth -= 1
        elif c == "," and depth == 0:
            args.append(inner[start:i])
            start = i + 1
        i += 1
    args.append(inner[start:])
    return args


def transform_body(body: str, multipart: bool) -> str:
    s = replace_httpx_calls(body)

    s = re.sub(
        r"if err := json\.NewDecoder\(r\.Body\)\.Decode\(&(\w+)\); err != nil",
        r"if err := appinput.BindBody(in, &\1); err != nil",
        s,
    )
    s = re.sub(
        r"_ = json\.NewDecoder\(r\.Body\)\.Decode\(&(\w+)\)",
        r"_ = appinput.BindBody(in, &\1)",
        s,
    )
    s = re.sub(
        r"json\.NewDecoder\(r\.Body\)\.Decode\(&(\w+)\)",
        r"appinput.BindBody(in, &\1)",
        s,
    )

    s = s.replace("middleware.ParsePage(r)", "in.Page()")
    s = re.sub(r'httpserver\.PathParam\(r,\s*"([^"]+)"\)', r'in.Path("\1")', s)
    s = re.sub(r'pathvar\.Vars\(r\)\["([^"]+)"\]', r'in.Path("\1")', s)
    s = re.sub(r"r\.URL\.Query\(\)\.Get\(([^)]+)\)", r"in.QueryGet(\1)", s)
    s = re.sub(r"r\.URL\.Query\(\)", "in.Query", s)
    s = re.sub(r"r\.FormValue\(([^)]+)\)", r"in.QueryGet(\1)", s)
    s = s.replace("r.Context()", "ctx")

    if multipart:
        s = re.sub(r"\br\.(ParseMultipartForm|FormFile|MultipartForm|Form)\b", r"in.Request.\1", s)
        # remaining bare r (careful)
        s = re.sub(r"(?<![.\w])r(?=[\.\)\,\s])", "in.Request", s)
        s = s.replace("in.Request.Context()", "ctx")

    return s


def patch_imports(text: str) -> str:
    # Use goimports later; ensure required imports present and drop httpx/httpserver if unused.
    if "mymall/pkg/appinput" not in text:
        if "import (" in text:
            text = text.replace("import (", 'import (\n\t"mymall/pkg/appinput"', 1)
        else:
            text = re.sub(
                r"(package \w+\n)",
                r'\1\nimport (\n\t"mymall/pkg/appinput"\n)\n',
                text,
                count=1,
            )
    if '"context"' not in text and "context.Context" in text:
        text = text.replace("import (", 'import (\n\t"context"', 1)

    # remove httpx import line
    text = re.sub(r'\n\t"github.com/zeromicro/go-zero/rest/httpx"', "", text)
    text = re.sub(r'\n\t"mymall/pkg/httpserver"', "", text)
    return text


def convert_require_perm(text: str) -> str:
    # func (h *X) requirePerm(w http.ResponseWriter, r *http.Request, code string) bool {
    for m in list(
        re.finditer(
            r"func \(([^)]+)\) requirePerm\(w http\.ResponseWriter, r \*http\.Request, code string\) bool ",
            text,
        )
    ):
        brace = text.find("{", m.end() - 1)
        b0, b1 = find_brace_block(text, brace)
        body = text[b0:b1]
        # typical: write error and return false -> return err; success return nil
        new_body = replace_httpx_calls(body)
        new_body = new_body.replace("r.Context()", "ctx")
        # return false -> need error; this is messy. Simpler rewrite:
        # if original returns false after ErrorCtx, ErrorCtx already became return nil, err — wrong for bool helper.
        # Manual pattern for catalog:
        pass

    text = re.sub(
        r"func \(([^)]+)\) shopUser\(r \*http\.Request\)",
        r"func (\1) shopUser(ctx context.Context)",
        text,
    )
    text = text.replace(".shopUser(r)", ".shopUser(ctx)")
    return text


def convert_app_file(path: Path) -> bool:
    text = path.read_text()
    if "http.ResponseWriter" not in text and "*http.Request" not in text:
        return False

    replacements: list[tuple[int, int, str]] = []

    for rx in (APP_METHOD, PLAIN_METHOD):
        for m in rx.finditer(text):
            # find opening brace of function
            brace = text.find("{", m.end())
            if brace < 0:
                continue
            b0, b1 = find_brace_block(text, brace)
            body = text[b0:b1]
            multipart = any(x in body for x in MULTIPART_MARKERS)
            name = m.group("name")
            if "recv" in m.groupdict() and m.group("recv"):
                new_sig = f"func ({m.group('recv')}) {name}(ctx context.Context, in appinput.CallInput) (any, error) "
            else:
                new_sig = f"func {name}(ctx context.Context, in appinput.CallInput) (any, error) "
            new_body = transform_body(body, multipart=multipart)
            if multipart:
                inner = new_body[1:-1]
                new_body = (
                    "{\n\tif in.Request == nil {\n"
                    '\t\treturn nil, xerr.New(http.StatusBadRequest, "缺少上传请求")\n\t}\n'
                    + inner
                    + "}"
                )
            if not re.search(r"\breturn\b", new_body):
                new_body = new_body[:-1] + "\n\treturn nil, nil\n}"
            replacements.append((m.start(), b1, new_sig + new_body))

    if not replacements:
        return False

    replacements.sort(key=lambda x: x[0], reverse=True)
    for start, end, repl in replacements:
        text = text[:start] + repl + text[end:]

    text = convert_require_perm(text)
    text = patch_imports(text)
    path.write_text(text)
    return True


def parse_run_args(call_inner: str) -> tuple[str, str, str, str, str]:
    """pathvars, query, body, handler from Run(ctx, method, path, pathVars, query, body, h)."""
    args = split_args(call_inner)
    if len(args) != 7:
        raise ValueError(f"expected 7 args, got {len(args)}: {call_inner[:160]}")
    return args[3].strip(), args[4].strip(), args[5].strip(), args[6].strip()


def convert_logic_file(path: Path) -> bool:
    text = path.read_text()
    if "httpinvoke.Run" not in text:
        return False

    # Replace each httpinvoke.Run(...) occurrence
    out = []
    i = 0
    changed = False
    while True:
        j = text.find("httpinvoke.Run(", i)
        if j < 0:
            out.append(text[i:])
            break
        out.append(text[i:j])
        open_paren = j + len("httpinvoke.Run") 
        # httpinvoke.Run(  -> open at position of (
        open_paren = text.find("(", j)
        close = find_matching(text, open_paren)
        inner = text[open_paren + 1 : close]
        try:
            pathvars, query, body, handler = parse_run_args(inner)
        except ValueError as e:
            print(f"WARN {path}: {e}")
            out.append(text[j : close + 1])
            i = close + 1
            continue

        parts = []
        if pathvars != "nil":
            parts.append(f"PathVars: {pathvars}")
        if query != "nil":
            parts.append(f"Query: {query}")
        if body != "nil":
            parts.append(f"Body: {body}")
        in_lit = "appinput.CallInput{" + ", ".join(parts) + "}"
        out.append(f"{handler}(ctx, {in_lit})")
        changed = True
        i = close + 1

    if not changed:
        return False

    new_text = "".join(out)

    # raw, err := CALL  -> data, err := CALL
    new_text = re.sub(r"\braw,\s*err\s*:=", "data, err :=", new_text)

    # Remove AnyResp decode boilerplate
    new_text = re.sub(
        r"\n\tvar data interface\{\}\n\tif err := httpinvoke\.Decode\(data, &data\); err != nil \{\n\t\treturn nil, err\n\t\}\n\treturn &types\.AnyResp\{Data: data\}, nil",
        "\n\treturn &types.AnyResp{Data: data}, nil",
        new_text,
    )
    # After raw->data, Decode(raw) might still say raw — fix Decode(raw
    new_text = new_text.replace("httpinvoke.Decode(raw,", "httpinvoke.Decode(data,")
    new_text = re.sub(
        r"\n\tvar data interface\{\}\n\tif err := httpinvoke\.Decode\(data, &data\); err != nil \{\n\t\treturn nil, err\n\t\}\n\treturn &types\.AnyResp\{Data: data\}, nil",
        "\n\treturn &types.AnyResp{Data: data}, nil",
        new_text,
    )

    # var data interface{}; Decode — when we already have data from call, remove shadowing
    new_text = re.sub(
        r"if err != nil \{\n\t\treturn nil, err\n\t\}\n\tvar data interface\{\}\n\tif err := httpinvoke\.Decode\(data, &data\); err != nil \{\n\t\treturn nil, err\n\t\}\n\treturn &types\.AnyResp\{Data: data\}, nil",
        "if err != nil {\n\t\treturn nil, err\n\t}\n\treturn &types.AnyResp{Data: data}, nil",
        new_text,
    )

    # list interface
    new_text = re.sub(
        r"if err != nil \{\n\t\treturn nil, err\n\t\}\n\tvar list interface\{\}\n\tif err := httpinvoke\.Decode\(data, &list\); err != nil \{\n\t\treturn nil, err\n\t\}\n\treturn &types\.PageListResp\{List: list\}, nil",
        "if err != nil {\n\t\treturn nil, err\n\t}\n\treturn &types.PageListResp{List: data}, nil",
        new_text,
    )

    # Typed decode: httpinvoke.Decode(data, &out) -> json roundtrip
    def typed_decode(m: re.Match) -> str:
        out_name = m.group("out")
        typ = m.group("typ")
        return (
            f"b, _ := json.Marshal(data)\n"
            f"\tvar {out_name} {typ}\n"
            f"\tif err := json.Unmarshal(b, &{out_name}); err != nil {{"
        )

    new_text = re.sub(
        r"var (?P<out>\w+) (?P<typ>types\.\w+)\n\tif err := httpinvoke\.Decode\(data, &(?P=out)\); err != nil \{",
        typed_decode,
        new_text,
    )

    # Fallback decode with comments (points etc)
    new_text = re.sub(
        r"var (?P<out>\w+) (?P<typ>types\.\w+)\n\tif err := httpinvoke\.Decode\(data, &(?P=out)\); err != nil \{",
        typed_decode,
        new_text,
    )

    # Remove remaining httpinvoke.Decode if any with json.Unmarshal already
    # Fix references to raw in fallbacks
    new_text = new_text.replace("json.Unmarshal(raw,", "json.Unmarshal(b,")
    new_text = new_text.replace("json.Unmarshal(raw,", "json.Unmarshal(")  # noop if fixed

    # For fallback that used raw bytes — use b from marshal; if no b, marshal data
    if "json.Unmarshal(raw" in new_text:
        new_text = new_text.replace("json.Unmarshal(raw,", "json.Unmarshal(mustJSON(data),")

    # imports
    new_text = re.sub(r'\n\t"mymall/pkg/httpinvoke"', "", new_text)
    if "appinput." in new_text and "mymall/pkg/appinput" not in new_text:
        new_text = new_text.replace("import (", 'import (\n\t"mymall/pkg/appinput"', 1)
    if "json." in new_text and '"encoding/json"' not in new_text:
        new_text = new_text.replace("import (", 'import (\n\t"encoding/json"', 1)

    # Drop empty httpinvoke usage
    if "httpinvoke" in new_text:
        # still has Decode — convert remaining Decode(data, &x) 
        new_text = re.sub(
            r"if err := httpinvoke\.Decode\(data, &(\w+)\); err != nil",
            r"b, _ := json.Marshal(data); if err := json.Unmarshal(b, &\1); err != nil",
            new_text,
        )
        new_text = re.sub(r'\n\t"mymall/pkg/httpinvoke"', "", new_text)

    path.write_text(new_text)
    return True


def app_dirs(svc: str) -> list[Path]:
    base = ROOT / "services" / svc
    dirs = [base / "internal" / "app"]
    for nest in ("product", "content", "shopops", "notify"):
        dirs.append(base / "internal" / nest / "app")
    return [d for d in dirs if d.is_dir()]


def main() -> int:
    services = sys.argv[1:] or [
        "user-service",
        "merchant-service",
        "order-service",
        "catalog-service",
    ]
    for svc in services:
        n_app = n_logic = 0
        for d in app_dirs(svc):
            for go in sorted(d.rglob("*.go")):
                if convert_app_file(go):
                    n_app += 1
                    print(f"app: {go.relative_to(ROOT)}")
        logic = ROOT / "services" / svc / "internal" / "logic"
        if logic.is_dir():
            for go in sorted(logic.rglob("*.go")):
                if convert_logic_file(go):
                    n_logic += 1
                    print(f"logic: {go.relative_to(ROOT)}")
        print(f"{svc}: {n_app} app, {n_logic} logic")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
