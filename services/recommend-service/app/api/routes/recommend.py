from __future__ import annotations

from fastapi import APIRouter, Header, Query

from app.reco import service as reco_service
from app.reco.schemas import RecommendListResp

router = APIRouter(tags=["recommend"])


@router.get("/also-bought", response_model=RecommendListResp)
async def also_bought(
    product_id: int = Query(..., ge=1),
    limit: int = Query(10, ge=1, le=50),
) -> RecommendListResp:
    """商详「买了又买」— 在 app.reco.service 中实现."""
    items = await reco_service.also_bought(product_id=product_id, limit=limit)
    return RecommendListResp(items=items, total=len(items))


@router.get("/for-you", response_model=RecommendListResp)
async def for_you(
    limit: int = Query(20, ge=1, le=50),
    x_user_id: str | None = Header(default=None, alias="X-User-Id"),
) -> RecommendListResp:
    """首页「猜你喜欢」— 在 app.reco.service 中实现."""
    user_id = int(x_user_id) if x_user_id and x_user_id.isdigit() else 0
    items = await reco_service.for_you(user_id=user_id, limit=limit)
    return RecommendListResp(items=items, total=len(items))
