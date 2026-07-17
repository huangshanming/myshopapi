-- 商家端分层菜单种子（可单独执行，幂等）
-- mysql -u homestead -p mymall < scripts/seed-shop-menus.sql
USE mymall;

INSERT INTO shop_menus (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(100, 0, '商品中心', 'dir', '', '', 'Goods', '', 10, 1, 1),
(101, 0, '库存管理', 'dir', '', '', 'Box', '', 20, 1, 1),
(102, 0, '订单中心', 'dir', '', '', 'List', '', 30, 1, 1),
(103, 0, '店铺设置', 'dir', '', '', 'Setting', '', 90, 1, 1),
(1, 100, '商品列表', 'menu', '/merchant/products', 'merchant/Products', 'Goods', 'product:list', 10, 1, 1),
(2, 100, '发布商品', 'menu', '/merchant/products/edit', 'merchant/ProductEdit', 'Edit', 'product:edit', 11, 1, 1),
(3, 100, '回收站', 'menu', '/merchant/products/recycle', 'merchant/ProductRecycle', 'Delete', 'product:recycle', 12, 1, 1),
(7, 100, '操作日志', 'menu', '/merchant/products/op-logs', 'merchant/OpLogs', 'Document', 'product:list', 13, 1, 1),
(4, 101, '库存预警', 'menu', '/merchant/stocks/warnings', 'merchant/StockWarnings', 'Warning', 'stock:warn', 10, 1, 1),
(5, 102, '店铺订单', 'menu', '/merchant/orders', 'merchant/Orders', 'List', 'order:list', 10, 1, 1),
(6, 103, '员工权限', 'menu', '/merchant/staff', 'merchant/Staff', 'User', 'shop:staff', 10, 1, 1),
(104, 0, '内容中心', 'dir', '', '', 'Document', '', 40, 1, 1),
(8, 104, '我的文章', 'menu', '/merchant/articles', 'merchant/Articles', 'Document', 'article:list', 10, 1, 1),
(9, 104, '发布文章', 'menu', '/merchant/articles/edit', 'merchant/ArticleEdit', 'Edit', 'article:edit', 11, 1, 1),
(11, 1, '商品新增', 'button', '', '', '', 'product:add', 1, 1, 1),
(12, 1, '商品编辑', 'button', '', '', '', 'product:edit', 2, 1, 1),
(13, 1, '商品上下架', 'button', '', '', '', 'product:status', 3, 1, 1),
(14, 1, '批量操作', 'button', '', '', '', 'product:batch', 4, 1, 1),
(15, 1, '导入导出', 'button', '', '', '', 'product:import', 5, 1, 1),
(16, 8, '文章新增', 'button', '', '', '', 'article:add', 1, 1, 1),
(17, 8, '文章编辑', 'button', '', '', '', 'article:edit', 2, 1, 1),
(18, 8, '文章删除', 'button', '', '', '', 'article:delete', 3, 1, 1)
ON DUPLICATE KEY UPDATE
  parent_id=VALUES(parent_id), name=VALUES(name), type=VALUES(type),
  path=VALUES(path), component=VALUES(component), icon=VALUES(icon),
  perms=VALUES(perms), sort=VALUES(sort), visible=VALUES(visible), status=VALUES(status);

-- 店主角色补挂全部菜单
INSERT IGNORE INTO shop_role_menus (role_id, menu_id)
SELECT r.id, m.id FROM shop_roles r
CROSS JOIN shop_menus m
WHERE r.code = 'shop_owner' AND m.status = 1;
