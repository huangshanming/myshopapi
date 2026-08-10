from fastapi import APIRouter

from app.config import get_settings

router = APIRouter(tags=["health"])


@router.get("/healthz")
def healthz() -> dict:
    return {"status": "ok", "service": "askdata-service", "product": "电商问数"}


@router.get("/readyz")
def readyz() -> dict:
    settings = get_settings()
    return {
        "status": "ready",
        "mode": settings.mode,
        "order": settings.order_http,
        "catalog": settings.catalog_http,
        "model": settings.askdata_model,
    }
