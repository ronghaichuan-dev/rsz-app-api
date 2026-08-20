-- Clearwave 合同任务定义、分配和 occurrence 事实表。
-- occurrence 使用 task/member/scheduled_date 复合自然键，禁止对外暴露 surrogate ID。

CREATE TABLE IF NOT EXISTS kids_contract_task (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    task_id VARCHAR(128) NOT NULL COMMENT '合同任务标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    title VARCHAR(160) NOT NULL COMMENT '任务标题',
    notes TEXT NULL COMMENT '任务备注',
    emoji VARCHAR(64) NULL COMMENT '任务图标',
    stars INT UNSIGNED NOT NULL COMMENT '星星数量',
    start_date DATE NOT NULL COMMENT '系列开始日期',
    zone_id VARCHAR(128) NOT NULL COMMENT '时区',
    repeat_rule JSON NOT NULL COMMENT '重复规则',
    end_rule JSON NOT NULL COMMENT '结束规则',
    time_limit_minute_of_day SMALLINT UNSIGNED NULL COMMENT '时限分钟',
    reminder_config JSON NULL COMMENT '提醒配置',
    photo_required TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '是否要求图片凭证',
    task_tag_id VARCHAR(128) NULL COMMENT '合同任务标签标识',
    series_revision BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '系列修订号',
    future_effective_from_date DATE NOT NULL COMMENT '当前未来生效日',
    status VARCHAR(32) NOT NULL COMMENT '任务状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_task_id (task_id),
    KEY idx_circle_status (circle_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同任务定义';

CREATE TABLE IF NOT EXISTS kids_contract_task_assignment (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    task_id VARCHAR(128) NOT NULL COMMENT '合同任务标识',
    member_id VARCHAR(128) NOT NULL COMMENT '合同成员标识',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_task_member (task_id, member_id),
    KEY idx_member_task (member_id, task_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同任务成员分配';

CREATE TABLE IF NOT EXISTS kids_contract_task_occurrence (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '内部主键',
    circle_id VARCHAR(128) NOT NULL COMMENT '合同圈子标识',
    task_id VARCHAR(128) NOT NULL COMMENT '合同任务标识',
    member_id VARCHAR(128) NOT NULL COMMENT '合同成员标识',
    scheduled_date DATE NOT NULL COMMENT '预定日',
    zone_id VARCHAR(128) NOT NULL COMMENT '时区',
    definition_revision BIGINT UNSIGNED NOT NULL COMMENT '定义修订号',
    title_snapshot VARCHAR(160) NOT NULL COMMENT '标题快照',
    notes_snapshot TEXT NULL COMMENT '备注快照',
    emoji_snapshot VARCHAR(64) NULL COMMENT '图标快照',
    stars_snapshot INT UNSIGNED NOT NULL COMMENT '星星快照',
    photo_required_snapshot TINYINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '图片要求快照',
    task_tag_id_snapshot VARCHAR(128) NULL COMMENT '标签快照',
    state VARCHAR(32) NOT NULL COMMENT 'occurrence 状态',
    completion_id VARCHAR(128) NULL COMMENT '完成事实标识',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_task_member_date (task_id, member_id, scheduled_date),
    KEY idx_circle_member_date (circle_id, member_id, scheduled_date),
    KEY idx_task_state_date (task_id, state, scheduled_date)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids合同任务 occurrence';
