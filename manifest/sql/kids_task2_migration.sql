-- task2 原型相关升级：补充群组管理、成员绑定、上传记录、通知设备和奖励兑换记录。

ALTER TABLE kids_circle
    ADD COLUMN deleted_at DATETIME NULL COMMENT '删除时间' AFTER updated_at,
    ADD KEY idx_deleted_at (deleted_at);

ALTER TABLE kids_circle_user
    ADD COLUMN deleted_at DATETIME NULL COMMENT '删除时间' AFTER updated_at,
    ADD COLUMN left_at DATETIME NULL COMMENT '退出时间' AFTER deleted_at,
    ADD KEY idx_user_deleted (user_id, deleted_at),
    ADD KEY idx_circle_role_deleted (circle_id, role, deleted_at);

ALTER TABLE kids_family_member
    ADD COLUMN bind_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '绑定用户ID' AFTER owner,
    ADD COLUMN bound_at DATETIME NULL COMMENT '绑定时间' AFTER bind_user_id,
    ADD COLUMN deleted_at DATETIME NULL COMMENT '删除时间' AFTER updated_at,
    ADD KEY idx_bind_user_id (bind_user_id),
    ADD KEY idx_circle_deleted (circle_id, deleted_at);

CREATE TABLE IF NOT EXISTS kids_upload_file (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '上传文件ID',
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '上传用户ID',
    member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '绑定成员ID',
    usage_type VARCHAR(32) NOT NULL DEFAULT '' COMMENT '文件用途：avatar/task_photo/reward',
    file_name VARCHAR(255) NOT NULL DEFAULT '' COMMENT '文件名称',
    file_url VARCHAR(512) NOT NULL COMMENT '文件访问地址',
    content_type VARCHAR(128) NOT NULL DEFAULT '' COMMENT '文件类型',
    file_size BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '文件大小字节数',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_user_usage (user_id, usage_type),
    KEY idx_member_id (member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids上传文件';

CREATE TABLE IF NOT EXISTS kids_device_notification (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '设备通知ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    device_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '设备ID',
    platform VARCHAR(32) NOT NULL DEFAULT '' COMMENT '平台：ios/android/web',
    device_token VARCHAR(512) NOT NULL DEFAULT '' COMMENT '推送设备令牌',
    authorized TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已授权通知',
    task_enabled TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否开启任务提醒',
    reward_enabled TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否开启奖励提醒',
    member_enabled TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否开启成员提醒',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_user_device (user_id, device_id),
    KEY idx_device_token (device_token(191))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids设备通知设置';

CREATE TABLE IF NOT EXISTS kids_reward_redeem_record (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖励兑换记录ID',
    reward_id BIGINT UNSIGNED NOT NULL COMMENT '奖励ID',
    kid_id BIGINT UNSIGNED NOT NULL COMMENT '儿童成员ID',
    user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '兑换用户ID',
    title VARCHAR(128) NOT NULL COMMENT '奖励标题',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图标标识',
    star_cost INT NOT NULL COMMENT '消耗星星数量',
    remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '兑换备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_kid_created (kid_id, created_at),
    KEY idx_reward_id (reward_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids奖励兑换记录';
