-- Clearwave 接口奖励、兑换、冷却与通知 outbox 事实表。
-- 所有表仅由 kids 微服务访问，接口对外标识使用 opaque ID。

CREATE TABLE IF NOT EXISTS kids_reward_definition (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    reward_id VARCHAR(128) NOT NULL COMMENT '接口奖励标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    title VARCHAR(160) NOT NULL COMMENT '奖励标题',
    description TEXT NOT NULL COMMENT '奖励描述',
    visual JSON NOT NULL COMMENT '奖励视觉引用',
    stars_required INT UNSIGNED NOT NULL COMMENT '兑换所需星星',
    repeat_rule VARCHAR(32) NOT NULL COMMENT '重复兑换规则',
    cooldown_days INT UNSIGNED NULL COMMENT '自定义冷却天数',
    zone_id VARCHAR(128) NOT NULL COMMENT '时区标识',
    status VARCHAR(32) NOT NULL COMMENT '奖励状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_reward_id (reward_id),
    KEY idx_circle_status (circle_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口奖励定义';

CREATE TABLE IF NOT EXISTS kids_reward_assignment (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    reward_id VARCHAR(128) NOT NULL COMMENT '接口奖励标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_reward_member (reward_id, member_id),
    KEY idx_member_reward (member_id, reward_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口奖励成员分配';

CREATE TABLE IF NOT EXISTS kids_reward_cooldown (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    reward_id VARCHAR(128) NOT NULL COMMENT '接口奖励标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    cooldown_until_at DATETIME NULL COMMENT '冷却结束时间',
    last_redeemed_at DATETIME NOT NULL COMMENT '最近兑换时间',
    permanently_unavailable TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否永久不可兑换',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_reward_member (reward_id, member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口奖励冷却';

CREATE TABLE IF NOT EXISTS kids_reward_exchange (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    exchange_id VARCHAR(128) NOT NULL COMMENT '接口兑换标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    member_name_snapshot VARCHAR(120) NOT NULL COMMENT '成员名称快照',
    member_avatar_snapshot JSON NOT NULL COMMENT '成员头像快照',
    reward_id VARCHAR(128) NOT NULL COMMENT '接口奖励标识',
    reward_title_snapshot VARCHAR(160) NOT NULL COMMENT '奖励标题快照',
    reward_visual_snapshot JSON NOT NULL COMMENT '奖励视觉快照',
    stars_deducted_snapshot INT UNSIGNED NOT NULL COMMENT '扣减星星快照',
    reward_repeat_rule_snapshot VARCHAR(32) NOT NULL COMMENT '重复规则快照',
    reward_cooldown_days_snapshot INT UNSIGNED NULL COMMENT '冷却天数快照',
    cooldown_until_at_snapshot DATETIME NULL COMMENT '冷却结束快照',
    permanently_unavailable_snapshot TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '永久不可兑换快照',
    ledger_id VARCHAR(128) NOT NULL COMMENT '账本流水标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '提交序列',
    exchanged_at DATETIME NOT NULL COMMENT '兑换时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_exchange_id (exchange_id),
    KEY idx_circle_sequence (circle_id, commit_sequence),
    KEY idx_member_created (member_id, exchanged_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口奖励兑换审计';

CREATE TABLE IF NOT EXISTS kids_notification_outbox (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    notification_id VARCHAR(128) NOT NULL COMMENT '接口通知标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    account_id VARCHAR(128) NOT NULL COMMENT '接收账号标识',
    exchange_id VARCHAR(128) NOT NULL COMMENT '接口兑换标识',
    event_type VARCHAR(64) NOT NULL COMMENT '通知事件类型',
    payload JSON NOT NULL COMMENT '通知载荷',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '提交序列',
    status VARCHAR(32) NOT NULL COMMENT '投递状态',
    attempt_count INT UNSIGNED NOT NULL DEFAULT 0 COMMENT '已尝试投递次数',
    next_attempt_at DATETIME NULL COMMENT '下次允许投递时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    PRIMARY KEY (id),
    UNIQUE KEY uk_notification_id (notification_id),
    KEY idx_account_status (account_id, status),
    KEY idx_exchange_id (exchange_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口通知出站箱';
