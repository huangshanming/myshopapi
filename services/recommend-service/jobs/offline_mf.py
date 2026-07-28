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
from implicit.als import AlternatingLeastSquares
from app.config import get_settings
import pandas as pd
from pymilvus import connections, FieldSchema, DataType, Collection, CollectionSchema, utility
import numpy as np
DAY_RANGE = 90  # 行为数据时间范围
VECTOR_DIM = 64   # 隐向量维度
 # 行为权重
BEHAVIOR_WEIGHT = {
    # "view": 1,
    # "click": 3,
    "cart": 8,
    "buy": 15,
    "collect": 4,
}

r = redis.from_url(get_settings().redis_url, decode_responses=True)
def main() -> None:
    pass
    mysql_conn.close()
    connections.disconnect("default")

# 构建/加载商品向量集合
def init_milvus_collection():
    collection_name = "item_mf_vector"
    if utility.has_collection(collection_name):
        utility.drop_collection(collection_name)  
    fields = [
        FieldSchema(name="product_id", dtype=DataType.INT64, max_length=64, is_primary=True),
        FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=VECTOR_DIM)
    ]
    schema = CollectionSchema(fields=fields, description="MF item vectors")
    coll = Collection(name=collection_name, schema=schema)
    # 如果不存在索引，创建索引
    if not coll.has_index():
        index_params = {
            "index_type": "FLAT",
            "metric_type": "IP", # 内积，匹配MF
            "params":{}
        }
        coll.create_index(field_name="vector", index_params=index_params)
    coll.load()
    return coll

  
  
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

    raw_data = []
    for user_id, products in user_behavior.items():
        for product_id, weight in products.items():
          raw_data.append([user_id, product_id, weight])

    df = pd.DataFrame(raw_data, columns=["user_id", "product_id", "weight"])
    user2id = {uid:i for i,uid in enumerate(df["user_id"].unique())}
    item2id = {pid:i for i,pid in enumerate(df["product_id"].unique())}
    id2item = {v: k for k, v in item2id.items()}
    df["uid"] = df["user_id"].map(user2id) # 给表格新增一列数字版用户ID
    df["iid"] = df["product_id"].map(item2id)  # 给表格新增一列数字版商品ID
    return df, user2id, item2id, id2item

def train_als():
    df, user2id, item2id, id2item = load_user_behavior() # 拿到处理好的数据
    
    # 初始化ALS模型：64维向量，正则化系数0.01（防过拟合），迭代15次
    model = AlternatingLeastSquares(factors=VECTOR_DIM, regularization=0.01, iterations=15)
    
    from scipy.sparse import coo_matrix
    # 把表格变成一个巨大的“稀疏矩阵”（行是用户，列是商品，值是权重）
    mat = coo_matrix((df["weight"], (df["uid"], df["iid"])))
    model.fit(mat)  # 【核心】开始训练！让模型学习用户和商品的隐藏特征

    item_vectors = model.item_factors  # 训练完，拿到所有商品的向量
    user_vectors = model.user_factors  # 训练完，拿到所有用户的向量

    # -------- 处理商品向量（存入Milvus） --------
    coll = init_milvus_collection()      # 获取Milvus表
    coll.delete(expr="product_id != 0")     # 清空表里的旧数据（因为每次训练都是全量更新）
    
    insert_data = []
    for iid, product_id in id2item.items():  # 遍历每个商品
        vec = item_vectors[iid].tolist() # 把numpy数组转成Python列表
        insert_data.append([product_id, vec])# 打包成 [商品ID, 向量]
        
    sku_list, vec_list = zip(*insert_data) # 把ID和向量拆分成两个列表
    coll.insert([list(sku_list), list(vec_list)]) # 批量写入Milvus

    # -------- 处理用户向量（存入Redis） --------
    for user_id, uid in user2id.items(): # 遍历每个用户
        vec_arr = user_vectors[uid].tolist() # 取出用户的向量
        # 以 `mf:user_vec:用户ID` 为键，向量字符串为值，存入Redis
        r.set(f"mf:user_vec:{user_id}", str(vec_arr)) 

    print("✅ ALS训练完成，向量入库成功")

if __name__ == "__main__":
    # from app.clients import redis_client
    mysql_conn = pymysql.connect(**get_settings().mysql_connect_kwargs())
    connections.connect(alias="default", uri=get_settings().milvus_uri)
    init_milvus_collection()
    train_als()  # 运行主程序
    mysql_conn.close()
    connections.disconnect("default")
