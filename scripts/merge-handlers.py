#!/usr/bin/env python3
"""Merge goctl one-file-per-route handlers into one file per leaf package.

goctl always emits create_handler.go / update_handler.go etc. This script
collapses each directory under internal/handler/ (except the root routes.go)
into a single <path_with_underscores>_handler.go, e.g.:

  internal/handler/admin/article/*.go  →  admin_article_handler.go

Usage: merge-handlers.py <service-dir>
"""
from __future__ import annotations

import ast
import re
import sys
from pathlib import Path


def merge_dir(dir_path: Path, handler_root: Path) -> None:
    files = sorted(
        p
        for p in dir_path.glob("*_handler.go")
        if p.is_file() and p.name != "routes.go"
    )
    # Also include any .go that is a handler file from goctl (all except routes)
    # Prefer *_handler.go only.
    if len(files) <= 1:
        # Already consolidated or empty — if single file has wrong name, rename
        if len(files) == 1:
            target = out_name(dir_path, handler_root)
            if files[0].name != target:
                text = files[0].read_text()
                files[0].unlink()
                (dir_path / target).write_text(text)
        return

    package = None
    imports: dict[str, str] = {}  # path -> full import line (with alias)
    funcs: list[str] = []

    for f in files:
        text = f.read_text()
        pm = re.search(r"^package\s+(\w+)", text, re.M)
        if not pm:
            raise SystemExit(f"no package in {f}")
        if package is None:
            package = pm.group(1)
        elif package != pm.group(1):
            raise SystemExit(f"package mismatch in {f}: {pm.group(1)} vs {package}")

        im = re.search(r"import\s*\(([\s\S]*?)\)", text)
        if im:
            for line in im.group(1).splitlines():
                s = line.strip()
                if not s or s.startswith("//"):
                    continue
                # "path" or alias "path"
                parts = s.split()
                path = parts[-1].strip('"')
                imports[path] = s

        # Extract top-level func blocks (handlers)
        for m in re.finditer(
            r"^func (\w+Handler)\(.*?^}\n?",
            text,
            re.M | re.S,
        ):
            funcs.append(m.group(0).rstrip() + "\n")

    if not funcs:
        raise SystemExit(f"no handler funcs in {dir_path}")

    # Format imports
    std, third, local = [], [], []
    for path, line in sorted(imports.items(), key=lambda x: x[0]):
        if path.startswith("mymall/") or "/internal/" in path:
            local.append(line)
        elif "." in path.split("/")[0]:
            third.append(line)
        else:
            std.append(line)

    sections = []
    for block in (std, third, local):
        if block:
            sections.append("\n".join("\t" + x for x in block))
    import_block = "\n\n".join(sections)

    out = f"""package {package}

import (
{import_block}
)

{chr(10).join(funcs)}"""
    # Ensure trailing newline
    if not out.endswith("\n"):
        out += "\n"

    target = dir_path / out_name(dir_path, handler_root)
    # Remove old files first
    for f in files:
        f.unlink()
    target.write_text(out)


def out_name(dir_path: Path, handler_root: Path) -> str:
    rel = dir_path.relative_to(handler_root)
    parts = list(rel.parts)
    if parts == ["."]:
        return "handlers.go"
    return "_".join(parts) + "_handler.go"


def main() -> None:
    if len(sys.argv) != 2:
        print("usage: merge-handlers.py <service-dir>", file=sys.stderr)
        sys.exit(2)
    service = Path(sys.argv[1])
    handler_root = service / "internal" / "handler"
    if not handler_root.is_dir():
        print(f"no handler dir: {handler_root}", file=sys.stderr)
        sys.exit(1)

    # Leaf dirs: any directory that contains *_handler.go
    dirs = set()
    for f in handler_root.rglob("*_handler.go"):
        if f.parent == handler_root:
            # flat handlers at root — rare; merge into handlers.go
            dirs.add(f.parent)
        else:
            dirs.add(f.parent)

    for d in sorted(dirs, key=lambda p: str(p)):
        merge_dir(d, handler_root)
        print(f"merged {d.relative_to(handler_root)}")


if __name__ == "__main__":
    main()
