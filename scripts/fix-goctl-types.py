#!/usr/bin/env python3
"""Post-process goctl-generated internal/types/types.go.

goctl emits DataResp / named *Resp as plain structs. This repo keeps opaque
entity JSON unwrapped via AnyResp.MarshalJSON, so we:

  1. Replace `type DataResp struct {...}` with `type DataResp = AnyResp`
  2. Drop structs whose names are aliased in entity_resp.go (`type Xxx = DataResp`)
  3. Fail if any type name is still declared in more than one .go file under types/

Usage: python3 scripts/fix-goctl-types.py <service-dir>
"""
from __future__ import annotations

import re
import sys
from collections import defaultdict
from pathlib import Path

STRUCT_RE = re.compile(
    r"^type\s+(\w+)\s+struct\s*\{",
    re.M,
)
ALIAS_RE = re.compile(
    r"^type\s+(\w+)\s*=\s*",
    re.M,
)
ENTITY_ALIAS_RE = re.compile(
    r"^type\s+(\w+)\s*=\s*(?:DataResp|AnyResp)\s*$",
    re.M,
)
DATA_RESP_STRUCT = (
    "type DataResp struct {\n"
    '\tData interface{} `json:"data,optional"`\n'
    "}\n"
)
ANY_RESP_STRUCT = (
    "type AnyResp struct {\n"
    '\tData interface{} `json:"data,optional"`\n'
    "}\n"
)
DATA_RESP_ALIAS = (
    "type AnyResp struct {\n"
    '\tData interface{} `json:"data,optional"`\n'
    "}\n\n"
    "// DataResp is the preferred name for opaque entity JSON bodies.\n"
    "type DataResp = AnyResp\n"
)


def entity_aliases(path: Path) -> set[str]:
    if not path.exists():
        return set()
    return set(ENTITY_ALIAS_RE.findall(path.read_text()))


def strip_named_struct(text: str, name: str) -> str:
    # Struct bodies are single-level; DOTALL so field lines match.
    pat = rf"(?ms)^type {re.escape(name)} struct \{{.*?\}}\n+"
    return re.sub(pat, "", text)


def alias_data_resp(text: str) -> str:
    text = text.replace(DATA_RESP_STRUCT, "")
    if "type DataResp =" in text:
        return text
    if ANY_RESP_STRUCT not in text:
        return text
    return text.replace(ANY_RESP_STRUCT, DATA_RESP_ALIAS, 1)


def declared_names(path: Path) -> set[str]:
    text = path.read_text()
    return set(STRUCT_RE.findall(text)) | set(ALIAS_RE.findall(text))


def check_duplicates(types_dir: Path) -> list[str]:
    by_name: dict[str, list[str]] = defaultdict(list)
    for p in sorted(types_dir.glob("*.go")):
        for name in declared_names(p):
            by_name[name].append(p.name)
    errs: list[str] = []
    for name, files in sorted(by_name.items()):
        if len(files) > 1:
            errs.append(f"  {name}: {', '.join(files)}")
    return errs


def main() -> int:
    if len(sys.argv) != 2:
        print(f"usage: {sys.argv[0]} <service-dir>", file=sys.stderr)
        return 2
    svc = Path(sys.argv[1]).resolve()
    types_dir = svc / "internal" / "types"
    types_go = types_dir / "types.go"
    if not types_go.exists():
        print(f"skip: no {types_go}", file=sys.stderr)
        return 0

    aliases = entity_aliases(types_dir / "entity_resp.go")
    text = types_go.read_text()
    for name in sorted(aliases):
        text = strip_named_struct(text, name)
    text = alias_data_resp(text)
    # collapse excess blank lines
    text = re.sub(r"\n{3,}", "\n\n", text)
    types_go.write_text(text)
    print(
        f"==> fix-goctl-types ({svc.name}): "
        f"DataResp alias, stripped {len(aliases)} entity resp struct(s)"
    )

    dups = check_duplicates(types_dir)
    if dups:
        print(f"ERROR: duplicate type names under {types_dir}:", file=sys.stderr)
        print("\n".join(dups), file=sys.stderr)
        print(
            "Keep API DTOs only in goctl types.go; biz_types.go may only hold "
            "types that are NOT declared in api/*.api.",
            file=sys.stderr,
        )
        return 1
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
