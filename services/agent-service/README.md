# agent-service — Python AI Agents (FastAPI)

智能导购等 Agent 服务，与 Go 业务服务解耦；经 HTTP 调用 catalog / order / merchant。

## 运行

需要 **Python 3.10+**（推荐 3.10 / 3.11）。

```bash
cd services/agent-service
python3.10 -m venv .venv   # 或 python3 -m venv .venv（需 ≥3.10）
source .venv/bin/activate
pip install -U pip
pip install -r requirements.txt
cp .env.example .env
uvicorn app.main:app --reload --host 0.0.0.0 --port 8886
```

健康检查：`GET http://127.0.0.1:8886/healthz`  
导购规划（占位）：`POST http://127.0.0.1:8886/api/v1/agents/shopping-guide/plan`

本地 mall-uni 已代理 `/api/v1/agents` → `:8886`。

## 目录

```
app/
  main.py           # FastAPI 入口
  config.py         # 环境变量配置
  api/routes/       # HTTP 路由
  agents/           # Agent 实现（shopping_guide 等）
  clients/          # 调 Go 服务的 HTTP 客户端
etc/
  agent-service.yaml
```
