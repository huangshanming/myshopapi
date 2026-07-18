-- Canal 复制账号（本机 MySQL 执行一次即可）
-- 前置：已开启 ROW binlog（见 README「Canal / 库存预热」）
--
-- CREATE USER 在 MySQL 8 上可能需改认证插件，按本机版本调整。

CREATE USER IF NOT EXISTS 'canal'@'%' IDENTIFIED BY 'canal';
GRANT SELECT, REPLICATION SLAVE, REPLICATION CLIENT ON *.* TO 'canal'@'%';
FLUSH PRIVILEGES;
