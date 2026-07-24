-- 售后申请窗口（确认收货后 X 天）；总后台「系统设置」可改
USE mymall;

INSERT INTO sys_config (config_key, config_value, remark)
VALUES ('order_after_sale_days', '7', '确认收货后可申请售后天数')
ON DUPLICATE KEY UPDATE remark=VALUES(remark);
