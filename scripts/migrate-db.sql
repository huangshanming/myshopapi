-- mymall 微服务数据库拆分脚本
-- 在本机 MySQL 执行：mysql -u homestead -p < scripts/migrate-db.sql

CREATE DATABASE IF NOT EXISTS user_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS catalog_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE DATABASE IF NOT EXISTS order_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- 若表仍在 mymall 库，迁移到独立库（已迁移可跳过）
-- RENAME TABLE mymall.users TO user_db.users;
-- RENAME TABLE mymall.products TO catalog_db.products;
-- RENAME TABLE mymall.product_categories TO catalog_db.product_categories;

-- Docker/K8s Pod 访问本机 MySQL 所需权限
GRANT ALL PRIVILEGES ON user_db.* TO 'homestead'@'%';
GRANT ALL PRIVILEGES ON catalog_db.* TO 'homestead'@'%';
GRANT ALL PRIVILEGES ON order_db.* TO 'homestead'@'%';
FLUSH PRIVILEGES;
