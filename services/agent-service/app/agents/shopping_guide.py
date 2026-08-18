from __future__ import annotations

from pydantic import BaseModel, Field

from app.clients.catalog import CatalogClient


class ShoppingGuidePlanRequest(BaseModel):
    query: str = Field(..., min_length=1, description="用户购物意图，如：给狗买冻干、预算200")
    budget: float | None = Field(default=None, ge=0, description="可选预算上限")
    pet_type: str | None = Field(default=None, description="可选：dog / cat 等")
    limit: int = Field(default=8, ge=1, le=30)


class ShoppingGuidePlanResponse(BaseModel):
    summary: str
    steps: list[str]
    product_ids: list[int] = Field(default_factory=list)
    products: list[dict] = Field(default_factory=list)
    source: str = "rule_stub"


class ShoppingGuideAgent:
    """智能导购规划 Agent。

    现状：关键词检索 catalog 列表 + 规则化步骤（不依赖 LLM）。
    后续可在此接入 LangGraph / OpenAI tools，并复用 CatalogClient。
    """

    def __init__(self, catalog: CatalogClient | None = None) -> None:
        self.catalog = catalog or CatalogClient()

    async def plan(self, req: ShoppingGuidePlanRequest) -> dict:
        products = await self.catalog.search_products(
            keyword=req.query,
            pet_type=req.pet_type,
            page_size=req.limit,
        )
        if req.budget is not None:
            products = [
                p
                for p in products
                if float(p.get("sale_price") or p.get("price") or 0) <= req.budget
            ]

        ids = [int(p["id"]) for p in products if p.get("id") is not None]
        steps = [
            f"理解需求：{req.query.strip()}",
            "在商品库中检索候选商品",
            "按预算与宠物类型过滤" if (req.budget is not None or req.pet_type) else "整理推荐清单",
            "给出可下单的商品建议",
        ]
        summary = (
            f"已为「{req.query.strip()}」规划导购方案，共 {len(products)} 个候选商品。"
            if products
            else f"暂未找到与「{req.query.strip()}」匹配的商品，可换关键词或放宽预算。"
        )
        resp = ShoppingGuidePlanResponse(
            summary=summary,
            steps=steps,
            product_ids=ids,
            products=products[: req.limit],
            source="rule_stub",
        )
        return resp.model_dump()
