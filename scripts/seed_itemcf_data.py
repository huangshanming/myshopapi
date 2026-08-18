#!/usr/bin/env python3
"""ItemCF 测试造数：用户 / 商品 / SKU / 订单 / 收藏 / 评价（带购买簇共现）。

规模（默认）
  - 300 用户  手机号 13800001000～13800001299  密码 123456
  - 120 商品  8 簇 × 15，product_no=ICF-P-000…
  - 2000 已完成订单（每单 2～4 件，80% 同簇）
  - ~800 收藏、~400 评价

依赖
  先执行 scripts/seed-admin-merchant.sql（需要「示例宠物店」）

用法
  pip install pymysql
  python scripts/seed_itemcf_data.py --clean

环境变量
  MYMALL_MYSQL_HOST / PORT / USERNAME / PASSWORD / DBNAME
  SEED_SHOP_ID  （可选，默认取「示例宠物店」）

验收 SQL 示例
  SELECT COUNT(*) FROM users WHERE mobile LIKE '13800001%';
  SELECT COUNT(*) FROM products WHERE product_no LIKE 'ICF-P-%';
  SELECT COUNT(*) FROM orders WHERE order_no LIKE 'ICF-O-%';
  -- 抽查共现：某商品同簇共购应高于跨簇
"""

from __future__ import annotations

import argparse
import csv
import hashlib
import os
import random
import sys
from collections import Counter
from datetime import datetime, timedelta
from pathlib import Path

try:
    import pymysql
except ImportError:
    print("请先安装: pip install pymysql", file=sys.stderr)
    sys.exit(1)

# ---- 规模 ----
N_USERS = 300
N_CLUSTERS = 8
PER_CLUSTER = 15
N_PRODUCTS = N_CLUSTERS * PER_CLUSTER  # 120
N_ORDERS = 2000
N_FAVORITES = 800
N_REVIEWS = 400
SAME_CLUSTER_RATIO = 0.80

MOBILE_START = 13800001000
PRODUCT_PREFIX = "ICF-P-"
SKU_PREFIX = "ICF-S-"
ORDER_PREFIX = "ICF-O-"
PASSWORD_PLAIN = "123456"
PASSWORD_SALT = "this is my mall"

CLUSTER_NAMES = [
    "猫粮精选",
    "猫砂清洁",
    "狗粮主食",
    "零食罐头",
    "玩具用品",
    "洗护美容",
    "出行牵引",
    "保健营养",
]

PET_TYPES = ["cat", "cat", "dog", "both", "both", "both", "dog", "both"]


def password_hash() -> str:
    return hashlib.md5(f"{PASSWORD_PLAIN}{PASSWORD_SALT}".encode()).hexdigest()


def env(name: str, default: str) -> str:
    return os.environ.get(name, default).strip() or default


def connect():
    return pymysql.connect(
        host=env("MYMALL_MYSQL_HOST", "127.0.0.1"),
        port=int(env("MYMALL_MYSQL_PORT", "3306")),
        user=env("MYMALL_MYSQL_USERNAME", "homestead"),
        password=env("MYMALL_MYSQL_PASSWORD", "secret"),
        database=env("MYMALL_MYSQL_DBNAME", "mymall"),
        charset="utf8mb4",
        autocommit=False,
        cursorclass=pymysql.cursors.DictCursor,
    )


def resolve_shop_id(cur) -> int:
    raw = env("SEED_SHOP_ID", "")
    if raw.isdigit():
        return int(raw)
    cur.execute("SELECT id FROM shops WHERE name=%s ORDER BY id ASC LIMIT 1", ("示例宠物店",))
    row = cur.fetchone()
    if not row:
        raise SystemExit(
            "未找到店铺「示例宠物店」。请先执行: mysql … < scripts/seed-admin-merchant.sql"
        )
    return int(row["id"])


def clean(cur) -> None:
    """按前缀删除本脚本写入的数据（不动 1390000000x 等真实种子）。"""
    print("cleaning previous ICF seed…")

    cur.execute(
        "SELECT id FROM orders WHERE order_no LIKE %s",
        (ORDER_PREFIX + "%",),
    )
    order_ids = [r["id"] for r in cur.fetchall()]
    if order_ids:
        ph = ",".join(["%s"] * len(order_ids))
        cur.execute(f"DELETE FROM product_reviews WHERE order_id IN ({ph})", order_ids)
        cur.execute(f"DELETE FROM order_items WHERE order_id IN ({ph})", order_ids)
        cur.execute(f"DELETE FROM orders WHERE id IN ({ph})", order_ids)

    cur.execute(
        "SELECT id FROM products WHERE product_no LIKE %s",
        (PRODUCT_PREFIX + "%",),
    )
    product_ids = [r["id"] for r in cur.fetchall()]
    if product_ids:
        ph = ",".join(["%s"] * len(product_ids))
        cur.execute(f"DELETE FROM product_favorites WHERE product_id IN ({ph})", product_ids)
        cur.execute(f"DELETE FROM product_reviews WHERE product_id IN ({ph})", product_ids)
        cur.execute(f"DELETE FROM product_skus WHERE product_id IN ({ph})", product_ids)
        cur.execute(f"DELETE FROM products WHERE id IN ({ph})", product_ids)

    # 手机号段用户及其残留收藏
    mobiles = [str(MOBILE_START + i) for i in range(N_USERS)]
    ph = ",".join(["%s"] * len(mobiles))
    cur.execute(f"SELECT id FROM users WHERE mobile IN ({ph})", mobiles)
    user_ids = [r["id"] for r in cur.fetchall()]
    if user_ids:
        phu = ",".join(["%s"] * len(user_ids))
        cur.execute(f"DELETE FROM product_favorites WHERE user_id IN ({phu})", user_ids)
        cur.execute(f"DELETE FROM product_reviews WHERE user_id IN ({phu})", user_ids)
        cur.execute(f"DELETE FROM users WHERE id IN ({phu})", user_ids)


def seed_users(cur) -> list[int]:
    pwd = password_hash()
    user_ids: list[int] = []
    for i in range(N_USERS):
        mobile = str(MOBILE_START + i)
        nick = f"ICF用户{i:03d}"
        cur.execute(
            """
            INSERT INTO users (mobile, password, nickname, status, role)
            VALUES (%s, %s, %s, 1, 'user')
            """,
            (mobile, pwd, nick),
        )
        user_ids.append(int(cur.lastrowid))
    print(f"users: {len(user_ids)}")
    return user_ids


def seed_products(cur, shop_id: int) -> tuple[list[int], list[int], dict[int, int]]:
    """Returns product_ids, sku_ids (aligned), cluster_of[product_index]."""
    cur.execute(
        "SELECT id FROM product_categories ORDER BY id ASC LIMIT 1"
    )
    cat = cur.fetchone()
    category_id = int(cat["id"]) if cat else 1

    product_ids: list[int] = []
    sku_ids: list[int] = []
    cluster_of: dict[int, int] = {}  # product index 0..119 -> cluster

    for idx in range(N_PRODUCTS):
        cluster = idx // PER_CLUSTER
        local = idx % PER_CLUSTER
        cluster_of[idx] = cluster
        name = f"{CLUSTER_NAMES[cluster]}-{local + 1:02d}"
        product_no = f"{PRODUCT_PREFIX}{idx:03d}"
        sku_no = f"{SKU_PREFIX}{idx:03d}"
        price = round(19.9 + (cluster * 7) + local * 1.5, 2)
        pet = PET_TYPES[cluster]

        cur.execute(
            """
            INSERT INTO products (
              shop_id, product_no, name, subtitle, sale_price, stock,
              category_id, status, pet_type, is_recommend, publish_time
            ) VALUES (
              %s, %s, %s, %s, %s, %s, %s, 'on_sale', %s, 1, NOW()
            )
            """,
            (
                shop_id,
                product_no,
                name,
                f"ItemCF簇{cluster}",
                price,
                9999,
                category_id,
                pet,
            ),
        )
        pid = int(cur.lastrowid)
        product_ids.append(pid)

        cur.execute(
            """
            INSERT INTO product_skus (
              product_id, shop_id, sku_no, spec_values, spec_key,
              sale_price, stock, status
            ) VALUES (
              %s, %s, %s, %s, 'default', %s, 9999, 'enabled'
            )
            """,
            (pid, shop_id, sku_no, '{"规格":"默认"}', price),
        )
        sku_ids.append(int(cur.lastrowid))

    print(f"products/skus: {len(product_ids)}")
    return product_ids, sku_ids, cluster_of


def pick_items_for_order(
    rng: random.Random,
    user_idx: int,
    cluster_of: dict[int, int],
) -> list[int]:
    """Return list of product indices (2～4) with same-cluster bias."""
    main = user_idx % N_CLUSTERS
    n = rng.randint(2, 4)
    same = rng.random() < SAME_CLUSTER_RATIO
    if same:
        pool = [i for i, c in cluster_of.items() if c == main]
    else:
        other = (main + rng.randint(1, N_CLUSTERS - 1)) % N_CLUSTERS
        pool = [i for i, c in cluster_of.items() if c == other]
    if len(pool) < n:
        n = len(pool)
    return rng.sample(pool, n)


def seed_orders(
    cur,
    shop_id: int,
    user_ids: list[int],
    product_ids: list[int],
    sku_ids: list[int],
    cluster_of: dict[int, int],
    rng: random.Random,
) -> list[tuple[int, str, int, list[int]]]:
    """Returns list of (order_id, order_no, user_id, product_ids_in_order)."""
    created: list[tuple[int, str, int, list[int]]] = []
    base_time = datetime.now() - timedelta(days=90)

    for oi in range(N_ORDERS):
        user_idx = oi % N_USERS
        user_id = user_ids[user_idx]
        pidxs = pick_items_for_order(rng, user_idx, cluster_of)

        lines = []
        goods = 0.0
        for pidx in pidxs:
            qty = rng.randint(1, 2)
            # fetch price from products already inserted — use deterministic price
            price = round(19.9 + (cluster_of[pidx] * 7) + (pidx % PER_CLUSTER) * 1.5, 2)
            lines.append((pidx, qty, price))
            goods += price * qty
        goods = round(goods, 2)

        order_no = f"{ORDER_PREFIX}{oi:05d}"
        ts = base_time + timedelta(minutes=oi * 17)
        cur.execute(
            """
            INSERT INTO orders (
              order_no, user_id, shop_id, total_amount, goods_amount,
              discount_amount, pay_amount, receiver_name, receiver_phone,
              receiver_address, status, completed_at, created_at, updated_at
            ) VALUES (
              %s, %s, %s, %s, %s, 0, %s, %s, %s, %s, 'completed', %s, %s, %s
            )
            """,
            (
                order_no,
                user_id,
                shop_id,
                goods,
                goods,
                goods,
                f"收货人{user_idx}",
                str(MOBILE_START + user_idx),
                "测试省测试市测试区 ItemCF地址",
                ts,
                ts,
                ts,
            ),
        )
        order_id = int(cur.lastrowid)
        pids_in_order = []
        for pidx, qty, price in lines:
            pid = product_ids[pidx]
            sid = sku_ids[pidx]
            pids_in_order.append(pid)
            cur.execute(
                """
                INSERT INTO order_items (
                  order_id, product_id, sku_id, product_name, price, quantity, created_at
                ) VALUES (%s, %s, %s, %s, %s, %s, %s)
                """,
                (
                    order_id,
                    pid,
                    sid,
                    f"{CLUSTER_NAMES[cluster_of[pidx]]}-{(pidx % PER_CLUSTER) + 1:02d}",
                    price,
                    qty,
                    ts,
                ),
            )
        created.append((order_id, order_no, user_id, pids_in_order))

    print(f"orders: {len(created)}")
    return created


def seed_favorites(
    cur,
    user_ids: list[int],
    product_ids: list[int],
    cluster_of: dict[int, int],
    rng: random.Random,
) -> int:
    seen: set[tuple[int, int]] = set()
    n = 0
    attempts = 0
    while n < N_FAVORITES and attempts < N_FAVORITES * 5:
        attempts += 1
        ui = rng.randrange(N_USERS)
        main = ui % N_CLUSTERS
        if rng.random() < 0.85:
            pool = [i for i, c in cluster_of.items() if c == main]
        else:
            pool = list(range(N_PRODUCTS))
        pidx = rng.choice(pool)
        key = (user_ids[ui], product_ids[pidx])
        if key in seen:
            continue
        seen.add(key)
        try:
            cur.execute(
                "INSERT INTO product_favorites (user_id, product_id) VALUES (%s, %s)",
                key,
            )
            n += 1
        except pymysql.err.IntegrityError:
            continue
    print(f"favorites: {n}")
    return n


def seed_reviews(
    cur,
    orders: list[tuple[int, str, int, list[int]]],
    shop_id: int,
    rng: random.Random,
) -> int:
    sample = rng.sample(orders, min(N_REVIEWS, len(orders)))
    n = 0
    for order_id, order_no, user_id, pids in sample:
        if not pids:
            continue
        pid = rng.choice(pids)
        rating = rng.choice([4, 5, 5, 5])
        cur.execute(
            """
            INSERT INTO product_reviews (
              order_id, order_no, user_id, shop_id, product_id, sku_id,
              rating, content, is_anonymous, status
            ) VALUES (%s, %s, %s, %s, %s, 0, %s, %s, 0, 'visible')
            """,
            (
                order_id,
                order_no,
                user_id,
                shop_id,
                pid,
                rating,
                "ItemCF种子评价，商品不错",
            ),
        )
        n += 1
    print(f"reviews: {n}")
    return n


def write_csv(
    path: Path,
    orders: list[tuple[int, str, int, list[int]]],
) -> None:
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w", newline="", encoding="utf-8") as f:
        w = csv.writer(f)
        w.writerow(["user_id", "product_id", "weight", "event"])
        for _, _, user_id, pids in orders:
            for pid in set(pids):
                w.writerow([user_id, pid, 5, "purchase"])
    print(f"csv: {path}")


def print_cooccurrence_sample(
    orders: list[tuple[int, str, int, list[int]]],
    product_ids: list[int],
    cluster_of: dict[int, int],
) -> None:
    # Build pair counts from baskets
    pair: Counter[tuple[int, int]] = Counter()
    for _, _, _, pids in orders:
        uniq = sorted(set(pids))
        for i in range(len(uniq)):
            for j in range(i + 1, len(uniq)):
                a, b = uniq[i], uniq[j]
                pair[(a, b)] += 1

    # Pick first product of cluster 0
    seed_pid = product_ids[0]
    neigh: list[tuple[int, int]] = []
    for (a, b), c in pair.items():
        if a == seed_pid:
            neigh.append((b, c))
        elif b == seed_pid:
            neigh.append((a, c))
    neigh.sort(key=lambda x: -x[1])
    print(f"\nco-occurrence sample for product_id={seed_pid} (簇0 第1件) Top5:")
    id_to_idx = {pid: i for i, pid in enumerate(product_ids)}
    for pid, cnt in neigh[:5]:
        idx = id_to_idx.get(pid, -1)
        cl = cluster_of.get(idx, -1)
        print(f"  neighbor product_id={pid} cluster={cl} co_count={cnt}")


def main() -> None:
    parser = argparse.ArgumentParser(description="Seed ItemCF test data")
    parser.add_argument(
        "--clean",
        action="store_true",
        help="删除此前 ICF 前缀数据后再插入（推荐）",
    )
    parser.add_argument("--seed", type=int, default=42, help="随机种子")
    args = parser.parse_args()
    rng = random.Random(args.seed)

    conn = connect()
    try:
        with conn.cursor() as cur:
            shop_id = resolve_shop_id(cur)
            print(f"shop_id={shop_id}")
            if args.clean:
                clean(cur)
                conn.commit()

            user_ids = seed_users(cur)
            product_ids, sku_ids, cluster_of = seed_products(cur, shop_id)
            orders = seed_orders(
                cur, shop_id, user_ids, product_ids, sku_ids, cluster_of, rng
            )
            fav_n = seed_favorites(cur, user_ids, product_ids, cluster_of, rng)
            rev_n = seed_reviews(cur, orders, shop_id, rng)
            conn.commit()

            # optional CSV for offline job
            root = Path(__file__).resolve().parent
            write_csv(root / "data" / "itemcf_interactions.csv", orders)

            print("\n=== summary ===")
            print(f"users={len(user_ids)} products={len(product_ids)} orders={len(orders)}")
            print(f"favorites≈{fav_n} reviews={rev_n}")
            print_cooccurrence_sample(orders, product_ids, cluster_of)
            print("\ndone. 密码均为 123456；登录示例手机 13800001000")
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()
