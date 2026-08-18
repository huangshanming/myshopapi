# recommend-service — Collaborative Recommend (FastAPI)

协同推荐服务骨架：ItemCF 离线算图 + Redis 在线召回。业务逻辑在 `app/reco/service.py` / `jobs/` 自行实现。

## 运行

需要 **Python 3.10+**。

```bash
cd services/recommend-service
python3 -m venv .venv
source .venv/bin/activate
pip install -U pip
pip install -r requirements.txt
cp .env.example .env
uvicorn app.main:app --reload --host 0.0.0.0 --port 8888
```

环境变量与其它服务对齐：`MYMALL_MYSQL_*`、`MYMALL_REDIS_*`、`MYMALL_MILVUS_*`
（本地默认 MySQL `homestead/secret@127.0.0.1/mymall`，Milvus `127.0.0.1:19530`）。

### Milvus

```bash
# 仅基础设施（含 Milvus standalone）
docker compose -f deploy/local/docker-compose.infra.yaml up -d milvus-etcd milvus-minio milvus

# 健康检查
curl -s http://127.0.0.1:19091/healthz
```

Job / 代码里：

```python
from app.clients import milvus_client

milvus_client.connect()
milvus_client.ensure_item_collection()  # 创建 reco_item_emb
# milvus_client.upsert_item_embeddings(ids, vectors)
# milvus_client.search_similar_items(query_vec, top_k=10)
```

Job 里可用：

```python
from app.config import get_settings
import pymysql

settings = get_settings()
conn = pymysql.connect(**settings.mysql_connect_kwargs(), cursorclass=pymysql.cursors.DictCursor)
```

| 接口 | 说明 |
|------|------|
| `GET /healthz` | 存活 |
| `GET /readyz` | 就绪 |
| `GET /api/v1/recommend/also-bought?product_id=&limit=` | 买了又买（空实现） |
| `GET /api/v1/recommend/for-you?limit=` | 猜你喜欢（空实现，可读 `X-User-Id`） |
| `POST /api/v1/recommend/track` | 行为上报（仅计数返回） |

mall-uni 本地已代理 `/api/v1/recommend` → `:8888`。

## ItemCF 测试造数

先有店铺种子，再跑造数脚本（300 用户 / 120 商品 / 2000 订单，带购买簇共现）：

```bash
mysql -uhomestead -psecret mymall < scripts/seed-admin-merchant.sql

cd services/recommend-service && source .venv/bin/activate
pip install pymysql
python ../../scripts/seed_itemcf_data.py --clean
```

- 用户手机：`13800001000`～`13800001299`，密码 `123456`
- 商品编号：`ICF-P-000`…
- 订单号：`ICF-O-00000`…
- 交互 CSV：`scripts/data/itemcf_interactions.csv`

验收：

```sql
SELECT COUNT(*) FROM users WHERE mobile LIKE '13800001%';       -- 300
SELECT COUNT(*) FROM products WHERE product_no LIKE 'ICF-P-%';  -- 120
SELECT COUNT(*) FROM orders WHERE order_no LIKE 'ICF-O-%';     -- 2000
```

## 离线 Job

```bash
cd services/recommend-service && source .venv/bin/activate

# A 路 ItemCF → Redis sim:spu:{id}
python -m jobs.itemcf_train

# B 路 ALS/MF → Milvus item 向量 + Redis 用户向量
python -m jobs.offline_mf

# C 路 SASRec → Redis sasrec:next:{user_id} + Milvus item_sasrec_vector + artifacts/sasrec.pt
python -m jobs.offline_sasrec
```

SASRec 需要 `torch`（已写入 `requirements.txt`），以及本机 MySQL / Redis / Milvus 可用。

## 目录

```
app/
  main.py
  config.py
  api/routes/       # health / recommend / track
  reco/             # 业务逻辑（keys + service stubs）
  clients/          # redis / catalog / milvus
jobs/
  itemcf_train.py   # ItemCF 离线
  offline_mf.py     # ALS/MF 离线
  offline_sasrec.py # SASRec 序列推荐离线
artifacts/          # sasrec.pt 等模型产物
etc/
  recommend-service.yaml
```
