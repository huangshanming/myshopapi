-- 优质商家种草文章种子（幂等：按标题去重）
-- 执行：mysql --default-character-set=utf8mb4 -u homestead -psecret mymall < scripts/seed-quality-merchant-articles.sql
USE mymall;

-- 分类
INSERT INTO community_article_category (parent_id, name, sort, status)
SELECT 0, '生鲜美食', 10, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM community_article_category WHERE name = '生鲜美食');
INSERT INTO community_article_category (parent_id, name, sort, status)
SELECT 0, '潮流穿搭', 20, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM community_article_category WHERE name = '潮流穿搭');
INSERT INTO community_article_category (parent_id, name, sort, status)
SELECT 0, '数码生活', 30, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM community_article_category WHERE name = '数码生活');
INSERT INTO community_article_category (parent_id, name, sort, status)
SELECT 0, '萌宠日常', 40, 1 FROM DUAL
WHERE NOT EXISTS (SELECT 1 FROM community_article_category WHERE name = '萌宠日常');

SET @cat_fresh := (SELECT id FROM community_article_category WHERE name = '生鲜美食' LIMIT 1);
SET @cat_fashion := (SELECT id FROM community_article_category WHERE name = '潮流穿搭' LIMIT 1);
SET @cat_digital := (SELECT id FROM community_article_category WHERE name = '数码生活' LIMIT 1);
SET @cat_pet := (SELECT id FROM community_article_category WHERE name = '萌宠日常' LIMIT 1);

SET @shop_fresh := (SELECT id FROM shops WHERE name = '本地生鲜连锁自营店' LIMIT 1);
SET @shop_fashion := (SELECT id FROM shops WHERE name = '潮流服饰全国连锁店' LIMIT 1);
SET @shop_xiaomi := (SELECT id FROM shops WHERE name = '小米数码官方旗舰店' LIMIT 1);
SET @shop_pet := (SELECT id FROM shops WHERE name = '示例宠物店' LIMIT 1);

SET @owner_fresh := (SELECT owner_user_id FROM shops WHERE id = @shop_fresh LIMIT 1);
SET @owner_fashion := (SELECT owner_user_id FROM shops WHERE id = @shop_fashion LIMIT 1);
SET @owner_xiaomi := (SELECT owner_user_id FROM shops WHERE id = @shop_xiaomi LIMIT 1);
SET @owner_pet := (SELECT owner_user_id FROM shops WHERE id = @shop_pet LIMIT 1);

-- 1 生鲜
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fresh, @cat_fresh,
  '清晨直采｜今天这波时令蔬果真的绝了',
  '/uploads/shops/seed/1080.jpg',
  CONCAT(
    '<p>凌晨四点半，产地直采车就到了仓库。今天主推的是<strong>云南高原小番茄</strong>和<strong>寿光黄瓜</strong>，口感脆甜、表皮干净，拿回家几乎不用二次清洗。</p>',
    '<p>店员会按「先到先得」分装，建议下班前下单，晚上 8 点前可同城送达。</p>',
    '<ul><li>小番茄：酸甜爆汁，适合当零食</li><li>黄瓜：水分足，凉拌最合适</li><li>今日加赠：新鲜香菜一把</li></ul>',
    '<p>如果你也在意食材新鲜度，来店里闻一闻叶子的清香就懂了。</p>'
  ),
  'approved', 'published', 1, 1280, 96, 210, 860, 42, NOW(), IFNULL(@owner_fresh, 0)
FROM DUAL
WHERE @shop_fresh IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '清晨直采｜今天这波时令蔬果真的绝了');

-- 2 生鲜
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fresh, @cat_fresh,
  '周末家庭餐桌：三道半小时搞定的轻负担菜',
  '/uploads/products/1/tangerine.png',
  CONCAT(
    '<h3>1. 蒜蓉清蒸鲈鱼</h3><p>鱼身划刀，铺姜片，大火蒸 8 分钟，淋热油与生抽，鲜到说话都小声。</p>',
    '<h3>2. 时蔬杂粮饭</h3><p>玉米粒、豌豆、胡萝卜丁和米饭一起焖，孩子也爱吃。</p>',
    '<h3>3. 番茄牛腩汤</h3><p>提前用我们店的牛腩预制包，回锅加番茄，汤色红亮不油腻。</p>',
    '<p>食材都在店内有现货，下单备注「周末餐桌套装」，我们帮你配齐。</p>'
  ),
  'approved', 'published', 0, 860, 71, 150, 520, 33, NOW(), IFNULL(@owner_fresh, 0)
FROM DUAL
WHERE @shop_fresh IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '周末家庭餐桌：三道半小时搞定的轻负担菜');

-- 3 生鲜
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fresh, @cat_fresh,
  '冷链看得见｜从产地到你家餐桌只要 12 小时',
  '/uploads/shops/seed/1080.jpg',
  CONCAT(
    '<p>很多人问我们为什么敢承诺「当日鲜」。答案很简单：<strong>产地采摘 → 产地预冷 → 全程冷链 → 门店分拣</strong>，每一环都有温度记录。</p>',
    '<p>你在 App 下单后，系统会自动匹配最近仓，尽量保证叶菜类当天送达、不捂黄、不发蔫。</p>',
    '<p>本周冷链专线重点保供：菠菜、生菜、草莓、冰鲜鸡胸。</p>'
  ),
  'approved', 'published', 0, 1020, 88, 180, 640, 51, NOW(), IFNULL(@owner_fresh, 0)
FROM DUAL
WHERE @shop_fresh IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '冷链看得见｜从产地到你家餐桌只要 12 小时');

-- 4 生鲜
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fresh, @cat_fresh,
  '夏日冰镇水果拼盘灵感：甜度刚刚好',
  '/uploads/products/1/wine-red.jpg',
  CONCAT(
    '<p>天气一热，冰箱里就不能只有矿泉水。推荐一套「低负担冰镇拼盘」：</p>',
    '<ol><li>西瓜切块 + 薄荷叶</li><li>青提去梗冰镇 1 小时</li><li>黄桃罐头换成新鲜黄桃片</li></ol>',
    '<p>店里水果按糖度分拣，甜度不够的会直接下架，不硬卖。</p>',
    '<p>下单备注「拼盘装」，我们按 2–3 人份帮你配好。</p>'
  ),
  'approved', 'published', 0, 740, 63, 120, 410, 28, NOW(), IFNULL(@owner_fresh, 0)
FROM DUAL
WHERE @shop_fresh IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '夏日冰镇水果拼盘灵感：甜度刚刚好');

-- 5 服饰
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fashion, @cat_fashion,
  '通勤三件套：一件衬衫撑起整周体面感',
  '/uploads/products/1/sweater.png',
  CONCAT(
    '<p>真正好穿的通勤装，不是堆单品，而是<strong>可混搭、耐洗、不起皱</strong>。</p>',
    '<p>本季主推：雾蓝牛津纺衬衫 + 高腰直筒西裤 + 极简小白鞋。办公室、约饭、周末咖啡馆都能穿。</p>',
    '<ul><li>面料：亲肤棉感，透气不闷</li><li>版型：微落肩，显肩窄更利落</li><li>尺码建议：偏宽松，按日常码选即可</li></ul>',
    '<p>门店可预约试穿，线上下单支持 7 天无理由。</p>'
  ),
  'approved', 'published', 1, 1560, 120, 260, 980, 67, NOW(), IFNULL(@owner_fashion, 0)
FROM DUAL
WHERE @shop_fashion IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '通勤三件套：一件衬衫撑起整周体面感');

-- 6 服饰
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fashion, @cat_fashion,
  '城市漫步Look｜轻薄外套这样叠才不土',
  '/uploads/shops/seed/1066.jpg',
  CONCAT(
    '<p>早晚温差大的时候，一件轻薄外套比厚卫衣更聪明。</p>',
    '<p>推荐：短款风衣 + 内搭针织背心 + 阔腿牛仔裤。关键是<strong>颜色不超过三种</strong>，整体更干净。</p>',
    '<p>店员整理了 4 套「出片率」很高的搭配，到店可以说「看漫步 Look 墙」。</p>'
  ),
  'approved', 'published', 0, 920, 77, 140, 560, 39, NOW(), IFNULL(@owner_fashion, 0)
FROM DUAL
WHERE @shop_fashion IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '城市漫步Look｜轻薄外套这样叠才不土');

-- 7 服饰
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_fashion, @cat_fashion,
  '选对基础色，衣柜少一半闲置单品',
  '/uploads/products/1/lipstick1.png',
  CONCAT(
    '<p>买衣服最亏的是「冲动色」。我们把畅销基础色收敛成四类：</p>',
    '<ul><li>米白：提亮显干净</li><li>燕麦灰：最百搭</li><li>墨绿：比黑更有质感</li><li>焦糖棕：秋冬友好</li></ul>',
    '<p>先定色，再选款，衣柜利用率会明显提高。本周基础色专区满 299 减 30。</p>'
  ),
  'approved', 'published', 0, 680, 54, 110, 390, 25, NOW(), IFNULL(@owner_fashion, 0)
FROM DUAL
WHERE @shop_fashion IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '选对基础色，衣柜少一半闲置单品');

-- 8 数码
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_xiaomi, @cat_digital,
  '智能家居入门：从客厅开始的三件套',
  '/uploads/shops/seed/1054.jpg',
  CONCAT(
    '<p>别一上来就全屋智能化。先从客厅三件套做起：</p>',
    '<ol><li><strong>智能音箱</strong>：语音控制灯光与窗帘</li><li><strong>空气净化器</strong>：回南天/雾霾天的刚需</li><li><strong>摄像头</strong>：出差也能看看家里状态</li></ol>',
    '<p>官方旗舰店支持以旧换新，旧设备抵扣后入手更划算。</p>',
    '<p>到店可体验「一句话关灯」场景，真的会上瘾。</p>'
  ),
  'approved', 'published', 1, 2100, 168, 320, 1400, 89, NOW(), IFNULL(@owner_xiaomi, 0)
FROM DUAL
WHERE @shop_xiaomi IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '智能家居入门：从客厅开始的三件套');

-- 9 数码
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_xiaomi, @cat_digital,
  '续航焦虑退散｜出差党的轻装充电方案',
  '/uploads/shops/seed/1054.jpg',
  CONCAT(
    '<p>出差最怕线材乱成一团。我们给商务党整理了「一只手装下」的方案：</p>',
    '<ul><li>65W 氮化镓充电器（双口）</li><li>1.5m 编织线 ×2</li><li>20000mAh 充电宝（支持登机）</li></ul>',
    '<p>重量控制在合理范围，手机、平板、笔记本都能兼顾。门店有现货，可开发票。</p>'
  ),
  'approved', 'published', 0, 1340, 101, 200, 780, 56, NOW(), IFNULL(@owner_xiaomi, 0)
FROM DUAL
WHERE @shop_xiaomi IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '续航焦虑退散｜出差党的轻装充电方案');

-- 10 宠物
INSERT INTO community_article (
  shop_id, category_id, title, cover_url, content,
  audit_status, status, is_top, view_count, like_count, audience_count, read_count, collect_count,
  published_at, created_by
)
SELECT @shop_pet, @cat_pet,
  '新手养猫避坑清单：第一周真正需要买什么',
  '/uploads/products/1/doll.png',
  CONCAT(
    '<p>刚接猫回家别冲动买一堆玩具。第一周真正重要的是：</p>',
    '<ol><li>封闭式猫砂盆 + 低尘猫砂</li><li>缓食碗 / 饮水机</li><li>基础粮（先别频繁换口味）</li><li>抓板：保护沙发也保护爪子</li></ol>',
    '<p>店内提供「新手开箱包」，按体重和年龄配齐，到店还能免费称重建档。</p>',
    '<p>记住：安全感比玩具更重要，给猫一个安静角落，比买十个逗猫棒有用。</p>'
  ),
  'approved', 'published', 0, 990, 82, 160, 610, 47, NOW(), IFNULL(@owner_pet, 0)
FROM DUAL
WHERE @shop_pet IS NOT NULL
  AND NOT EXISTS (SELECT 1 FROM community_article WHERE title = '新手养猫避坑清单：第一周真正需要买什么');

-- 同步配图（每篇 1–2 张，幂等按 article+url）
INSERT INTO community_article_img (article_id, shop_id, url, sort)
SELECT a.id, a.shop_id, a.cover_url, 0
FROM community_article a
WHERE a.title IN (
  '清晨直采｜今天这波时令蔬果真的绝了',
  '周末家庭餐桌：三道半小时搞定的轻负担菜',
  '冷链看得见｜从产地到你家餐桌只要 12 小时',
  '夏日冰镇水果拼盘灵感：甜度刚刚好',
  '通勤三件套：一件衬衫撑起整周体面感',
  '城市漫步Look｜轻薄外套这样叠才不土',
  '选对基础色，衣柜少一半闲置单品',
  '智能家居入门：从客厅开始的三件套',
  '续航焦虑退散｜出差党的轻装充电方案',
  '新手养猫避坑清单：第一周真正需要买什么'
)
AND a.cover_url <> ''
AND NOT EXISTS (
  SELECT 1 FROM community_article_img i WHERE i.article_id = a.id AND i.url = a.cover_url
);
