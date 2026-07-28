"""Business logic stubs — implement ItemCF recall / track here."""

from __future__ import annotations

from typing import Any

from app.reco.schemas import RecommendItem


async def also_bought(*, product_id: int, limit: int) -> list[RecommendItem]:
    # TODO: Redis I2I ZREVRANGE + filter + catalog enrich
    _ = (product_id, limit)
    return []


async def for_you(*, user_id: int, limit: int) -> list[RecommendItem]:
    # TODO: seeds → multi I2I merge → filter → hot fill
    _ = (user_id, limit)
    return []


async def track(*, user_id: int, events: list[Any]) -> int:
    # TODO: publish MQ / write warehouse
    _ = user_id
    return len(events or [])
