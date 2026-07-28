"""Milvus client helpers for recommend-service (sync; suitable for jobs + sync routes)."""

from __future__ import annotations

from typing import Sequence

from pymilvus import (
    Collection,
    CollectionSchema,
    DataType,
    FieldSchema,
    connections,
    utility,
)

from app.config import get_settings

_ALIAS = "default"
_connected = False


def connect(*, alias: str = _ALIAS) -> None:
    """Connect to Milvus (idempotent)."""
    global _connected
    if _connected and connections.has_connection(alias):
        return
    settings = get_settings()
    kwargs: dict = {
        "alias": alias,
        "host": settings.milvus_host,
        "port": str(settings.milvus_port),
        "db_name": settings.milvus_db_name,
    }
    if settings.milvus_user:
        kwargs["user"] = settings.milvus_user
        kwargs["password"] = settings.milvus_password
    connections.connect(**kwargs)
    _connected = True


def disconnect(*, alias: str = _ALIAS) -> None:
    global _connected
    if connections.has_connection(alias):
        connections.disconnect(alias)
    _connected = False


def ping() -> bool:
    """Return True if Milvus is reachable."""
    try:
        connect()
        utility.get_server_version()
        return True
    except Exception:
        return False


def ensure_item_collection(
    *,
    collection_name: str | None = None,
    dim: int | None = None,
    metric_type: str = "IP",
) -> Collection:
    """
    Create item embedding collection if missing.
    Schema: id (INT64 pk = product_id), embedding (FLOAT_VECTOR).
    """
    connect()
    settings = get_settings()
    name = collection_name or settings.milvus_item_collection
    dim = dim or settings.milvus_embedding_dim

    if utility.has_collection(name):
        col = Collection(name)
        col.load()
        return col

    fields = [
        FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=False),
        FieldSchema(name="embedding", dtype=DataType.FLOAT_VECTOR, dim=dim),
    ]
    schema = CollectionSchema(fields=fields, description="recommend item embeddings")
    col = Collection(name=name, schema=schema)
    index_params = {
        "metric_type": metric_type,
        "index_type": "IVF_FLAT",
        "params": {"nlist": 128},
    }
    col.create_index(field_name="embedding", index_params=index_params)
    col.load()
    return col


def upsert_item_embeddings(
    ids: Sequence[int],
    embeddings: Sequence[Sequence[float]],
    *,
    collection_name: str | None = None,
) -> int:
    """Upsert product vectors. Returns number of entities inserted/upserted."""
    if len(ids) != len(embeddings):
        raise ValueError("ids and embeddings length mismatch")
    if not ids:
        return 0
    col = ensure_item_collection(collection_name=collection_name)
    # delete existing then insert (compatible across pymilvus versions)
    expr = f"id in [{','.join(str(int(i)) for i in ids)}]"
    col.delete(expr)
    col.insert([list(ids), list(embeddings)])
    col.flush()
    return len(ids)


def search_similar_items(
    query_embedding: Sequence[float],
    *,
    top_k: int = 10,
    collection_name: str | None = None,
    exclude_ids: Sequence[int] | None = None,
) -> list[tuple[int, float]]:
    """
    ANN search. Returns [(product_id, score), ...] sorted by score desc (IP).
    """
    col = ensure_item_collection(collection_name=collection_name)
    settings = get_settings()
    search_params = {"metric_type": "IP", "params": {"nprobe": 16}}
    # over-fetch if we need to filter exclusions
    fetch_k = top_k + (len(exclude_ids) if exclude_ids else 0)
    fetch_k = max(fetch_k, top_k)
    results = col.search(
        data=[list(query_embedding)],
        anns_field="embedding",
        param=search_params,
        limit=fetch_k,
        output_fields=["id"],
    )
    exclude = set(int(x) for x in (exclude_ids or []))
    out: list[tuple[int, float]] = []
    if not results:
        return out
    for hit in results[0]:
        pid = int(hit.id)
        if pid in exclude:
            continue
        out.append((pid, float(hit.score)))
        if len(out) >= top_k:
            break
    _ = settings  # keep settings import useful for future filters
    return out
