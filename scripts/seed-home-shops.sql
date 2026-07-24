-- 首页品牌/优质商户种子（幂等）
-- 门头图需先落盘：bash scripts/seed-home-shop-images.sh
-- 密码均为 123456
-- 执行：mysql -u homestead -p mymall < scripts/seed-home-shops.sql

USE mymall;

SET @pwd := MD5(CONCAT('123456', 'this is my mall'));

-- 店主账号
INSERT INTO users (mobile, password, nickname, status, role)
VALUES ('13900000011', @pwd, '小米旗舰店主', 1, 'merchant_owner')
ON DUPLICATE KEY UPDATE role = 'merchant_owner', password = @pwd, nickname = '小米旗舰店主', status = 1;

INSERT INTO users (mobile, password, nickname, status, role)
VALUES ('13900000012', @pwd, '生鲜连锁店主', 1, 'merchant_owner')
ON DUPLICATE KEY UPDATE role = 'merchant_owner', password = @pwd, nickname = '生鲜连锁店主', status = 1;

INSERT INTO users (mobile, password, nickname, status, role)
VALUES ('13900000013', @pwd, '潮流服饰店主', 1, 'merchant_owner')
ON DUPLICATE KEY UPDATE role = 'merchant_owner', password = @pwd, nickname = '潮流服饰店主', status = 1;

SET @owner11 := (SELECT id FROM users WHERE mobile = '13900000011' LIMIT 1);
SET @owner12 := (SELECT id FROM users WHERE mobile = '13900000012' LIMIT 1);
SET @owner13 := (SELECT id FROM users WHERE mobile = '13900000013' LIMIT 1);

-- 小米数码官方旗舰店
INSERT INTO shops (
  name, logo, contact_name, contact_phone, description, category,
  province, city, district, address,
  storefront_image, owner_user_id, status
)
SELECT
  '小米数码官方旗舰店', '/uploads/shops/seed/1054.jpg', '店主', '13900000011', '品牌数码官方直营', '品牌直营',
  '湖南省', '长沙市', '岳麓区', '岳麓大道88号',
  '/uploads/shops/seed/1054.jpg', @owner11, 'approved'
FROM DUAL
WHERE @owner11 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shops WHERE name = '小米数码官方旗舰店');

-- 本地生鲜连锁自营店
INSERT INTO shops (
  name, logo, contact_name, contact_phone, description, category,
  province, city, district, address,
  storefront_image, owner_user_id, status
)
SELECT
  '本地生鲜连锁自营店', '/uploads/shops/seed/1080.jpg', '店主', '13900000012', '生鲜连锁自营，全城配送', '连锁品牌',
  '湖南省', '长沙市', '芙蓉区', '五一大道66号',
  '/uploads/shops/seed/1080.jpg', @owner12, 'approved'
FROM DUAL
WHERE @owner12 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shops WHERE name = '本地生鲜连锁自营店');

-- 潮流服饰全国连锁店
INSERT INTO shops (
  name, logo, contact_name, contact_phone, description, category,
  province, city, district, address,
  storefront_image, owner_user_id, status
)
SELECT
  '潮流服饰全国连锁店', '/uploads/shops/seed/1066.jpg', '店主', '13900000013', '线下门店同步，全场促销', '线下门店同步',
  '广东省', '广州市', '天河区', '天河路100号',
  '/uploads/shops/seed/1066.jpg', @owner13, 'approved'
FROM DUAL
WHERE @owner13 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shops WHERE name = '潮流服饰全国连锁店');

-- 若店铺已存在则补齐 logo / 门头图 / 类目 / 城市
UPDATE shops
SET storefront_image = '/uploads/shops/seed/1054.jpg',
    logo = '/uploads/shops/seed/1054.jpg',
    category = IF(category = '' OR category IS NULL, '品牌直营', category),
    description = IF(description = '' OR description IS NULL, '品牌数码官方直营', description),
    province = '湖南省',
    city = '长沙市',
    district = IF(district = '' OR district IS NULL, '岳麓区', district),
    address = IF(address = '' OR address IS NULL, '岳麓大道88号', address)
WHERE name = '小米数码官方旗舰店';

UPDATE shops
SET storefront_image = '/uploads/shops/seed/1080.jpg',
    logo = '/uploads/shops/seed/1080.jpg',
    category = IF(category = '' OR category IS NULL, '连锁品牌', category),
    description = IF(description = '' OR description IS NULL, '生鲜连锁自营，全城配送', description),
    province = '湖南省',
    city = '长沙市',
    district = IF(district = '' OR district IS NULL, '芙蓉区', district),
    address = IF(address = '' OR address IS NULL, '五一大道66号', address)
WHERE name = '本地生鲜连锁自营店';

UPDATE shops
SET storefront_image = '/uploads/shops/seed/1066.jpg',
    logo = '/uploads/shops/seed/1066.jpg',
    category = IF(category = '' OR category IS NULL, '线下门店同步', category),
    description = IF(description = '' OR description IS NULL, '线下门店同步，全场促销', description),
    province = '广东省',
    city = '广州市',
    district = IF(district = '' OR district IS NULL, '天河区', district),
    address = IF(address = '' OR address IS NULL, '天河路100号', address)
WHERE name = '潮流服饰全国连锁店';

SET @shop11 := (SELECT id FROM shops WHERE name = '小米数码官方旗舰店' LIMIT 1);
SET @shop12 := (SELECT id FROM shops WHERE name = '本地生鲜连锁自营店' LIMIT 1);
SET @shop13 := (SELECT id FROM shops WHERE name = '潮流服饰全国连锁店' LIMIT 1);

INSERT INTO shop_members (shop_id, user_id, member_role)
SELECT @shop11, @owner11, 'owner'
FROM DUAL
WHERE @shop11 IS NOT NULL AND @owner11 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shop_members WHERE shop_id = @shop11 AND user_id = @owner11);

INSERT INTO shop_members (shop_id, user_id, member_role)
SELECT @shop12, @owner12, 'owner'
FROM DUAL
WHERE @shop12 IS NOT NULL AND @owner12 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shop_members WHERE shop_id = @shop12 AND user_id = @owner12);

INSERT INTO shop_members (shop_id, user_id, member_role)
SELECT @shop13, @owner13, 'owner'
FROM DUAL
WHERE @shop13 IS NOT NULL AND @owner13 IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shop_members WHERE shop_id = @shop13 AND user_id = @owner13);

-- 本地商家坐标（幂等）
UPDATE shops SET latitude=28.228209, longitude=112.938814, local_enabled=1
WHERE name = '小米数码官方旗舰店';
UPDATE shops SET latitude=28.201500, longitude=112.982000, local_enabled=1
WHERE name = '本地生鲜连锁自营店';
UPDATE shops SET latitude=28.215000, longitude=112.950000, local_enabled=1
WHERE name = '示例宠物店';
UPDATE shops SET latitude=23.129110, longitude=113.264385, local_enabled=1
WHERE name = '潮流服饰全国连锁店';
