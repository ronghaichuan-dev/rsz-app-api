-- kids 登录和业务表。每个 kids 微服务只拥有并访问自己的 rslytics_kids_* 数据库。

CREATE TABLE IF NOT EXISTS kids_circle (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '圈子ID',
    name VARCHAR(64) NOT NULL COMMENT '圈子名称',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '圈子图标标识',
    owner_user_id BIGINT UNSIGNED NOT NULL COMMENT '创建者用户ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_owner_user_id (owner_user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids圈子';

CREATE TABLE IF NOT EXISTS kids_circle_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '圈子用户关系ID',
    circle_id BIGINT UNSIGNED NOT NULL COMMENT '圈子ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    role VARCHAR(16) NOT NULL COMMENT '圈子角色：admin/member',
    member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '绑定的家庭成员ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_circle_user (circle_id, user_id),
    KEY idx_user_id (user_id),
    KEY idx_member_id (member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids圈子用户关系';

CREATE TABLE IF NOT EXISTS kids_invite_code (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '邀请码ID',
    code CHAR(6) NOT NULL COMMENT '六位邀请码',
    circle_id BIGINT UNSIGNED NOT NULL COMMENT '圈子ID',
    inviter_user_id BIGINT UNSIGNED NOT NULL COMMENT '邀请人用户ID',
    invite_role VARCHAR(16) NOT NULL COMMENT '邀请加入角色：admin/member',
    target_member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标家庭成员ID，0表示不指定',
    expired_at DATETIME NOT NULL COMMENT '过期时间',
    used_at DATETIME NULL COMMENT '使用时间',
    used_by_user_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '使用者用户ID',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_code (code),
    KEY idx_circle_role (circle_id, invite_role),
    KEY idx_target_member_id (target_member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids邀请码';

CREATE TABLE IF NOT EXISTS kids_user (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    device_id VARCHAR(128) NOT NULL DEFAULT '' COMMENT '最近登录设备ID',
    provider VARCHAR(32) NOT NULL DEFAULT 'guest' COMMENT '当前登录方式：guest/google/apple',
    email VARCHAR(255) NOT NULL DEFAULT '' COMMENT '授权服务商返回的邮箱',
    nickname VARCHAR(128) NOT NULL DEFAULT '' COMMENT '昵称',
    avatar VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像地址',
    is_guest TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否游客账号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_device_id (device_id),
    KEY idx_is_guest (is_guest)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids用户';

CREATE TABLE IF NOT EXISTS kids_user_auth (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '授权记录ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT 'kids用户ID',
    provider VARCHAR(32) NOT NULL COMMENT '授权服务商：google/apple',
    open_id VARCHAR(255) NOT NULL COMMENT '服务商开放ID或主体标识',
    email VARCHAR(255) NOT NULL DEFAULT '' COMMENT '服务商邮箱',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_provider_open_id (provider, open_id),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids用户授权绑定';

CREATE TABLE IF NOT EXISTS kids_user_token (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '令牌ID',
    user_id BIGINT UNSIGNED NOT NULL COMMENT 'kids用户ID',
    token VARCHAR(255) NOT NULL COMMENT '访问令牌',
    expired_at DATETIME NOT NULL COMMENT '过期时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_token (token),
    KEY idx_user_id (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids用户访问令牌';

CREATE TABLE IF NOT EXISTS kids_family_member (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '家庭成员ID',
    circle_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '所属圈子ID',
    name VARCHAR(64) NOT NULL COMMENT '显示名称',
    gender VARCHAR(16) NOT NULL DEFAULT '' COMMENT '性别：male/female',
    avatar VARCHAR(512) NOT NULL DEFAULT '' COMMENT '头像地址或预设标识',
    avatar_style VARCHAR(64) NOT NULL DEFAULT '' COMMENT '虚拟形象风格标识',
    relation VARCHAR(64) NOT NULL DEFAULT '' COMMENT '家庭关系',
    owner TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否家庭拥有者',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_circle_id (circle_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids家庭成员';

CREATE TABLE IF NOT EXISTS kids_task_tag (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务标签ID',
    name VARCHAR(32) NOT NULL COMMENT '标签名称',
    color VARCHAR(32) NOT NULL DEFAULT '' COMMENT '标签颜色',
    sort_order INT NOT NULL DEFAULT 0 COMMENT '排序值',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_deleted_sort (deleted_at, sort_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids任务标签';

CREATE TABLE IF NOT EXISTS kids_task (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务ID',
    title VARCHAR(128) NOT NULL COMMENT '任务标题',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图标标识',
    note VARCHAR(512) NOT NULL DEFAULT '' COMMENT '任务说明',
    star INT NOT NULL COMMENT '星星奖励数量',
    task_date DATE NOT NULL COMMENT '计划日期',
    completion_mode VARCHAR(32) NOT NULL COMMENT '完成模式：single/rotation/anyone/everyone',
    repeat_rule VARCHAR(64) NOT NULL DEFAULT 'none' COMMENT '重复规则',
    repeat_end_type VARCHAR(16) NOT NULL DEFAULT 'never' COMMENT '重复结束类型：never/date/count',
    repeat_end_date DATE NULL COMMENT '重复结束日期',
    repeat_end_count INT NOT NULL DEFAULT 0 COMMENT '重复结束次数',
    time_limit_type VARCHAR(16) NOT NULL DEFAULT 'all_day' COMMENT '时间限制类型：all_day/range',
    time_limit_start CHAR(5) NOT NULL DEFAULT '' COMMENT '开始时间，格式HH:mm',
    time_limit_end CHAR(5) NOT NULL DEFAULT '' COMMENT '结束时间，格式HH:mm',
    reminder_type VARCHAR(16) NOT NULL DEFAULT 'none' COMMENT '提醒类型：none/at_time/before_start',
    reminder_at CHAR(5) NOT NULL DEFAULT '' COMMENT '提醒时间，格式HH:mm',
    reminder_offset_minutes INT NOT NULL DEFAULT 0 COMMENT '提前提醒分钟数',
    need_photo_proof TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否需要照片凭证',
    tag_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '标签ID',
    completed TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已完成',
    completed_by BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '完成任务的儿童成员ID',
    photo_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '照片凭证地址',
    completed_at DATETIME NULL COMMENT '完成时间',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_task_date (task_date),
    KEY idx_completed_by (completed_by),
    KEY idx_tag_date (tag_id, task_date),
    KEY idx_deleted_date (deleted_at, task_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids任务';

CREATE TABLE IF NOT EXISTS kids_task_assignee (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务分配ID',
    task_id BIGINT UNSIGNED NOT NULL COMMENT '任务ID',
    kid_id BIGINT UNSIGNED NOT NULL COMMENT '儿童成员ID',
    assignee_order INT NOT NULL DEFAULT 0 COMMENT '分配顺序，用于轮流模式',
    completed TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '该儿童是否已完成',
    photo_url VARCHAR(512) NOT NULL DEFAULT '' COMMENT '该儿童照片凭证地址',
    completed_at DATETIME NULL COMMENT '该儿童完成时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_task_kid (task_id, kid_id),
    KEY idx_kid_id (kid_id),
    KEY idx_task_completed (task_id, completed)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids任务分配';

CREATE TABLE IF NOT EXISTS kids_reward (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '奖励ID',
    title VARCHAR(128) NOT NULL COMMENT '奖励标题',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图标标识',
    star_cost INT NOT NULL COMMENT '所需星星数量',
    stock INT NOT NULL DEFAULT -1 COMMENT '可用库存，-1表示不限量',
    description VARCHAR(512) NOT NULL DEFAULT '' COMMENT '奖励描述',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids奖励';

CREATE TABLE IF NOT EXISTS kids_star_record (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '星星流水ID',
    kid_id BIGINT UNSIGNED NOT NULL COMMENT '儿童成员ID',
    change_amount INT NOT NULL COMMENT '星星变动数量',
    balance INT NOT NULL COMMENT '变动后余额',
    record_type VARCHAR(32) NOT NULL COMMENT '流水类型：task/reward/adjustment',
    title VARCHAR(255) NOT NULL COMMENT '流水标题',
    remark VARCHAR(512) NOT NULL DEFAULT '' COMMENT '流水备注',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    KEY idx_kid_id_created_at (kid_id, created_at),
    KEY idx_record_type (record_type)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids星星流水';

CREATE TABLE IF NOT EXISTS kids_notification (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '通知ID',
    member_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '目标成员ID，0表示全家庭',
    notification_type VARCHAR(32) NOT NULL COMMENT '通知类型',
    title VARCHAR(128) NOT NULL COMMENT '通知标题',
    content VARCHAR(512) NOT NULL DEFAULT '' COMMENT '通知内容',
    is_read TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否已读',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_member_id_read (member_id, is_read)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids通知';


CREATE TABLE IF NOT EXISTS kids_task_preset (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '任务预设ID',
    title VARCHAR(128) NOT NULL COMMENT '预设任务标题',
    icon VARCHAR(128) NOT NULL DEFAULT '' COMMENT '图标标识',
    star INT NOT NULL COMMENT '默认星星数量',
    need_photo TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否建议照片凭证',
    description VARCHAR(512) NOT NULL DEFAULT '' COMMENT '预设描述',
    enabled TINYINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '是否启用预设',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    KEY idx_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids任务预设';

INSERT INTO kids_task_preset (id, title, icon, star, need_photo, description, enabled) VALUES
    (1, 'Remove toys', 'toy', 2, 0, 'Put toys back in place', 1),
    (2, 'Brush teeth', 'toothbrush', 1, 0, 'Brush teeth morning or night', 1),
    (3, 'Remove shoes', 'shoes', 1, 0, 'Put shoes neatly by the door', 1),
    (4, 'Make the bed', 'bed', 2, 1, 'Make the bed and submit photo proof', 1),
    (5, 'Read book', 'book', 3, 0, 'Read a book independently', 1),
    (6, 'Wash the dishes', 'dishes', 2, 0, 'Help wash dishes', 1)
ON DUPLICATE KEY UPDATE
    title = VALUES(title),
    icon = VALUES(icon),
    star = VALUES(star),
    need_photo = VALUES(need_photo),
    description = VALUES(description),
    enabled = VALUES(enabled);
