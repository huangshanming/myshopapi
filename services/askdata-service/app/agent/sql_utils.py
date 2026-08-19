"""从 LLM 输出中抽出可执行 SQL。"""

import re

_FENCE_RE = re.compile(r"```(?:sql)?\s*([\s\S]*?)```", re.IGNORECASE)


def extract_sql(text: str) -> str:
    """去掉 markdown 代码块，只保留纯 SQL。"""
    if not text:
        return ""
    text = text.strip()
    match = _FENCE_RE.search(text)
    if match:
        text = match.group(1)
    return text.strip()
