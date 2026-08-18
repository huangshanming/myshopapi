-- 示例抽奖活动（需先执行 ddl/lottery-tables.sql + lottery-prize-stock-strict.sql）
-- 封面图需存在于 uploads/lottery/（见仓库 uploads/lottery/*.png）
USE mymall;

INSERT INTO lottery_activities (id, title, status, cost_points, daily_limit, start_at, end_at)
VALUES (1, '新人九宫格抽奖', 1, 10, 3, NOW() - INTERVAL 1 DAY, NOW() + INTERVAL 365 DAY)
ON DUPLICATE KEY UPDATE title=VALUES(title), status=1, cost_points=VALUES(cost_points), daily_limit=VALUES(daily_limit);

INSERT INTO lottery_prizes (activity_id, slot, name, cover_url, prize_type, points_amount, weight, stock, stock_strict) VALUES
(1, 0, '10积分', '/uploads/lottery/points-coin.png', 'points', 10, 25, -1, 0),
(1, 1, '谢谢参与', '', 'thanks', 0, 35, -1, 0),
(1, 2, '50积分', '/uploads/lottery/points-redpacket.png', 'points', 50, 12, -1, 0),
(1, 3, '蓝牙耳机', '/uploads/lottery/physical-headphones.png', 'physical', 0, 6, 20, 1),
(1, 4, '100积分', '/uploads/lottery/points-diamond.png', 'points', 100, 8, -1, 0),
(1, 5, '精美礼盒', '/uploads/lottery/physical-gift.png', 'physical', 0, 8, 30, 1),
(1, 6, '20积分', '/uploads/lottery/points-coin.png', 'points', 20, 18, -1, 0),
(1, 7, '谢谢参与', '', 'thanks', 0, 35, -1, 0),
(1, 8, '智能手机', '/uploads/lottery/physical-phone.png', 'physical', 0, 2, 5, 1)
ON DUPLICATE KEY UPDATE
  name=VALUES(name), cover_url=VALUES(cover_url), prize_type=VALUES(prize_type),
  points_amount=VALUES(points_amount), weight=VALUES(weight), stock=VALUES(stock),
  stock_strict=VALUES(stock_strict);
