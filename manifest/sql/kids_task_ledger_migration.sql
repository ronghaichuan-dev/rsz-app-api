-- Clearwave 接口任务完成、取消、星星流水与余额事实表。

CREATE TABLE IF NOT EXISTS kids_task_completion (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    completion_id VARCHAR(128) NOT NULL COMMENT '接口完成标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    task_id VARCHAR(128) NOT NULL COMMENT '接口任务标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    scheduled_date DATE NOT NULL COMMENT '预定日',
    zone_id VARCHAR(128) NOT NULL COMMENT '时区',
    proof_asset_id VARCHAR(128) NULL COMMENT '凭证资产标识',
    title_snapshot VARCHAR(160) NOT NULL COMMENT '标题快照',
    stars_snapshot INT UNSIGNED NOT NULL COMMENT '星星快照',
    completed_by JSON NOT NULL COMMENT '完成操作者快照',
    completed_at DATETIME NOT NULL COMMENT '完成时间',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '提交序列',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_completion_id (completion_id),
    UNIQUE KEY uk_task_member_date (task_id, member_id, scheduled_date),
    KEY idx_circle_member_time (circle_id, member_id, completed_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口任务完成审计';

CREATE TABLE IF NOT EXISTS kids_task_cancellation (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    cancellation_id VARCHAR(128) NOT NULL COMMENT '接口取消标识',
    completion_id VARCHAR(128) NOT NULL COMMENT '接口完成标识',
    reason_code VARCHAR(64) NOT NULL COMMENT '取消原因',
    cancelled_by JSON NOT NULL COMMENT '取消操作者快照',
    cancelled_at DATETIME NOT NULL COMMENT '取消时间',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '提交序列',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_cancellation_id (cancellation_id),
    UNIQUE KEY uk_completion_id (completion_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口任务取消审计';

CREATE TABLE IF NOT EXISTS kids_star_ledger (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    ledger_id VARCHAR(128) NOT NULL COMMENT '接口流水标识',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    source JSON NOT NULL COMMENT '来源快照',
    delta INT NOT NULL COMMENT '余额变化',
    reason VARCHAR(500) NULL COMMENT '原因',
    actor JSON NOT NULL COMMENT '操作者快照',
    reversal_of_ledger_id VARCHAR(128) NULL COMMENT '原流水标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '提交序列',
    created_at DATETIME NOT NULL COMMENT '创建时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_ledger_id (ledger_id),
    KEY idx_circle_member_time (circle_id, member_id, created_at),
    KEY idx_reversal (reversal_of_ledger_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口星星审计流水';

CREATE TABLE IF NOT EXISTS kids_star_balance (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键',
    circle_id VARCHAR(128) NOT NULL COMMENT '接口圈子标识',
    member_id VARCHAR(128) NOT NULL COMMENT '接口成员标识',
    balance INT NOT NULL DEFAULT 0 COMMENT '星星余额',
    version BIGINT UNSIGNED NOT NULL DEFAULT 1 COMMENT '版本号',
    source_commit_id VARCHAR(128) NOT NULL COMMENT '来源提交标识',
    source_commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '来源提交序列',
    updated_at DATETIME NOT NULL COMMENT '更新时间',
    PRIMARY KEY (id),
    UNIQUE KEY uk_circle_member (circle_id, member_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口成员星星余额';
