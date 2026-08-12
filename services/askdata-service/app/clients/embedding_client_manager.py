import asyncio
from typing import Any, Optional

from huggingface_hub import AsyncInferenceClient

from app.conf.app_config import EmbeddingConfig, app_config


class TeiEmbeddings:
    """本地 TEI（Text Embeddings Inference）客户端，兼容 LangChain 的 aembed_* 用法。

    TEI 原生接口是 POST /embed，不是 HuggingFace Hub 上的模型 ID，
    因此不能用 HuggingFaceEndpointEmbeddings(model=repo_id)。
    """

    def __init__(self, base_url: str, model: str):
        self._client = AsyncInferenceClient(base_url=base_url)
        self.model = model

    @staticmethod
    def _to_vector(raw: Any) -> list[float]:
        # feature_extraction 可能返回 list / ndarray / 嵌套 batch 结果
        if hasattr(raw, "tolist"):
            raw = raw.tolist()
        if isinstance(raw, list) and raw and isinstance(raw[0], list):
            raw = raw[0]
        return [float(x) for x in raw]

    async def aembed_query(self, text: str) -> list[float]:
        raw = await self._client.feature_extraction(text)
        return self._to_vector(raw)

    async def aembed_documents(self, texts: list[str]) -> list[list[float]]:
        # TEI / HF client 对单条更稳；批量时逐条请求，避免输入格式差异
        return [await self.aembed_query(t) for t in texts]


class EmbeddingClientManager:
    def __init__(self, config: EmbeddingConfig):
        # 客户端在模块导入阶段先不立即创建，避免启动时就发起外部依赖连接
        self.client: Optional[TeiEmbeddings] = None
        self.config = config
        self.base_url = f"http://{self.config.host}:{self.config.port}"

    def init(self):
        # 在应用启动阶段显式调用，完成真正的客户端初始化
        self.client = TeiEmbeddings(base_url=self.base_url, model=self.config.model)


# 模块级单例，供其他模块按需复用同一个客户端管理器
embedding_client_manager = EmbeddingClientManager(app_config.embedding)


if __name__ == "__main__":
    embedding_client_manager.init()
    client = embedding_client_manager.client
    assert client is not None

    async def test():
        text = "What is deep learning?"
        query_result = await client.aembed_query(text)
        print(len(query_result), query_result[:3])

    asyncio.run(test())
