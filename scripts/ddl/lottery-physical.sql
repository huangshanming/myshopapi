-- 抽奖实物履约字段（增量）
USE mymall;

ALTER TABLE lottery_prizes
  MODIFY COLUMN prize_type VARCHAR(32) NOT NULL DEFAULT 'thanks' COMMENT 'points|thanks|physical';

ALTER TABLE lottery_draw_records
  ADD COLUMN fulfill_status VARCHAR(16) NOT NULL DEFAULT 'none' COMMENT 'none|need_address|pending|shipped' AFTER status,
  ADD COLUMN address_id BIGINT UNSIGNED NOT NULL DEFAULT 0 AFTER fulfill_status,
  ADD COLUMN receiver_name VARCHAR(64) NOT NULL DEFAULT '' AFTER address_id,
  ADD COLUMN receiver_phone VARCHAR(32) NOT NULL DEFAULT '' AFTER receiver_name,
  ADD COLUMN receiver_address VARCHAR(512) NOT NULL DEFAULT '' AFTER receiver_phone,
  ADD COLUMN ship_company VARCHAR(64) NOT NULL DEFAULT '' AFTER receiver_address,
  ADD COLUMN ship_no VARCHAR(64) NOT NULL DEFAULT '' AFTER ship_company,
  ADD COLUMN shipped_at DATETIME NULL AFTER ship_no,
  ADD KEY idx_fulfill (fulfill_status, created_at);
