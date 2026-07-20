-- 积分商城示例商品
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-points-mall-products.sql
USE mymall;

INSERT INTO points_products
  (name, cover_url, description, points_price, stock, per_user_limit, status, sort)
SELECT * FROM (
  SELECT '精美马克杯' AS name,
         'https://picsum.photos/id/30/400/400' AS cover_url,
         '陶瓷马克杯，容量约 350ml，适合日常喝水喝咖啡。' AS description,
         200 AS points_price, 100 AS stock, 2 AS per_user_limit, 'on' AS status, 100 AS sort
  UNION ALL SELECT '帆布托特包', 'https://picsum.photos/id/1011/400/400',
         '简约帆布袋，可单肩或手提，日常通勤出街都好搭。', 350, 80, 1, 'on', 95
  UNION ALL SELECT '蓝牙耳机（基础款）', 'https://picsum.photos/id/160/400/400',
         '轻便入耳式蓝牙耳机，续航约 4 小时，适合通勤听歌。', 1200, 50, 1, 'on', 90
  UNION ALL SELECT '充电宝 10000mAh', 'https://picsum.photos/id/2/400/400',
         '小巧便携充电宝，双口输出，出行应急充电。', 800, 60, 1, 'on', 85
  UNION ALL SELECT '不锈钢保温杯', 'https://picsum.photos/id/225/400/400',
         '真空保温杯，保温约 6 小时，办公室/户外都适用。', 450, 90, 2, 'on', 80
  UNION ALL SELECT '品牌帆布鞋（优惠券）', 'https://picsum.photos/id/21/400/400',
         '兑换后发放满减券码（虚拟商品，客服核销）。', 500, 200, 3, 'on', 75
  UNION ALL SELECT '护理手霜套装', 'https://picsum.photos/id/292/400/400',
         '保湿护手霜两支装，清爽不油腻。', 180, 120, 2, 'on', 70
  UNION ALL SELECT '桌面小夜灯', 'https://picsum.photos/id/201/400/400',
         '触控调光小夜灯，三档亮度，宿舍/床头适用。', 280, 70, 1, 'on', 65
  UNION ALL SELECT '运动毛巾', 'https://picsum.photos/id/96/400/400',
         '速干运动毛巾，健身跑步擦汗好物。', 120, 150, 3, 'on', 60
  UNION ALL SELECT '零食大礼包', 'https://picsum.photos/id/1080/400/400',
         '休闲零食组合装，约 8～10 件，口味随机。', 300, 100, 2, 'on', 55
  UNION ALL SELECT '猫狗零食小礼盒', 'https://picsum.photos/id/1025/400/400',
         '宠物零食小礼盒（犬猫通用口味随机）。', 220, 80, 2, 'on', 50
  UNION ALL SELECT '定制帆布徽章', 'https://picsum.photos/id/119/400/400',
         '平台限定金属徽章，收藏纪念款。', 80, 300, 5, 'on', 45
) AS t
WHERE NOT EXISTS (
  SELECT 1 FROM points_products p WHERE p.name = t.name
);
