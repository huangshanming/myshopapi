-- 当前活跃场次写入 5 条自动续费秒杀报名（每店上限 5），供首页展示。
-- 用法：mysql -u homestead -psecret mymall < scripts/seed-seckill-auto-renew.sql

SET @session_id := (
  SELECT id FROM seckill_sessions WHERE status = 'active' ORDER BY id DESC LIMIT 1
);

DELETE FROM seckill_entries WHERE session_id = @session_id AND shop_id = 1;

INSERT INTO seckill_entries
  (session_id, shop_id, product_id, product_name, product_image, origin_price, seckill_price, seckill_stock, fee_amount, status, auto_renew, created_at, updated_at)
VALUES
  (@session_id, 1, 3,  '晨雾花语淡香水 50ml',     '/uploads/products/1/perfume.png',   259.00, 159.00, 30, 500.00, 'active', 1, NOW(3), NOW(3)),
  (@session_id, 1, 4,  '丝绒哑光口红 · 豆沙玫瑰', '/uploads/products/1/lipstick1.png', 129.00,  69.00, 50, 500.00, 'active', 1, NOW(3), NOW(3)),
  (@session_id, 1, 7,  '软糯圆领针织毛衣',         '/uploads/products/1/sweater.png',   229.00, 129.00, 40, 500.00, 'active', 1, NOW(3), NOW(3)),
  (@session_id, 1, 8,  '动物奶油鲜奶蛋糕 6 寸',   '/uploads/products/1/cake.png',      98.00,  49.00, 20, 500.00, 'active', 1, NOW(3), NOW(3)),
  (@session_id, 1, 11, '软萌布偶娃娃 毛绒公仔',   '/uploads/products/1/doll.png',     119.00,  59.00, 35, 500.00, 'active', 1, NOW(3), NOW(3));

UPDATE shop_wallets
SET balance = balance - 2500.00, updated_at = NOW(3)
WHERE shop_id = 1 AND balance >= 2500.00;
