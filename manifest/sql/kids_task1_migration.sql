-- task1 原型相关升级：补充任务标签、重复、提醒、时间限制、说明和软删除字段。

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

ALTER TABLE kids_task
    ADD COLUMN note VARCHAR(512) NOT NULL DEFAULT '' COMMENT '任务说明' AFTER icon,
    ADD COLUMN repeat_end_type VARCHAR(16) NOT NULL DEFAULT 'never' COMMENT '重复结束类型：never/date/count' AFTER repeat_rule,
    ADD COLUMN repeat_end_date DATE NULL COMMENT '重复结束日期' AFTER repeat_end_type,
    ADD COLUMN repeat_end_count INT NOT NULL DEFAULT 0 COMMENT '重复结束次数' AFTER repeat_end_date,
    ADD COLUMN time_limit_type VARCHAR(16) NOT NULL DEFAULT 'all_day' COMMENT '时间限制类型：all_day/range' AFTER repeat_end_count,
    ADD COLUMN time_limit_start CHAR(5) NOT NULL DEFAULT '' COMMENT '开始时间，格式HH:mm' AFTER time_limit_type,
    ADD COLUMN time_limit_end CHAR(5) NOT NULL DEFAULT '' COMMENT '结束时间，格式HH:mm' AFTER time_limit_start,
    ADD COLUMN reminder_type VARCHAR(16) NOT NULL DEFAULT 'none' COMMENT '提醒类型：none/at_time/before_start' AFTER time_limit_end,
    ADD COLUMN reminder_at CHAR(5) NOT NULL DEFAULT '' COMMENT '提醒时间，格式HH:mm' AFTER reminder_type,
    ADD COLUMN reminder_offset_minutes INT NOT NULL DEFAULT 0 COMMENT '提前提醒分钟数' AFTER reminder_at,
    ADD COLUMN tag_id BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '标签ID' AFTER need_photo_proof,
    ADD COLUMN deleted_at DATETIME NULL COMMENT '删除时间' AFTER completed_at,
    ADD KEY idx_tag_date (tag_id, task_date),
    ADD KEY idx_deleted_date (deleted_at, task_date);
