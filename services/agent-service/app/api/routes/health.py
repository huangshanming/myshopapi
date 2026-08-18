from fastapi import APIRouter

from app.config import get_settings

router = APIRouter(tags=["health"])


@router.get("/healthz")
def healthz() -> dict:
    return {"status": "ok", "service": "agent-service"}


@router.get("/readyz")
def readyz() -> dict:
    settings = get_settings()
    return {
        "status": "ready",
        "mode": settings.mode,
        "catalog": settings.catalog_http,
    }
