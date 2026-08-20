-- Clearwave 接口任务标签事实表。

CREATE TABLE IF NOT EXISTS kids_task_tag_definition (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    task_tag_id VARCHAR(128) NOT NULL COMMENT '接口任务标签标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    name VARCHAR(64) NOT NULL COMMENT '标签名称',
    status VARCHAR(32) NOT NULL COMMENT '标签状态',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    deleted_at DATETIME NULL COMMENT '删除时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_task_tag_id (task_tag_id),
    UNIQUE KEY uk_circle_name_active (circle_id, name, status),
    KEY idx_circle_status (circle_id, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口任务标签';
