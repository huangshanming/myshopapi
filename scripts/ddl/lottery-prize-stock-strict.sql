-- 奖品：强控库存开关
USE mymall;
ALTER TABLE lottery_prizes
  ADD COLUMN stock_strict TINYINT NOT NULL DEFAULT 0 COMMENT '1=强控库存 0=不强制' AFTER stock;
