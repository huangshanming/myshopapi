from typing import TypedDict

from langchain_huggingface import HuggingFaceEndpointEmbeddings

from app.repositories.es.value_es_repository import ValueESRepository
from app.repositories.mysql.meta.meta_mysql_repository import MetaMySQLRepository
from app.repositories.qdrant.column_qdrant_repository import ColumnQdrantRepository
from app.repositories.qdrant.metric_qdrant_repository import MetricQdrantRepository

class DataAgentContext(TypedDict):
    """LangGraph Runtime 中传递的上下文对象"""
    # 字段向量仓储，负责根据向量从 Qdrant 检索候选字段
    column_qdrant_repository: ColumnQdrantRepository
    # Embedding 客户端，字段召回和指标召回都会复用
    embedding_client: HuggingFaceEndpointEmbeddings
    # 指标向量仓储，负责从 Qdrant 检索候选指标
    metric_qdrant_repository: MetricQdrantRepository
    # 字段取值仓储，负责从 Elasticsearch 检索真实字段值
    value_es_repository: ValueESRepository
    # 元数据库仓储，合并阶段用它按 id 补齐字段、表、主外键信息
    meta_mysql_repository: MetaMySQLRepository