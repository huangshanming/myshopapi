#!/usr/bin/env python3
"""将 images/ 素材填充为店铺 #1 的完整上架商品（幂等）。

用法（仓库根目录）:
  python3 scripts/seed-shop-demo-products.py

会：
  1. 解压 images/红酒.zip
  2. 复制图片到 uploads/products/1/（仓库根，catalog 会回退读取）
  3. 写入分类 + products + product_skus + product_images
"""

from __future__ import annotations

import hashlib
import json
import os
import shutil
import zipfile
from datetime import datetime
from pathlib import Path

import pymysql
import yaml

ROOT = Path(__file__).resolve().parents[1]
IMAGES = ROOT / "images"
UPLOAD_DIR = ROOT / "uploads" / "products" / "1"
SHOP_ID = 1
CFG = ROOT / "services" / "catalog-service" / "etc" / "catalog-service.yaml"


def md5_spec_key(values: dict) -> str:
    if not values:
        return "default"
    parts = [f"{k}={values[k]}" for k in sorted(values.keys())]
    return hashlib.md5("|".join(parts).encode()).hexdigest()


def sku_no(product_no: str, spec_key: str) -> str:
    short = spec_key[:8]
    return f"{product_no}-{short}"


def load_db():
    cfg = yaml.safe_load(CFG.read_text())
    m = cfg["mysql"]
    return pymysql.connect(
        host=m.get("host", "127.0.0.1"),
        port=int(m.get("port", 3306)),
        user=m["username"],
        password=m["password"],
        database=m.get("database", "mymall"),
        charset="utf8mb4",
        autocommit=False,
    )


def ensure_categories(cur) -> dict[str, int]:
    """返回 名称 -> category_id"""
    wanted = [
        ("美妆护肤", "口红/香水/洁面等"),
        ("服饰鞋包", "毛衣等服饰"),
        ("休闲食品", "饼干/橘红等"),
        ("生鲜烘焙", "蛋糕等冷藏配送"),
        ("酒水饮料", "葡萄酒等"),
        ("玩具文创", "玩偶等"),
    ]
    out: dict[str, int] = {}
    for name, desc in wanted:
        cur.execute("SELECT id FROM product_categories WHERE name=%s LIMIT 1", (name,))
        row = cur.fetchone()
        if row:
            out[name] = int(row[0])
            continue
        cur.execute(
            """
            INSERT INTO product_categories
              (parent_id, name, description, sort_order, level, is_show, product_count)
            VALUES (0, %s, %s, 10, 1, 1, 0)
            """,
            (name, desc),
        )
        out[name] = int(cur.lastrowid)
        print(f"  + category {name} -> {out[name]}")
    return out


def prepare_assets() -> dict[str, str]:
    """复制素材到 uploads，返回 slug -> 相对 URL"""
    UPLOAD_DIR.mkdir(parents=True, exist_ok=True)
    mapping: dict[str, str] = {}

    # 红酒 zip
    wine_zip = IMAGES / "红酒.zip"
    wine_jpg = UPLOAD_DIR / "wine-red.jpg"
    if wine_zip.exists():
        tmp = ROOT / "tmp" / "wine-extract"
        if tmp.exists():
            shutil.rmtree(tmp)
        tmp.mkdir(parents=True, exist_ok=True)
        with zipfile.ZipFile(wine_zip, "r") as zf:
            zf.extractall(tmp)
        jpg = next(tmp.rglob("*.jpg"), None) or next(tmp.rglob("*.JPG"), None)
        if jpg:
            shutil.copy2(jpg, wine_jpg)
            mapping["wine"] = f"/uploads/products/1/{wine_jpg.name}"
            print(f"  + wine image -> {wine_jpg.name}")
        shutil.rmtree(tmp, ignore_errors=True)

    files = {
        "perfume": "香水.png",
        "lipstick1": "口红1.png",
        "lipstick2": "口红2.png",
        "cleanser": "洗面奶.png",
        "sweater": "毛衣.png",
        "cake": "奶油蛋糕.png",
        "cookie": "曲奇饼干.png",
        "tangerine": "化州橘红.png",
        "doll": "布偶娃娃.png",
    }
    for slug, fname in files.items():
        src = IMAGES / fname
        if not src.exists():
            print(f"  ! missing {fname}")
            continue
        ext = src.suffix.lower()
        dest = UPLOAD_DIR / f"{slug}{ext}"
        shutil.copy2(src, dest)
        mapping[slug] = f"/uploads/products/1/{dest.name}"
        print(f"  + {fname} -> {dest.name}")
    return mapping


def upsert_product(cur, p: dict, assets: dict[str, str], cats: dict[str, int]):
    product_no = p["product_no"]
    cur.execute("SELECT id FROM products WHERE product_no=%s LIMIT 1", (product_no,))
    row = cur.fetchone()
    now = datetime.now().strftime("%Y-%m-%d %H:%M:%S")

    main_key = p["main_asset"]
    main_url = assets.get(main_key, "")
    gallery_keys = p.get("gallery_assets") or [main_key]
    gallery_urls = [assets[k] for k in gallery_keys if k in assets]
    if main_url and main_url not in gallery_urls:
        gallery_urls.insert(0, main_url)

    cat_id = cats[p["category"]]
    spec_json = json.dumps(p.get("spec_json") or [], ensure_ascii=False)
    skus = p.get("skus") or [{"spec": {}, "sale_price": p["sale_price"], "stock": p["stock"]}]
    total_stock = sum(int(s["stock"]) for s in skus)
    min_price = min(float(s["sale_price"]) for s in skus)

    fields = dict(
        shop_id=SHOP_ID,
        name=p["name"],
        subtitle=p["subtitle"],
        description=p["description"],
        main_image=main_url,
        market_price=p.get("market_price", min_price * 1.2),
        sale_price=min_price,
        stock=total_stock,
        stock_warn=10,
        category_id=cat_id,
        product_type=p.get("product_type", "physical"),
        spec_json=spec_json,
        status="on_sale",
        pet_type="other",
        shelf_life=p.get("shelf_life"),
        storage_condition=p.get("storage_condition"),
        unit=p.get("unit") or "件",
        is_new=1,
        publish_time=now,
        deleted_at=None,
    )

    if row:
        pid = int(row[0])
        sets = ", ".join(f"{k}=%s" for k in fields)
        cur.execute(
            f"UPDATE products SET {sets}, updated_at=%s WHERE id=%s",
            (*fields.values(), now, pid),
        )
        print(f"  ~ update {product_no} id={pid}")
    else:
        cols = ["product_no"] + list(fields.keys())
        placeholders = ", ".join(["%s"] * len(cols))
        cur.execute(
            f"INSERT INTO products ({', '.join(cols)}) VALUES ({placeholders})",
            (product_no, *fields.values()),
        )
        pid = int(cur.lastrowid)
        print(f"  + create {product_no} id={pid}")

    # SKUs：演示数据直接物理删除后重建，避免 sku_no 唯一键冲突
    cur.execute("DELETE FROM product_skus WHERE product_id=%s", (pid,))
    for s in skus:
        spec = s.get("spec") or {}
        sk = md5_spec_key(spec)
        sn = sku_no(product_no, sk)
        sv = json.dumps(spec, ensure_ascii=False)
        cur.execute(
            """
            INSERT INTO product_skus
              (product_id, shop_id, sku_no, spec_values, spec_key, sale_price, market_price,
               stock, stock_warn, status, sold_count, deleted_at)
            VALUES (%s,%s,%s,%s,%s,%s,%s,%s,10,'enabled',0,NULL)
            """,
            (
                pid,
                SHOP_ID,
                sn,
                sv,
                sk,
                float(s["sale_price"]),
                float(s.get("market_price") or float(s["sale_price"]) * 1.15),
                int(s["stock"]),
            ),
        )

    cur.execute("DELETE FROM product_images WHERE product_id=%s", (pid,))
    for i, url in enumerate(gallery_urls):
        typ = "main" if i == 0 else "gallery"
        cur.execute(
            """
            INSERT INTO product_images (product_id, shop_id, url, typ, sort)
            VALUES (%s,%s,%s,%s,%s)
            """,
            (pid, SHOP_ID, url, typ, i),
        )
    return pid


PRODUCTS = [
    {
        "product_no": "DEMO-PERFUME-001",
        "category": "美妆护肤",
        "main_asset": "perfume",
        "name": "晨雾花语淡香水 50ml",
        "subtitle": "清新花果调，日常通勤与约会都适合",
        "sale_price": 199,
        "market_price": 259,
        "stock": 80,
        "unit": "瓶",
        "description": """【产品亮点】
晨雾花语淡香水以清新花果调为主，前调清甜、中调柔和、尾调干净留香，适合日常通勤、周末约会等场景。

【香调与规格】
- 香调：花果调（前调柑橘/浆果，中调花香，尾调麝香）
- 净含量：50ml
- 质地：喷雾型淡香水（EDT）

【使用建议】
喷于手腕内侧、耳后、锁骨等脉搏部位，距离皮肤约 15cm。先少量试喷，确认无不适再常规使用。

【注意事项】
避免接触眼睛与破损皮肤；孕妇及敏感肌请先局部试敏。请置于阴凉避光处，远离明火与高温。

【发货说明】
48 小时内发货，密封包装。支持七天无理由（未拆封）。""",
    },
    {
        "product_no": "DEMO-LIP-ROSE-001",
        "category": "美妆护肤",
        "main_asset": "lipstick1",
        "gallery_assets": ["lipstick1", "lipstick2"],
        "name": "丝绒哑光口红 · 豆沙玫瑰",
        "subtitle": "轻薄雾面，一抹显气色，日常百搭色",
        "sale_price": 89,
        "market_price": 129,
        "stock": 120,
        "unit": "支",
        "description": """【产品亮点】
丝绒哑光质地，上嘴轻盈不拔干，豆沙玫瑰色显白显气色，办公室与日常妆都好搭。

【色号与妆效】
- 色号：豆沙玫瑰（偏暖的中性玫瑰色）
- 妆效：雾面哑光，轻薄贴唇

【使用方法】
直接涂抹或用唇刷勾勒唇形。可先用润唇膏打底，提升舒适度与持妆。

【成分配方说明】
含滋润油脂与成膜成分，降低干燥感。具体成分表见包装标识。

【注意事项】
如有不适请立即停用。请置于阴凉干燥处，开封后建议 12 个月内用完。

【发货说明】
全新密封包装发货，支持未拆封七天无理由。""",
    },
    {
        "product_no": "DEMO-LIP-CORAL-002",
        "category": "美妆护肤",
        "main_asset": "lipstick2",
        "name": "丝绒哑光口红 · 蜜桃珊瑚",
        "subtitle": "活力珊瑚色，提亮肤色，春夏主打",
        "sale_price": 89,
        "market_price": 129,
        "stock": 100,
        "unit": "支",
        "description": """【产品亮点】
蜜桃珊瑚色系，显白又有活力，适合春夏妆容与轻约会场景。丝绒雾面，不易飞粉。

【色号与妆效】
- 色号：蜜桃珊瑚
- 妆效：哑光雾面，显色均匀

【使用建议】
可叠涂提升饱和度；内唇轻涂可打造咬唇妆。

【注意事项】
敏感唇请先试色。避免高温暴晒存放。

【发货说明】
独立包装，48 小时内发出。""",
    },
    {
        "product_no": "DEMO-CLEANSER-001",
        "category": "美妆护肤",
        "main_asset": "cleanser",
        "name": "氨基酸温和洁面乳 120g",
        "subtitle": "低刺激洁面，洗后不紧绷，敏感肌友好",
        "sale_price": 69,
        "market_price": 99,
        "stock": 150,
        "unit": "支",
        "description": """【产品亮点】
氨基酸表活体系，清洁力温和，洗后皮肤柔软不紧绷，适合每日早晚洁面。

【适用肤质】
干性、中性、混合性及轻度敏感肌（建议先局部试用）。

【使用方法】
取适量于掌心加水揉出泡沫，轻柔按摩面部 30～60 秒，用清水冲洗干净。

【储存与保质】
常温避光保存。开封后建议 6 个月内用完。

【发货说明】
全新正品密封，支持未拆封退换。""",
    },
    {
        "product_no": "DEMO-SWEATER-001",
        "category": "服饰鞋包",
        "main_asset": "sweater",
        "name": "软糯圆领针织毛衣",
        "subtitle": "亲肤保暖，宽松版型，秋冬叠穿必备",
        "sale_price": 159,
        "market_price": 229,
        "stock": 0,  # 由 SKU 汇总
        "unit": "件",
        "product_type": "physical",
        "spec_json": [{"name": "尺码", "values": ["S", "M", "L"]}],
        "skus": [
            {"spec": {"尺码": "S"}, "sale_price": 159, "stock": 30},
            {"spec": {"尺码": "M"}, "sale_price": 159, "stock": 45},
            {"spec": {"尺码": "L"}, "sale_price": 169, "stock": 25},
        ],
        "description": """【产品亮点】
柔软针织面料，圆领基础款，版型宽松不挑身材，可单穿或内搭衬衫/打底。

【面料与工艺】
- 主要面料：针织纱线（亲肤柔软）
- 工艺：圆领罗纹收边，不易变形

【尺码参考】
- S：适合 155–162cm / 45–52kg
- M：适合 160–168cm / 52–60kg
- L：适合 165–175cm / 60–70kg
（仅供参考，以自身习惯版型为准）

【洗涤说明】
建议反面轻柔手洗或冷水机洗网袋，平铺晾干，避免暴晒与高温烘干。

【发货说明】
48 小时内发货，支持七天无理由（不影响二次销售）。""",
    },
    {
        "product_no": "DEMO-CAKE-001",
        "category": "生鲜烘焙",
        "main_asset": "cake",
        "name": "动物奶油鲜奶蛋糕 6 寸",
        "subtitle": "当日烘焙，口感轻盈，生日聚会优选",
        "sale_price": 68,
        "market_price": 98,
        "stock": 40,
        "unit": "个",
        "product_type": "fresh",
        "shelf_life": 2,
        "storage_condition": "冷藏 0–4℃，开封后尽快食用",
        "description": """【产品亮点】
动物奶油裱花，蛋糕体松软，甜度适中，适合生日、下午茶与小聚。

【规格】
- 尺寸：约 6 寸（适合 2–4 人分享）
- 口味：原味鲜奶（以店铺当日出品为准）

【食用与储存】
收到后请冷藏保存，建议 24～48 小时内食用完毕。含乳制品，过敏者慎选。

【配送说明】
生鲜商品，优先同城冷链/保温袋发货；偏远地区请下单前咨询客服。

【售后】
非质量问题不支持无理由退货；运输导致严重破损请当日 intra联系客服并保留凭证。""",
    },
    {
        "product_no": "DEMO-COOKIE-001",
        "category": "休闲食品",
        "main_asset": "cookie",
        "name": "黄油曲奇饼干礼盒 200g",
        "subtitle": "浓郁奶香，酥脆不腻，伴手礼优选",
        "sale_price": 39,
        "market_price": 59,
        "stock": 200,
        "unit": "盒",
        "shelf_life": 90,
        "storage_condition": "常温阴凉干燥，开封后密封保存",
        "description": """【产品亮点】
以黄油烘焙，奶香浓郁、口感酥脆，独立礼盒适合自用或送礼。

【配料与过敏原】
含小麦粉、黄油、糖等（详见包装）。含麸质与乳制品，相关过敏者请勿食用。

【食用建议】
开袋即食。搭配咖啡、红茶风味更佳。

【储存】
常温避光防潮，开封后请密封并尽快食用。

【发货说明】
全新包装，保质期内发货。""",
    },
    {
        "product_no": "DEMO-JUHong-001",
        "category": "休闲食品",
        "main_asset": "tangerine",
        "name": "化州橘红切片 精选装 100g",
        "subtitle": "道地风味，可泡茶或直接咀嚼，甘香回味",
        "sale_price": 48,
        "market_price": 68,
        "stock": 160,
        "unit": "袋",
        "shelf_life": 365,
        "storage_condition": "密封干燥处保存，防潮防异味",
        "description": """【产品亮点】
精选化州橘红切片，果香清雅，可泡茶或直接细嚼，日常润喉茶饮好物。

【食用方法】
- 泡茶：取适量切片，热水冲泡 3–5 分钟，可反复冲泡
- 咀嚼：少量细嚼，慢慢品味甘香

【注意事项】
本产品为食品原料/休闲食品，不能替代药物。特殊体质请咨询专业人士。孕妇及儿童请在成人指导下食用。

【储存】
开封后请密封，置于干燥阴凉处。

【发货说明】
独立袋装，防潮包装发货。""",
    },
    {
        "product_no": "DEMO-DOLL-001",
        "category": "玩具文创",
        "main_asset": "doll",
        "name": "软萌布偶娃娃 毛绒公仔",
        "subtitle": "亲肤短毛，可爱造型，送礼与装饰两相宜",
        "sale_price": 79,
        "market_price": 119,
        "stock": 90,
        "unit": "个",
        "description": """【产品亮点】
柔软短毛面料，造型可爱，可作床头装饰、拍照道具或节日礼物。

【规格】
- 材质：短毛绒 + 填充棉
- 尺寸：以实物为准（约手持大小到中号，见主图比例）

【清洁与保养】
建议局部擦拭或放入洗衣袋轻柔机洗，阴干。避免高温暴晒以防褪色变形。

【安全提示】
请远离明火。三岁以下儿童请在成人看护下使用，注意填充物与细小配件。

【发货说明】
压缩包装可能略有压痕，取出抖动/晾放后可恢复蓬松。""",
    },
    {
        "product_no": "DEMO-WINE-001",
        "category": "酒水饮料",
        "main_asset": "wine",
        "name": "精选干红葡萄酒 750ml",
        "subtitle": "果香饱满，单宁柔顺，聚餐佐餐之选",
        "sale_price": 128,
        "market_price": 198,
        "stock": 60,
        "unit": "瓶",
        "shelf_life": 1825,
        "storage_condition": "阴凉避光横放，适宜 12–18℃",
        "description": """【产品亮点】
精选干红，果香层次清晰，单宁较为柔顺，适合牛排、硬质芝士与朋友小聚。

【品鉴建议】
开瓶后醒酒 15–30 分钟更佳。建议侍酒温度 16–18℃。

【规格】
- 净含量：750ml
- 类型：干红葡萄酒
- 酒精度：以瓶身标签为准

【购买须知】
本商品仅向年满 18 周岁的顾客销售。请勿酒后驾驶。孕妇及不宜饮酒人群请勿购买。

【储存与发货】
避光阴凉保存。酒类运输可能遇物流限制，偏远地区请先咨询客服。""",
    },
]


def main():
    print("== prepare assets ==")
    assets = prepare_assets()
    if "wine" not in assets:
        print("  ! wine image missing, DEMO-WINE-001 may have empty main_image")

    print("== write db ==")
    conn = load_db()
    try:
        cur = conn.cursor()
        cats = ensure_categories(cur)
        for p in PRODUCTS:
            if p["main_asset"] not in assets and p["main_asset"] != "wine":
                print(f"  ! skip {p['product_no']}: missing asset {p['main_asset']}")
                continue
            if p["main_asset"] == "wine" and "wine" not in assets:
                print(f"  ! skip {p['product_no']}: wine image missing")
                continue
            upsert_product(cur, p, assets, cats)
        conn.commit()
        print("== done ==")
        cur.execute(
            """
            SELECT product_no, name, status, sale_price, stock, main_image
            FROM products
            WHERE shop_id=%s AND product_no LIKE 'DEMO-%%'
            ORDER BY id
            """,
            (SHOP_ID,),
        )
        for r in cur.fetchall():
            print(" ", r)
    except Exception:
        conn.rollback()
        raise
    finally:
        conn.close()


if __name__ == "__main__":
    main()
