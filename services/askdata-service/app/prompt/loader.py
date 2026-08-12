from __future__ import annotations

from pathlib import Path

from app.conf import ROOT_DIR

PROMPTS_DIR = ROOT_DIR / "prompts"


def load_prompt(name: str) -> str:
    """从项目根 prompts/ 加载静态提示词模板。"""
    path = PROMPTS_DIR / name
    if not path.is_file():
        return ""
    return path.read_text(encoding="utf-8")
