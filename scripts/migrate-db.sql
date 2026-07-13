-- mymall 本地开发：三微服务共用同一库 mymall
-- 在本机 MySQL 依次执行：
--   mysql -u homestead -p < scripts/migrate-db.sql
--   mysql -u homestead -p mymall < scripts/init-schema.sql
--   mysql -u homestead -p mymall < scripts/init-order-tables.sql

CREATE DATABASE IF NOT EXISTS mymall CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

GRANT ALL PRIVILEGES ON mymall.* TO 'homestead'@'%';
FLUSH PRIVILEGES;
