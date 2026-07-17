-- 已有库增量：文章管理菜单（幂等）
USE mymall;

INSERT INTO sys_menu (id, parent_id, name, type, path, component, icon, perms, sort, visible, status) VALUES
(90, 0, '文章管理', 'dir', '', '', 'Document', '', 20, 1, 1),
(91, 90, '文章列表', 'menu', '/admin/articles', 'admin/articles/Articles', 'Document', 'community:article:list', 21, 1, 1),
(92, 90, '分类管理', 'menu', '/admin/articles/categories', 'admin/articles/Categories', 'Menu', 'community:article:category', 22, 1, 1),
(93, 90, '评论管理', 'menu', '/admin/articles/comments', 'admin/articles/Comments', 'ChatDotRound', 'community:article:comment', 23, 1, 1),
(94, 90, '文章回收站', 'menu', '/admin/articles/recycle', 'admin/articles/Recycle', 'Delete', 'community:article:recycle', 24, 1, 1),
(95, 90, '文章统计', 'menu', '/admin/articles/stats', 'admin/articles/Stats', 'DataAnalysis', 'community:article:stats', 25, 1, 1),
(96, 91, '文章新增', 'button', '', '', '', 'community:article:add', 1, 1, 1),
(97, 91, '文章编辑', 'button', '', '', '', 'community:article:edit', 2, 1, 1),
(98, 91, '文章审核', 'button', '', '', '', 'community:article:audit', 3, 1, 1),
(99, 91, '文章置顶', 'button', '', '', '', 'community:article:top', 4, 1, 1)
ON DUPLICATE KEY UPDATE name=VALUES(name), path=VALUES(path), component=VALUES(component),
  icon=VALUES(icon), perms=VALUES(perms), sort=VALUES(sort), visible=1, status=1;

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT 1, id FROM sys_menu WHERE id BETWEEN 90 AND 99;
