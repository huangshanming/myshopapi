from __future__ import annotations

import httpx

from app.config import get_settings


class CatalogClient:
    """HTTP client to catalog-service — flesh out batch get as needed."""

    def __init__(self, base_url: str | None = None, timeout: float = 8.0) -> None:
        settings = get_settings()
        self.base_url = (base_url or settings.catalog_http).rstrip("/")
        self.timeout = timeout

    async def get_products_by_ids(self, product_ids: list[int]) -> list[dict]:
        # TODO: call catalog batch API when available
        _ = product_ids
        return []
