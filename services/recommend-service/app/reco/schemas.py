from __future__ import annotations

from pydantic import BaseModel, ConfigDict, Field


class RecommendItem(BaseModel):
    product_id: int
    score: float = 0.0
    reason: str = ""


class RecommendListResp(BaseModel):
    model_config = ConfigDict(populate_by_name=True)

    items: list[RecommendItem] = Field(default_factory=list, alias="list")
    total: int = 0
