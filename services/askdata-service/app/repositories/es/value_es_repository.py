"""
字段取值 ES 仓储

把字段真实取值组织成 Elasticsearch 全文索引，并提供索引创建 批量写入和关键词检索能力

Service 层负责决定哪些字段需要同步
Repository 只关心索引是否存在 ValueInfo 如何写进 ES 以及如何按关键词召回
"""

import asyncio
from dataclasses import asdict

from elastic_transport import ConnectionTimeout
from elasticsearch import AsyncElasticsearch

from app.core.log import logger
from app.entities.value_info import ValueInfo


class ValueESRepository:
    """负责字段取值全文索引的创建 写入和基础检索"""

    index_name = "value_index"
    # value 字段使用 IK 分词，这样地区 会员等级 品类等中文值才能按全文方式检索
    index_mappings = {
        "dynamic": False,
        "properties": {
            "id": {"type": "keyword"},
            "value": {
                "type": "text",
                "analyzer": "ik_max_word",
                "search_analyzer": "ik_max_word",
            },
            "column_id": {"type": "keyword"},
        },
    }

    def __init__(self, client: AsyncElasticsearch):
        self.client = client

    async def ensure_index(self):
        """确保字段取值索引已经创建好"""
        if not await self.client.indices.exists(index=self.index_name):
            await self.client.indices.create(
                index=self.index_name, mappings=self.index_mappings
            )

    async def recreate_index(self):
        """删除并重建取值索引，避免全量重建后残留已废弃取值"""
        if await self.client.indices.exists(index=self.index_name):
            await self.client.indices.delete(index=self.index_name)
        await self.client.indices.create(
            index=self.index_name, mappings=self.index_mappings
        )

    async def index(
        self,
        value_infos: list[ValueInfo],
        batch_size: int = 100,
        max_retries: int = 3,
    ):
        """分批写入字段取值；超时自动重试，避免偶发 ConnectionTimeout 打断全量构建"""
        if not value_infos:
            return

        total = len(value_infos)
        for i in range(0, total, batch_size):
            batch_value_infos = value_infos[i : i + batch_size]
            batch_operations = []
            for value_info in batch_value_infos:
                # 用 ValueInfo.id 作为文档 id，这样重复构建时会覆盖同一条值记录
                batch_operations.append(
                    {"index": {"_index": self.index_name, "_id": value_info.id}}
                )
                batch_operations.append(asdict(value_info))

            last_error: Exception | None = None
            for attempt in range(max_retries):
                try:
                    await self.client.bulk(
                        operations=batch_operations, request_timeout=120
                    )
                    last_error = None
                    break
                except ConnectionTimeout as exc:
                    last_error = exc
                    wait_seconds = 2**attempt
                    logger.warning(
                        f"ES bulk 超时，{wait_seconds}s 后重试 "
                        f"({attempt + 1}/{max_retries})，进度 {i}/{total}"
                    )
                    await asyncio.sleep(wait_seconds)

            if last_error is not None:
                raise last_error

            if (i // batch_size) % 20 == 0:
                logger.info(f"ES 取值索引写入进度：{min(i + batch_size, total)}/{total}")

    async def search(
        self, keyword: str, score_threshold: float = 0.6, limit: int = 20
    ) -> list[ValueInfo]:
        """按关键词全文检索字段取值，并还原为 ValueInfo 实体"""

        resp = await self.client.search(
            index=self.index_name,
            # value 字段启用了 IK 分词，match 查询可以处理中文短语和枚举值匹配
            query={"match": {"value": keyword}},
            size=limit,
            # 过滤掉相关度过低的命中，避免把明显无关的取值带入后续上下文
            min_score=score_threshold,
        )
        # ES 文档 _source 中保存的是 ValueInfo 的字段结构，业务层继续使用实体对象
        return [ValueInfo(**hit["_source"]) for hit in resp["hits"]["hits"]]
