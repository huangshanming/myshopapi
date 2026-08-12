# askdata-service — 电商问数（FastAPI + uv）

自然语言查数 Agent。目录按「api / services / agent / repositories / clients」分层。

## 运行

需要 **Python 3.10–3.13**（推荐 3.13）。

```bash
cd services/askdata-service
uv sync
cp .env.example .env
uv run uvicorn app.main:app --reload --host 0.0.0.0 --port 8889
```

- 健康检查：`GET http://127.0.0.1:8889/healthz`
- 问数：`POST http://127.0.0.1:8889/api/v1/askdata/query`

```bash
curl -s http://127.0.0.1:8889/api/v1/askdata/query \
  -H 'Content-Type: application/json' \
  -d '{"question":"近7天 GMV 是多少？","days":7}'
```

## 目录结构

```
askdata-service/
├── app/                 # 后端源码
│   ├── agent/           # 问数 Agent / LangGraph 编排
│   ├── api/             # FastAPI 路由、schemas、依赖注入
│   ├── clients/         # MySQL / ES / Qdrant / Embedding 客户端
│   ├── conf/            # 配置加载（读根目录 conf/*.yaml）
│   ├── core/            # 日志、生命周期
│   ├── entities/        # 业务实体
│   ├── models/          # ORM 模型
│   ├── prompt/          # 加载 prompts/ 静态模板
│   ├── repositories/    # 数据访问层
│   ├── scripts/         # 初始化 / 知识库构建脚本
│   ├── services/        # 业务编排
│   └── main.py
├── conf/                # YAML：app / mysql / es / qdrant / llm / logging
├── docker/              # 本服务相关说明与 SQL（基建容器见 deploy/local）
├── logs/                # 本地日志输出
├── prompts/             # 静态 Prompt 模板
├── pyproject.toml
└── README.md
```

## 本地基建（Docker）

```bash
docker compose -f deploy/local/docker-compose.yaml up -d elasticsearch kibana qdrant
# Embedding 见 deploy/local/embedding/README.md
```

| 组件 | 地址 |
|------|------|
| Elasticsearch | http://127.0.0.1:9200 |
| Kibana | http://127.0.0.1:5601 |
| Qdrant | http://127.0.0.1:6333 |
| Embedding TEI | http://127.0.0.1:8082 |
| MySQL | 宿主机 :3306（homestead） |

## 与 agent-service

| | agent-service :8886 | askdata-service :8889 |
|--|---------------------|------------------------|
| 产品 | 智能导购 | 电商问数 |
| 典型能力 | 搜商品 | 指标问答 / NL2SQL |
