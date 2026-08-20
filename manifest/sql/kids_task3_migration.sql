-- task3 原型相关升级：补充奖励商城、奖励预设、奖励指派和重复兑换字段。

ALTER TABLE kids_reward
    ADD COLUMN circle_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属圈子ID' AFTER id,
    ADD COLUMN image_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '奖励图片地址' AFTER icon,
    ADD COLUMN repeat_rule VARCHAR(32) NOT NULL DEFAULT 'none' COMMENT '重复兑换规则：none/daily/weekly/monthly/custom' AFTER description,
    ADD COLUMN repeat_interval_days INT NOT NULL DEFAULT 0 COMMENT '自定义重复兑换间隔天数' AFTER repeat_rule,
    ADD COLUMN deleted_at DATETIME NULL COMMENT '删除时间' AFTER updated_at,
    ADD KEY idx_circle_deleted (circle_id, deleted_at);

CREATE TABLE IF NOT EXISTS kids_reward_assignee (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖励指派ID',
    reward_id BIGINT UNSIGNED NOT NULL COMMENT '奖励ID',
    kid_id BIGINT UNSIGNED NOT NULL COMMENT '儿童成员ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_reward_kid (reward_id, kid_id),
    KEY idx_kid_id (kid_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids奖励指派';

CREATE TABLE IF NOT EXISTS kids_reward_preset (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖励预设ID',
    title VARCHAR(128) NOT NULL COMMENT '预设奖励标题',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图标标识',
    image_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '奖励图片地址',
    star_cost INT NOT NULL COMMENT '默认所需星星数量',
    description VARCHAR(512) NOT NULL DEFAULT '' COMMENT '预设描述',
    enabled TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否启用预设',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids奖励预设';

ALTER TABLE kids_reward_redeem_record
    ADD COLUMN circle_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '圈子ID' AFTER id,
    ADD COLUMN image_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '奖励图片地址' AFTER icon,
    ADD KEY idx_circle_created (circle_id, created_at);

INSERT INTO kids_reward_preset (id, title, icon, image_url, star_cost, description, enabled) VALUES
    (1, 'Buy toys', '🧩', '', 10, '兑换玩具奖励', 1),
    (2, 'Money', '💰', '', 10, '兑换零花钱', 1),
    (3, 'Pizza', '🍕', '', 15, '兑换披萨', 1),
    (4, '1 hour of game time', '🎮', '', 20, '兑换一小时游戏时间', 1),
    (5, 'See a movie', '🎬', '', 20, '兑换看电影', 1),
    (6, 'Buy new clothes', '👗', '', 20, '兑换新衣服', 1),
    (7, 'camping', '🏕️', '', 30, '兑换露营活动', 1)
ON DUPLICATE KEY UPDATE
    title = VALUES(title),
    icon = VALUES(icon),
    image_url = VALUES(image_url),
    star_cost = VALUES(star_cost),
    description = VALUES(description),
    enabled = VALUES(enabled);
