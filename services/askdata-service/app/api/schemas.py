from __future__ import annotations

from pydantic import BaseModel, Field


class AskDataRequest(BaseModel):
    question: str = Field(
        ..., min_length=1, description="自然语言问数，如：昨天订单 GMV 是多少？"
    )
    merchant_id: int | None = Field(default=None, description="可选商家维度")
    days: int = Field(default=7, ge=1, le=90, description="默认统计窗口（天）")


class AskDataResponse(BaseModel):
    answer: str
    steps: list[str] = Field(default_factory=list)
    metrics: dict = Field(default_factory=dict)
    sql_preview: str | None = None
    source: str = "rule_stub"
