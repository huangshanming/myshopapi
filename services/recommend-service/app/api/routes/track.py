from __future__ import annotations

from fastapi import APIRouter, Header
from pydantic import BaseModel, Field

from app.reco import service as reco_service

router = APIRouter(tags=["track"])


class TrackEvent(BaseModel):
    event: str = Field(..., description="expose | click | cart | favorite | purchase")
    product_id: int = Field(..., ge=1)
    ts: int | None = Field(default=None, description="unix ms, optional")
    extra: dict | None = None


class TrackReq(BaseModel):
    events: list[TrackEvent] = Field(default_factory=list)


class TrackResp(BaseModel):
    accepted: int = 0


@router.post("/track", response_model=TrackResp)
async def track(
    body: TrackReq,
    x_user_id: str | None = Header(default=None, alias="X-User-Id"),
) -> TrackResp:
    """行为上报 — 在 app.reco.service 中实现落 MQ / 落库."""
    user_id = int(x_user_id) if x_user_id and x_user_id.isdigit() else 0
    n = await reco_service.track(user_id=user_id, events=body.events)
    return TrackResp(accepted=n)
