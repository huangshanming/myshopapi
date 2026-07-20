-- 总后台菜单拆分：由 4 大类细分为更清晰的目录
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/seed-admin-menu-restructure.sql
USE mymall;

-- 新目录（固定 id）
INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(140, 0, '商家管理', 'dir', '', '', 'Shop', '', 10, 1, 1),
(141, 0, '商品中心', 'dir', '', '', 'Goods', '', 15, 1, 1),
(142, 0, '交易中心', 'dir', '', '', 'List', '', 20, 1, 1),
(143, 0, '首页运营', 'dir', '', '', 'Picture', '', 35, 1, 1),
(144, 0, '用户触达', 'dir', '', '', 'Message', '', 50, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), icon=VALUES(icon), sort=VALUES(sort), visible=1, status=1, parent_id=0;

-- 原「业务管理」退役（子菜单已迁走）
UPDATE sys_menu SET name='业务管理(旧)', visible=0, sort=999 WHERE id=1;

-- 「文章管理」更名为「内容社区」
UPDATE sys_menu SET name='内容社区', sort=25 WHERE id=90;

-- 「营销中心」改为「营销玩法」，仅保留玩法类
UPDATE sys_menu SET name='营销玩法', sort=40 WHERE id=16;

-- 商家管理
UPDATE sys_menu SET parent_id=140, sort=10 WHERE id=2;  -- 入驻审核
UPDATE sys_menu SET parent_id=140, sort=20 WHERE id=3;  -- 店铺管理

-- 商品中心
UPDATE sys_menu SET parent_id=141, sort=10 WHERE id=5;  -- 全站商品
UPDATE sys_menu SET parent_id=141, sort=20 WHERE id=6;  -- 商品分类

-- 交易中心
UPDATE sys_menu SET parent_id=142, sort=10 WHERE id=4;   -- 全站订单
UPDATE sys_menu SET parent_id=142, sort=20 WHERE id=7;   -- 售后管理
UPDATE sys_menu SET parent_id=142, sort=30 WHERE id=8;   -- 物流管理
UPDATE sys_menu SET parent_id=142, sort=40 WHERE id=110; -- 评价管理

-- 首页运营
UPDATE sys_menu SET parent_id=143, sort=10 WHERE id=115; -- 首页 Banner
UPDATE sys_menu SET parent_id=143, sort=20 WHERE id=118; -- 主题集市
UPDATE sys_menu SET parent_id=143, sort=30 WHERE id=112; -- 首页展位

-- 营销中心（玩法）
UPDATE sys_menu SET parent_id=16, sort=10 WHERE id=121; -- 优惠券
UPDATE sys_menu SET parent_id=16, sort=20 WHERE id=17;  -- 秒杀规则
UPDATE sys_menu SET parent_id=16, sort=30 WHERE id=18;  -- 秒杀场次
UPDATE sys_menu SET parent_id=16, sort=40 WHERE id=125; -- 任务中心
UPDATE sys_menu SET parent_id=16, sort=50 WHERE id=127; -- 积分商城

-- 用户触达
UPDATE sys_menu SET parent_id=144, sort=10 WHERE id=124; -- 用户消息

-- 内容社区内排序微调
UPDATE sys_menu SET sort=10 WHERE id=91;
UPDATE sys_menu SET sort=20 WHERE id=92;
UPDATE sys_menu SET sort=30 WHERE id=93;
UPDATE sys_menu SET sort=40 WHERE id=100;
UPDATE sys_menu SET sort=50 WHERE id=94;
UPDATE sys_menu SET sort=60 WHERE id=95;

-- 超管绑定新目录
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id IN (140, 141, 142, 143, 144);
