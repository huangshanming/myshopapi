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
    target_name = out_name(dir_path, handler_root)
    target_path = dir_path / target_name
    files = sorted(
        p
        for p in dir_path.glob("*_handler.go")
        if p.is_file() and p.name != "routes.go"
    )
    leaf = [f for f in files if f.name != target_name]

    # Already consolidated or empty — if single file has wrong name, rename
    if not leaf:
        if len(files) == 1 and files[0].name != target_name:
            text = files[0].read_text()
            files[0].unlink()
            target_path.write_text(text)
        return

    def extract(path: Path) -> tuple[str, dict[str, str], dict[str, str]]:
        text = path.read_text()
        pm = re.search(r"^package\s+(\w+)", text, re.M)
        if not pm:
            raise SystemExit(f"no package in {path}")
        imports: dict[str, str] = {}
        im = re.search(r"import\s*\(([\s\S]*?)\)", text)
        if im:
            for line in im.group(1).splitlines():
                s = line.strip()
                if not s or s.startswith("//"):
                    continue
                parts = s.split()
                path_s = parts[-1].strip('"')
                imports[path_s] = s
        funcs: dict[str, str] = {}
        for m in re.finditer(
            r"^func (\w+Handler)\(.*?^}\n?",
            text,
            re.M | re.S,
        ):
            funcs[m.group(1)] = m.group(0).rstrip() + "\n"
        return pm.group(1), imports, funcs

    package = None
    imports: dict[str, str] = {}
    funcs_by_name: dict[str, str] = {}

    # Prefer previous merge target bodies (hand-tuned, e.g. multipart upload).
    if target_path.is_file():
        package, imports, funcs_by_name = extract(target_path)

    for f in leaf:
        pkg, imps, funcs = extract(f)
        if package is None:
            package = pkg
        elif package != pkg:
            raise SystemExit(f"package mismatch in {f}: {pkg} vs {package}")
        added_new = False
        for name, body in funcs.items():
            if name not in funcs_by_name:
                funcs_by_name[name] = body
                added_new = True
        # Only pull leaf imports when introducing new handlers (avoids unused imports).
        if added_new or not funcs_by_name:
            imports.update(imps)
        elif not target_path.is_file():
            imports.update(imps)

    if not funcs_by_name:
        raise SystemExit(f"no handler funcs in {dir_path}")

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

    func_bodies = [funcs_by_name[k] for k in sorted(funcs_by_name)]
    out = f"""package {package}

import (
{import_block}
)

{chr(10).join(func_bodies)}"""
    if not out.endswith("\n"):
        out += "\n"

    for f in dir_path.glob("*_handler.go"):
        if f.is_file() and f.name != "routes.go":
            f.unlink()
    target_path.write_text(out)


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
