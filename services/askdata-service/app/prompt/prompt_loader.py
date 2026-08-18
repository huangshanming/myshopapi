from __future__ import annotations

from pathlib import Path

from app.conf import ROOT_DIR


def load_prompt(name: str) -> str:
    prompt_path = Path(__file__).parents[2] / "prompts" / f"{name}.prompt"
    return prompt_path.read_text(encoding="utf-8")


