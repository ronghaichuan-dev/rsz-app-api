-- Clearwave 接口提交序列。
-- 单行计数器在同一事务内锁定，确保 commit_sequence 全局单调且跨进程持久化。

CREATE TABLE IF NOT EXISTS kids_sync_sequence (
    id TINYINT UNSIGNED NOT NULL COMMENT '固定序列标识',
    next_commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '下一个提交序列',
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='kids接口提交序列';

INSERT INTO kids_sync_sequence (id, next_commit_sequence)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE next_commit_sequence = next_commit_sequence;
