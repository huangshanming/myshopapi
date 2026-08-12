"""Qdrant 客户端占位。"""
from app.conf import get_config
import asyncio
import random
from typing import Optional

from qdrant_client import AsyncQdrantClient, models

from app.conf.app_config import QdrantConfig, app_config

class QdrantClientManager:
    def __init__(self, qdrant_config: QdrantConfig):
        # 保存配置对象，后面初始化客户端时要从这里读取 host 和 port
        self.qdrant_config = qdrant_config
        # 先把 client 声明出来，真正初始化放到 init() 中进行
        self.client: Optional[AsyncQdrantClient] = None

    def _get_url(self):
        # 根据配置文件拼出 Qdrant 服务地址
        return f"http://{self.qdrant_config.host}:{self.qdrant_config.port}"

    def init(self):
        # 创建异步客户端 
        # 这里不在 __init__ 中直接初始化，是为了和项目的生命周期管理保持一致
        self.client = AsyncQdrantClient(url=self._get_url())

    async def close(self):
        # 项目关闭时统一关闭客户端连接
        await self.client.close()


# 创建一个全局的管理器对象
# 后续项目中的其他模块都通过它来获取同一套 Qdrant 客户端
qdrant_client_manager = QdrantClientManager(app_config.qdrant)