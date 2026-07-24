-- 店铺地理坐标 + 本地商家开关 + 门店图
-- 执行：mysql -uhomestead -psecret mymall < scripts/alter-shop-location.sql

USE mymall;

-- 若列已存在会报错，可忽略后继续
ALTER TABLE shops
  ADD COLUMN latitude DECIMAL(10,7) NULL COMMENT '纬度' AFTER address,
  ADD COLUMN longitude DECIMAL(10,7) NULL COMMENT '经度' AFTER latitude,
  ADD COLUMN local_enabled TINYINT NOT NULL DEFAULT 0 COMMENT '是否展示在本地商家' AFTER longitude;

CREATE TABLE IF NOT EXISTS shop_images (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shop_id BIGINT UNSIGNED NOT NULL,
  url VARCHAR(500) NOT NULL,
  sort INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_shop_images_shop (shop_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
