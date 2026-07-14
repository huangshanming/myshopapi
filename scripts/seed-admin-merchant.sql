-- 种子账号：平台超管 + 示例商家店主与已审核店铺
-- 密码均为 123456
-- 执行：mysql -u homestead -p mymall < scripts/seed-admin-merchant.sql

USE mymall;

SET @pwd := MD5(CONCAT('123456', 'this is my mall'));

INSERT INTO users (mobile, password, nickname, status, role)
VALUES ('13900000001', @pwd, '平台超管', 1, 'platform_admin')
ON DUPLICATE KEY UPDATE role = 'platform_admin', password = @pwd, nickname = '平台超管', status = 1;

INSERT INTO users (mobile, password, nickname, status, role)
VALUES ('13900000002', @pwd, '示例店主', 1, 'merchant_owner')
ON DUPLICATE KEY UPDATE role = 'merchant_owner', password = @pwd, nickname = '示例店主', status = 1;

SET @owner_id := (SELECT id FROM users WHERE mobile = '13900000002' LIMIT 1);

INSERT INTO shops (name, logo, contact_name, contact_phone, description, owner_user_id, status)
SELECT '示例宠物店', '', '店主', '13900000002', '种子店铺', @owner_id, 'approved'
FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM shops WHERE owner_user_id = @owner_id AND name = '示例宠物店');

SET @shop_id := (SELECT id FROM shops WHERE owner_user_id = @owner_id ORDER BY id ASC LIMIT 1);

INSERT INTO shop_members (shop_id, user_id, member_role)
SELECT @shop_id, @owner_id, 'owner'
FROM DUAL
WHERE @shop_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM shop_members WHERE shop_id = @shop_id AND user_id = @owner_id);

INSERT INTO product_categories (parent_id, name, sort_order, level, is_show)
SELECT 0, '猫粮', 1, 1, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM product_categories WHERE name = '猫粮');

SET @cat_id := (SELECT id FROM product_categories WHERE name = '猫粮' LIMIT 1);

INSERT INTO products (shop_id, product_no, name, sale_price, stock, category_id, status, pet_type)
SELECT @shop_id, 'P-SEED-001', '种子商品-营养猫粮', 99.00, 100, IFNULL(@cat_id, 1), 'on_sale', 'cat'
FROM DUAL
WHERE @shop_id IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM products WHERE product_no = 'P-SEED-001');
