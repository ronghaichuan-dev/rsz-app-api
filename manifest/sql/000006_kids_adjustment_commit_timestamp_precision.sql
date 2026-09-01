-- 为星星调整的 ledger、余额、commit 与 receipt 保留毫秒精度，确保重启后的审计投影可与 mutation bundle 对齐。
-- 本迁移只提高时间精度，不修改既有业务事实或提交序列。

ALTER TABLE kids_star_ledger
    MODIFY COLUMN created_at DATETIME(3) NOT NULL COMMENT '创建时间';

ALTER TABLE kids_star_balance
    MODIFY COLUMN updated_at DATETIME(3) NOT NULL COMMENT '更新时间';

ALTER TABLE kids_sync_commit
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间';

ALTER TABLE kids_mutation_receipt
    MODIFY COLUMN committed_at DATETIME(3) NOT NULL COMMENT '提交时间';
