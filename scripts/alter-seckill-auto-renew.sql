-- 秒杀报名自动续费
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-seckill-auto-renew.sql
USE mymall;

ALTER TABLE seckill_entries
  ADD COLUMN IF NOT EXISTS auto_renew TINYINT NOT NULL DEFAULT 0 COMMENT '1=场次到期自动续报下一场' AFTER status;

-- 兼容不支持 IF NOT EXISTS 的 MySQL
-- ALTER TABLE seckill_entries ADD COLUMN auto_renew TINYINT NOT NULL DEFAULT 0 COMMENT '1=场次到期自动续报下一场' AFTER status;
