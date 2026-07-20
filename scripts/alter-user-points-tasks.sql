-- 积分账本 + 任务中心
-- mysql --default-character-set=utf8mb4 -uhomestead -psecret mymall < scripts/alter-user-points-tasks.sql
USE mymall;

CREATE TABLE IF NOT EXISTS user_points (
    user_id BIGINT UNSIGNED NOT NULL PRIMARY KEY COMMENT '用户ID',
    points BIGINT NOT NULL DEFAULT 0 COMMENT '当前积分',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分账户';

CREATE TABLE IF NOT EXISTS user_point_logs (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    change_type VARCHAR(32) NOT NULL DEFAULT 'task_claim' COMMENT 'task_claim/admin_adjust',
    delta INT NOT NULL COMMENT '变动积分(可负)',
    points_after BIGINT NOT NULL DEFAULT 0,
    remark VARCHAR(255) NOT NULL DEFAULT '',
    ref_type VARCHAR(32) NOT NULL DEFAULT '',
    ref_id BIGINT UNSIGNED NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    KEY idx_user_created (user_id, created_at),
    KEY idx_ref (ref_type, ref_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户积分流水';

CREATE TABLE IF NOT EXISTS task_definitions (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(64) NOT NULL COMMENT '任务编码',
    title VARCHAR(100) NOT NULL,
    description VARCHAR(500) NOT NULL DEFAULT '',
    icon VARCHAR(64) NOT NULL DEFAULT '',
    period VARCHAR(16) NOT NULL DEFAULT 'daily' COMMENT 'daily|once',
    enabled TINYINT NOT NULL DEFAULT 1,
    reward_points INT NOT NULL DEFAULT 0 COMMENT '0=不奖励积分',
    target_count INT NOT NULL DEFAULT 1 COMMENT '完成所需进度',
    daily_limit INT NOT NULL DEFAULT 0 COMMENT '每日最多领取次数,0=不限; once任务忽略',
    sort INT NOT NULL DEFAULT 0,
    rules_json VARCHAR(1000) NOT NULL DEFAULT '{}',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_code (code),
    KEY idx_enabled_sort (enabled, sort)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务定义';

CREATE TABLE IF NOT EXISTS user_task_progress (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    task_code VARCHAR(64) NOT NULL,
    biz_date DATE NOT NULL COMMENT '日任务按自然日; once=1970-01-01',
    progress INT NOT NULL DEFAULT 0,
    claim_count INT NOT NULL DEFAULT 0 COMMENT '当日已领取次数',
    status VARCHAR(16) NOT NULL DEFAULT 'ongoing' COMMENT 'ongoing|claimable|claimed',
    claimed_at DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY uk_user_task_date (user_id, task_code, biz_date),
    KEY idx_user_status (user_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户任务进度';

CREATE TABLE IF NOT EXISTS user_task_dedupe (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL,
    task_code VARCHAR(64) NOT NULL,
    biz_date DATE NOT NULL,
    ref_key VARCHAR(64) NOT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE KEY uk_dedupe (user_id, task_code, biz_date, ref_key)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务事件去重';

INSERT INTO task_definitions (code, title, description, icon, period, enabled, reward_points, target_count, daily_limit, sort, rules_json) VALUES
('daily_checkin', '每日签到', '每天签到一次可领取积分', 'checkin', 'daily', 1, 5, 1, 1, 10, '{}'),
('publish_article', '发布种草笔记', '发布笔记并通过审核后可领取；可配置每日奖励次数', 'article', 'daily', 1, 20, 1, 1, 20, '{}'),
('comment_article', '发表评论', '在种草社区发表有效评论', 'comment', 'daily', 1, 5, 3, 1, 30, '{}'),
('like_article', '点赞文章', '为喜欢的笔记点赞', 'like', 'daily', 1, 3, 5, 1, 40, '{}'),
('favorite_article', '收藏文章', '收藏种草笔记', 'star', 'daily', 1, 3, 3, 1, 50, '{}'),
('browse_products', '浏览商品', '浏览不同商品详情页', 'browse', 'daily', 1, 5, 5, 1, 60, '{}'),
('place_order', '下单购物', '支付成功一笔订单', 'order', 'daily', 1, 30, 1, 1, 70, '{}'),
('first_profile', '完善资料', '完善昵称与头像', 'profile', 'once', 1, 50, 1, 0, 80, '{}'),
('first_favorite_product', '首次收藏商品', '完成第一次商品收藏', 'favorite', 'once', 1, 20, 1, 0, 90, '{}'),
('invite_placeholder', '邀请好友', '邀请好友注册（即将开放）', 'invite', 'daily', 0, 100, 1, 1, 100, '{}')
ON DUPLICATE KEY UPDATE title=VALUES(title), description=VALUES(description), sort=VALUES(sort);
