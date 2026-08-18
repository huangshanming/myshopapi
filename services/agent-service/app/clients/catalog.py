from __future__ import annotations

import httpx

from app.config import get_settings


class CatalogClient:
    def __init__(self, base_url: str | None = None, timeout: float = 8.0) -> None:
        settings = get_settings()
        self.base_url = (base_url or settings.catalog_http).rstrip("/")
        self.timeout = timeout

    async def search_products(
        self,
        *,
        keyword: str,
        pet_type: str | None = None,
        page: int = 1,
        page_size: int = 8,
    ) -> list[dict]:
        params: dict = {
            "page": page,
            "page_size": page_size,
            "name": keyword,
        }
        if pet_type:
            params["pet_type"] = pet_type
        try:
            async with httpx.AsyncClient(timeout=self.timeout) as client:
                r = await client.get(f"{self.base_url}/api/v1/products/list", params=params)
                r.raise_for_status()
                data = r.json()
        except Exception:
            return []

        # go-zero success body may be DTO directly or wrapped
        if isinstance(data, dict):
            if isinstance(data.get("list"), list):
                return [x for x in data["list"] if isinstance(x, dict)]
            if isinstance(data.get("data"), dict) and isinstance(data["data"].get("list"), list):
                return [x for x in data["data"]["list"] if isinstance(x, dict)]
        return []
