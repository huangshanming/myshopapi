from fastapi import APIRouter

from app.clients import milvus_client
from app.config import get_settings

router = APIRouter(tags=["health"])


@router.get("/healthz")
def healthz() -> dict:
    return {"status": "ok", "service": "recommend-service"}


@router.get("/readyz")
def readyz() -> dict:
    settings = get_settings()
    milvus_ok = milvus_client.ping()
    return {
        "status": "ready" if milvus_ok else "degraded",
        "mode": settings.mode,
        "catalog": settings.catalog_http,
        "mysql": f"{settings.mysql_host}:{settings.mysql_port}/{settings.mysql_dbname}",
        "redis": f"{settings.redis_host}:{settings.redis_port}/{settings.redis_db}",
        "milvus": settings.milvus_uri,
        "milvus_ok": milvus_ok,
    }
