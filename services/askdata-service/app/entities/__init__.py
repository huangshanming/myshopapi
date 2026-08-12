"""业务实体（比 ORM 更贴近问数领域）。"""

from __future__ import annotations

from dataclasses import dataclass


@dataclass
class MetricQuery:
    name: str
    window_days: int
    merchant_id: int | None = None
