from fastapi import APIRouter

from app.conf import get_config

router = APIRouter(tags=["health"])


@router.get("/healthz")
def healthz() -> dict:
    return {"status": "ok", "service": "askdata-service", "product": "电商问数"}


@router.get("/readyz")
def readyz() -> dict:
    cfg = get_config()
    return {
        "status": "ready",
        "mode": cfg.get("mode", "dev"),
        "order": cfg.get("upstream", {}).get("order_http"),
        "catalog": cfg.get("upstream", {}).get("catalog_http"),
        "model": cfg.get("llm", {}).get("model"),
        "elasticsearch": cfg.get("elasticsearch", {}).get("url"),
        "qdrant": cfg.get("qdrant", {}).get("url"),
    }
