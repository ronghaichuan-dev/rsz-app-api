-- 修复历史部署遗漏的接口提交序列初始化行。
-- 该脚本可重复执行，已有序列不会被重置，避免影响既有同步游标。

INSERT INTO kids_sync_sequence (id, next_commit_sequence)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE next_commit_sequence = next_commit_sequence;
