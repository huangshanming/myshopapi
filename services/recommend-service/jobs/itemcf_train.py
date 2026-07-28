#!/usr/bin/env python3
"""Offline ItemCF / Swing job stub.

Usage (after implementing):
  cd services/recommend-service
  source .venv/bin/activate
  python -m jobs.itemcf_train
"""

from __future__ import annotations
from collections import defaultdict
import redis
import pymysql
from datetime import datetime, timedelta
from app.config import get_settings
import numpy as np

# from app.clients import redis_client
mysql_conn = pymysql.connect(**get_settings().mysql_connect_kwargs())

TOP_SIM_NUM = 50
DAY_RANGE = 90
# 置信度阈值：共同交互人数少于这个值，丢弃该商品对（避免噪声）
MIN_CO_OCCUR_USER = 2
    
 # 行为权重
BEHAVIOR_WEIGHT = {
    # "view": 1,
    # "click": 3,
    "cart": 8,
    "buy": 15,
    "collect": 4,
}
settings = get_settings()
r = redis.from_url(settings.redis_url, decode_responses=True)
def main() -> None:
    # TODO:
    # 1. load interactions from MySQL / ClickHouse
    # 2. build co-occurrence / Swing scores
    # 3. write Top-N neighbors to Redis ZSET (see app.reco.keys)
    # redis = redis_client.get_redis()
    # 加载用户行为数据
    user_behavior = load_user_behavior()
    sim_data = calc_itemcf(user_behavior)
    save_to_redis(sim_data)
    # 计算用户行为相似度
    mysql_conn.close()

def load_user_behavior():
    if not mysql_conn:
        raise ValueError("MySQL connection not established")

    # 获取用户行为数据
    with mysql_conn.cursor() as cursor:
        cursor.execute("""
        SELECT i.product_id, o.user_id  FROM orders o left join order_items i on o.id = i.order_id 
        WHERE o.status IN ('confirmed','shipped','completed','reviewed') and o.created_at >= %s
        GROUP BY i.product_id, o.user_id
        """, (datetime.now() - timedelta(days=DAY_RANGE)).strftime("%Y-%m-%d"))
        order_data =cursor.fetchall() # 用户下单数据

        cursor.execute("""
        SELECT product_id, user_id  FROM product_favorites
        WHERE created_at >= %s
        GROUP BY product_id, user_id
        """, (datetime.now() - timedelta(days=DAY_RANGE)).strftime("%Y-%m-%d"))
        favorite_data =cursor.fetchall() # 用户收藏数据
    user_behavior = {}
    for product_id, user_id in order_data:
        if user_id not in user_behavior:
            user_behavior[user_id] = {}
        if product_id not in user_behavior[user_id]:
            user_behavior[user_id][product_id] = 0
        user_behavior[user_id][product_id] += BEHAVIOR_WEIGHT["buy"]

    for product_id, user_id in favorite_data:
        if user_id not in user_behavior or product_id not in user_behavior[user_id]:
            user_behavior[user_id][product_id] = 0
        user_behavior[user_id][product_id] += BEHAVIOR_WEIGHT["collect"]
    return user_behavior

def calc_itemcf(user_behavior):
     # co_occur[a,b] = 累计共现权重
    co_occur = defaultdict(float)
    # item_pop[a] = 商品总权重（用户交互权重累加）
    item_pop = defaultdict(float)
    # item_co_user_count：记录两个商品共同交互的用户数量（用于置信度过滤）
    item_co_user = defaultdict(int)
    for user_id, item_dict in user_behavior.items():
        items = list(item_dict.items())
        n = len(items)
        for a, wa in items:
            for i in range(n):
                item_pop[a] += wa
                for j in range(i + 1, n):
                    b, wb = items[j]
                    score = min(wa, wb)
                    co_occur[(a, b)] += score
                    co_occur[(b, a)] += score
                    item_co_user[(a, b)] += 1
                    item_co_user[(b, a)] += 1
    # 构建每个商品的相似列表
    sim_raw = defaultdict(list)
    all_item_ids = list(item_pop.keys())

    # 遍历所有共现对，计算相似度（不再遍历全部商品笛卡尔积！！！）
    for (a, b), total_weight in co_occur.items():
        # 置信度过滤：共同用户太少直接跳过
        if item_co_user[(a, b)] < MIN_CO_OCCUR_USER:
            continue
        # 余弦相似度
        denom = np.sqrt(item_pop[a] * item_pop[b])
        if denom <= 0:
            continue
        sim_score = total_weight / denom
        sim_raw[a].append((b, sim_score))

    # 排序取topN
    sim_result = {}
    for item_id, sim_list in sim_raw.items():
        sim_list.sort(key=lambda x: x[1], reverse=True)
        top = [iid for iid, score in sim_list[:TOP_SIM_NUM]]
        sim_result[item_id] = top

    return sim_result

def save_to_redis(sim_dict):
    prefix = "sim:spu:"
    count = 0
    for item_id, similar_list in sim_dict.items():
        key = f"{prefix}{item_id}"
        r.delete(key)
        r.rpush(key, *similar_list)
        count += 1
    print(f"✅ 相似度数据写入Redis完成，共处理 {count} 个商品")
if __name__ == "__main__":
    main()
