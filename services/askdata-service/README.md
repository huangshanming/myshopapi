# askdata-service — 电商问数（FastAPI + uv）

自然语言查数 Agent：订单 GMV、单量、库存等指标问答。与导购 `agent-service` 解耦。

当前为 **规则 stub**（意图识别 + SQL 预览占位），后续可接 LangGraph / NL2SQL / order-service。

## 运行

需要 **Python 3.10+**（本机可用 3.13；由 uv 管理环境）。

```bash
cd services/askdata-service
uv sync
cp .env.example .env
uv run uvicorn app.main:app --reload --host 0.0.0.0 --port 8889
```

健康检查：`GET http://127.0.0.1:8889/healthz`  
问数（占位）：`POST http://127.0.0.1:8889/api/v1/askdata/query`

```bash
curl -s http://127.0.0.1:8889/api/v1/askdata/query \
  -H 'Content-Type: application/json' \
  -d '{"question":"近7天 GMV 是多少？","days":7}'
```

## 目录

```
app/
  main.py           # FastAPI 入口
  config.py         # 环境变量配置
  api/routes/       # health + askdata
  agents/           # AskDataAgent（问数逻辑）
  clients/          # 后续调 Go / 数仓的 HTTP 客户端
etc/
  askdata-service.yaml
pyproject.toml      # uv 项目与依赖
uv.lock
```

## 与 agent-service 的关系

| | agent-service :8886 | askdata-service :8889 |
|--|---------------------|------------------------|
| 产品 | 智能导购 | 电商问数 |
| 典型能力 | 搜商品、推荐清单 | 指标问答、NL2SQL |
| 依赖数据 | catalog 为主 | order / 指标 / 数仓为主 |
