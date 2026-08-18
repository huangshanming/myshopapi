#!/usr/bin/env python3
"""Offline SASRec: sequence next-item → Redis TopK + Milvus item emb + checkpoint.

与 ItemCF / ALS 的差别：必须保留「时间顺序」的用户行为序列，学「下一步买什么」。

Usage:
  cd services/recommend-service
  source .venv/bin/activate
  python -m jobs.offline_sasrec
"""

# 延后求值类型注解，避免循环引用、写法更干净
from __future__ import annotations

# defaultdict：按 user 聚事件时少写 if key not in dict
from collections import defaultdict
from datetime import datetime, timedelta
# 拼 artifacts/sasrec.pt 的绝对路径
from pathlib import Path
from typing import Any

# L2 归一化 item 向量（适配 Milvus 内积 IP）
import numpy as np
# 读 MySQL 订单 / 收藏
import pymysql
# 写用户下一跳 TopK
import redis
# 张量计算
import torch
import torch.nn as nn
# Dataset：一条样本；DataLoader：按 batch 喂模型
from torch.utils.data import DataLoader, Dataset

# 复用已验证的 connect / ensure / upsert（host 用 milvus_host，不要传 uri）
from app.clients import milvus_client
from app.config import get_settings

# ========================= 超参 =========================
DAY_RANGE = 90          # 只取最近 90 天交互
MAX_LEN = 50            # 序列最大长度；更长则截最近 MAX_LEN
EMBED_DIM = 64          # item / 位置向量维度 d
NUM_BLOCKS = 2          # TransformerEncoder 堆叠层数
NUM_HEADS = 2           # 多头注意力头数（须整除 EMBED_DIM）
DROPOUT = 0.2           # 防过拟合
BATCH_SIZE = 128        # 每批用户序列数
EPOCHS = 20             # 训练轮数
LR = 1e-3               # Adam 学习率
TOP_K = 50              # 写入 Redis 的下一跳候选数
PAD_ID = 0              # padding 专用 id；真实 item 从 1 起编
REDIS_PREFIX = "sasrec:next:"                 # Redis key 前缀
MILVUS_COLLECTION = "item_sasrec_vector"      # 商品向量集合名
# __file__=jobs/offline_sasrec.py → parent=jobs → parent.parent=recommend-service
ARTIFACT = Path(__file__).resolve().parent.parent / "artifacts" / "sasrec.pt"


def pad_left(seq: list[int], length: int, pad: int = PAD_ID) -> list[int]:
    """左侧补 pad，使长度 == length；过长则保留右边（最近）length 个。"""
    if len(seq) >= length:
        return seq[-length:]
    return [pad] * (length - len(seq)) + list(seq)


def load_sequences(conn: pymysql.Connection) -> dict[int, list[int]]:
    """
    从 MySQL 拉「有时间戳」的行为，聚成 user → [product_id, ...]（升序）。
    关键：不 GROUP BY 成权重矩阵（那是 MF）；必须保留顺序，SASRec 才有意义。
    """
    # 窗口起点：例如今天往前 90 天
    since = (datetime.now() - timedelta(days=DAY_RANGE)).strftime("%Y-%m-%d")
    rows: list[tuple[Any, Any, Any]] = []
    with conn.cursor() as cur:
        # 订单行：用订单创建时间当事件时间
        cur.execute(
            """
            SELECT o.user_id, i.product_id, o.created_at
            FROM orders o
            JOIN order_items i ON o.id = i.order_id
            WHERE o.status IN ('confirmed','shipped','completed','reviewed')
              AND o.created_at >= %s
              AND i.product_id IS NOT NULL
            """,
            (since,),
        )
        rows.extend(cur.fetchall())
        # 收藏：同样带时间，稍后与订单混排
        cur.execute(
            """
            SELECT user_id, product_id, created_at
            FROM product_favorites
            WHERE created_at >= %s
            """,
            (since,),
        )
        rows.extend(cur.fetchall())

    # user_id → [(ts, product_id), ...]
    by_user: dict[int, list[tuple[Any, int]]] = defaultdict(list)
    for user_id, product_id, ts in rows:
        by_user[int(user_id)].append((ts, int(product_id)))

    user_seqs: dict[int, list[int]] = {}
    for user_id, events in by_user.items():
        # 先按时间，同秒再按 product_id，保证可复现
        events.sort(key=lambda x: (x[0], x[1]))
        seq: list[int] = []
        for _, pid in events:
            # 连续重复同一商品丢掉（噪声）
            if seq and seq[-1] == pid:
                continue
            seq.append(pid)
        # 长度 < 2 无法构造「下一个」监督信号
        if len(seq) >= 2:
            user_seqs[user_id] = seq
    return user_seqs


def build_item_vocab(
    user_seqs: dict[int, list[int]],
) -> tuple[dict[int, int], dict[int, int]]:
    """真实 product_id ↔ 模型内部连续 id（1..N）；0 留给 PAD。"""
    # 集合推导：所有出现过的商品，排序保证每次训练映射稳定
    items = sorted({pid for seq in user_seqs.values() for pid in seq})
    item2id = {pid: i + 1 for i, pid in enumerate(items)}  # 从 1 开始
    id2item = {i: pid for pid, i in item2id.items()}       # 反查，导出 Redis 要用
    return item2id, id2item


class SeqDataset(Dataset):
    """把每条用户序列变成 (输入 x, 标签 y)，供 DataLoader 取 batch。"""

    def __init__(self, user_seqs: dict[int, list[int]], item2id: dict[int, int]) -> None:
        self.samples: list[tuple[list[int], list[int]]] = []
        for seq in user_seqs.values():
            # 真实 pid → 内部 id
            ids = [item2id[p] for p in seq]
            if len(ids) < 2:
                continue
            # 只留最近 MAX_LEN 个行为
            ids = ids[-MAX_LEN:]
            # next-item：用前缀预测下一个
            # 例 ids=[a,b,c] → inp=[a,b] tgt=[b,c]
            inp = ids[:-1]
            tgt = ids[1:]
            # 左 pad 到定长，方便组成 batch 矩阵 [B, MAX_LEN]
            x = pad_left(inp, MAX_LEN)
            y = pad_left(tgt, MAX_LEN)
            self.samples.append((x, y))

    def __len__(self) -> int:
        return len(self.samples)

    def __getitem__(self, idx: int) -> tuple[torch.Tensor, torch.Tensor]:
        # DataLoader 每次取下标 idx；long = 整数索引，给 Embedding 用
        x, y = self.samples[idx]
        return torch.tensor(x, dtype=torch.long), torch.tensor(y, dtype=torch.long)


class SASRec(nn.Module):
    """
    SASRec 风格：因果 Self-Attention 编码序列。
    每个时间步输出一个隐向量 h，再与全体 item embedding 点积 → 下一跳打分。
    """

    def __init__(self, num_items: int) -> None:
        super().__init__()
        # 词表大小 = 商品数 + 1（pad）；padding_idx 使 pad 向量恒为 0、无梯度
        self.item_emb = nn.Embedding(num_items + 1, EMBED_DIM, padding_idx=PAD_ID)
        # 位置 0..MAX_LEN-1 的可学习位置编码（Transformer 本身无顺序，靠它注入）
        self.pos_emb = nn.Embedding(MAX_LEN, EMBED_DIM)
        # 一层：多头自注意力 + FFN；batch_first → 张量形状 [B, L, d]
        encoder_layer = nn.TransformerEncoderLayer(
            d_model=EMBED_DIM,
            nhead=NUM_HEADS,
            dim_feedforward=EMBED_DIM * 4,
            dropout=DROPOUT,
            batch_first=True,
            activation="gelu",
        )
        # 堆 NUM_BLOCKS 层
        self.encoder = nn.TransformerEncoder(encoder_layer, num_layers=NUM_BLOCKS)
        self.dropout = nn.Dropout(DROPOUT)
        self._init_weights()

    def _init_weights(self) -> None:
        # 小方差正态初始化，训练更稳
        nn.init.normal_(self.item_emb.weight, std=0.02)
        nn.init.normal_(self.pos_emb.weight, std=0.02)
        # 显式把 pad 行置零（normal_ 会覆盖掉 Embedding 默认的 pad 零向量）
        with torch.no_grad():
            self.item_emb.weight[PAD_ID].zero_()

    def forward(self, seq: torch.Tensor) -> torch.Tensor:
        """
        seq: [B, L] 内部 item id
        return logits: [B, L, N+1] 每个位置对所有商品（含 pad 行）的分数
        """
        _b, length = seq.shape
        # 位置下标 [0,1,...,L-1]，扩成与 batch 同形 [B, L]
        pos = torch.arange(length, device=seq.device).unsqueeze(0).expand(seq.size(0), -1)
        # 内容向量 × √d + 位置向量（论文常用缩放）；再 dropout
        x = self.item_emb(seq) * (EMBED_DIM**0.5) + self.pos_emb(pos)
        x = self.dropout(x)
        # 因果 mask：上三角为 -inf → 位置 i 不能看 j>i（不能偷看未来）
        # 只用 float mask，不用 src_key_padding_mask：避免「整行被 mask」导致 softmax→NaN
        causal = torch.triu(
            torch.full((length, length), float("-inf"), device=seq.device),
            diagonal=1,
        )
        # h: [B, L, d] 每个位置的上下文表示
        h = self.encoder(x, mask=causal)
        # 与全体 item 向量点积打分（含第 0 行 pad；loss 会 ignore）
        logits = h @ self.item_emb.weight.T  # [B, L, N+1]
        return logits


def train_model(model: SASRec, loader: DataLoader, device: torch.device) -> None:
    """标准监督学习循环：前向 → CE loss → 反传 → 更新。"""
    opt = torch.optim.Adam(model.parameters(), lr=LR)
    # ignore_index=PAD_ID：y 里是 pad 的位置不算进 loss
    loss_fn = nn.CrossEntropyLoss(ignore_index=PAD_ID)
    model.train()  # 打开 dropout
    for epoch in range(EPOCHS):
        total = 0.0
        n_batches = 0
        for x, y in loader:
            # 搬到 GPU/CPU
            x, y = x.to(device), y.to(device)
            logits = model(x)  # [B, L, N+1]
            # 摊平时间步：当成「每个位置一个多分类」
            loss = loss_fn(logits.reshape(-1, logits.size(-1)), y.reshape(-1))
            if not torch.isfinite(loss):
                raise RuntimeError(f"non-finite loss={loss.item()}")
            opt.zero_grad()          # 清上一步梯度
            loss.backward()          # 反传
            nn.utils.clip_grad_norm_(model.parameters(), 5.0)  # 防梯度爆炸
            opt.step()               # 更新参数
            total += float(loss.item())
            n_batches += 1
        avg = total / max(n_batches, 1)
        print(f"epoch {epoch + 1}/{EPOCHS} loss={avg:.4f}")


@torch.no_grad()  # 导出阶段不需要梯度，省显存/加速
def export_artifacts(
    model: SASRec,
    user_seqs: dict[int, list[int]],
    item2id: dict[int, int],
    id2item: dict[int, int],
    r: redis.Redis,
    device: torch.device,
) -> None:
    """三份产出：Milvus 商品向量 / Redis 用户 TopK / 本地 checkpoint。"""
    model.eval()  # 关掉 dropout
    num_items = len(item2id)

    # ---------- 1) item embedding → Milvus ----------
    # 跳过第 0 行 pad；detach 脱离计算图；转 numpy float32
    emb = model.item_emb.weight[1 : num_items + 1].detach().cpu().numpy().astype(np.float32)
    if not np.isfinite(emb).all():
        emb = np.nan_to_num(emb, nan=0.0, posinf=0.0, neginf=0.0)
    # L2 归一化后，内积 ≈ 余弦相似度
    norms = np.linalg.norm(emb, axis=1, keepdims=True).clip(min=1e-8)
    emb = emb / norms
    # 内部 id 1..N → 真实 product_id
    pids = [int(id2item[i]) for i in range(1, num_items + 1)]
    # schema: id INT64 + embedding FLOAT_VECTOR（见 milvus_client）
    milvus_client.ensure_item_collection(collection_name=MILVUS_COLLECTION, dim=EMBED_DIM)
    n_up = milvus_client.upsert_item_embeddings(
        pids, emb.tolist(), collection_name=MILVUS_COLLECTION
    )
    print(f"Milvus upsert {n_up} vectors → {MILVUS_COLLECTION}")

    # ---------- 2) 每用户「下一跳」TopK → Redis ----------
    exported = 0
    sample_key = ""
    sample_vals: list[str] = []
    for user_id, seq in user_seqs.items():
        # 用完整历史（截断）做一次前向
        ids = [item2id[p] for p in seq if p in item2id][-MAX_LEN:]
        if len(ids) < 1:
            continue
        # batch 维 = 1
        x = torch.tensor([pad_left(ids, MAX_LEN)], dtype=torch.long, device=device)
        logits = model(x)  # [1, L, N+1]
        # 找最后一个非 pad 位置（序列真实末尾）
        nonzero = (x[0] != PAD_ID).nonzero(as_tuple=False)
        if nonzero.numel() == 0:
            continue
        last = int(nonzero[-1].item())
        scores = logits[0, last].clone()
        scores[PAD_ID] = float("-inf")  # 永不推荐 pad
        # 已看过的商品降权，减少「再推一遍刚买的」
        seen = set(ids)
        for iid in seen:
            scores[iid] -= 1e3
        k = min(TOP_K, num_items)
        top_idx = torch.topk(scores, k=k).indices.tolist()
        # 内部 id → 真实 product_id，在线服务直接用
        top_pids = [id2item[i] for i in top_idx if i in id2item]
        key = f"{REDIS_PREFIX}{user_id}"
        r.delete(key)  # 全量覆盖旧结果
        if top_pids:
            r.rpush(key, *[str(pid) for pid in top_pids])
            exported += 1
            if not sample_key:
                sample_key = key
                sample_vals = [str(pid) for pid in top_pids[:5]]

    print(f"Redis wrote {exported} keys ({REDIS_PREFIX}*)")
    if sample_key:
        print(f"  sample {sample_key} → {sample_vals}...")

    # ---------- 3) checkpoint（以后可挂在线 Torch 推理）----------
    ARTIFACT.parent.mkdir(parents=True, exist_ok=True)
    torch.save(
        {
            "state_dict": model.state_dict(),  # 权重
            "item2id": item2id,              # 词典，推理必须一致
            "id2item": id2item,
            "hyperparams": {
                "max_len": MAX_LEN,
                "embed_dim": EMBED_DIM,
                "num_blocks": NUM_BLOCKS,
                "num_heads": NUM_HEADS,
                "dropout": DROPOUT,
                "top_k": TOP_K,
            },
        },
        ARTIFACT,
    )
    print(f"checkpoint → {ARTIFACT}")


def main() -> None:
    """入口：连接 → 造序列 → 训练 → 导出；连接放这里，避免 import 就连库。"""
    settings = get_settings()
    conn = pymysql.connect(**settings.mysql_connect_kwargs())
    r = redis.from_url(settings.redis_url, decode_responses=True)
    milvus_client.connect()
    try:
        user_seqs = load_sequences(conn)
        if not user_seqs:
            raise RuntimeError(
                f"no sequences in last {DAY_RANGE} days; check seed data / DAY_RANGE"
            )
        item2id, id2item = build_item_vocab(user_seqs)
        ds = SeqDataset(user_seqs, item2id)
        if len(ds) == 0:
            raise RuntimeError("SeqDataset empty after filtering")
        print(
            f"users={len(user_seqs)} items={len(item2id)} samples={len(ds)} "
            f"device_will_use={'cuda' if torch.cuda.is_available() else 'cpu'}"
        )
        # shuffle=True：打乱用户顺序，减轻 batch 偏差
        loader = DataLoader(ds, batch_size=BATCH_SIZE, shuffle=True)
        device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        model = SASRec(num_items=len(item2id)).to(device)
        train_model(model, loader, device)
        export_artifacts(model, user_seqs, item2id, id2item, r, device)
        print("✅ SASRec train done")
    finally:
        # 无论成功失败都关连接
        conn.close()
        milvus_client.disconnect()


# 仅直接 python -m jobs.offline_sasrec 时跑；被 import 时不自动训练
if __name__ == "__main__":
    main()
