-- 修复历史星星流水使用的 ledger:v1 标识，恢复冻结接口要求的 star-transaction:v1 命名空间。
-- 执行期间应停止旧版 kids 进程，避免旧代码并发写入旧命名空间；全部关联更新在同一事务内完成。

START TRANSACTION;

UPDATE kids_star_ledger
SET reversal_of_ledger_id = CONCAT('star-transaction:v1:', SUBSTRING(reversal_of_ledger_id, CHAR_LENGTH('ledger:v1:') + 1))
WHERE reversal_of_ledger_id REGEXP '^ledger:v1:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$';

UPDATE kids_reward_exchange
SET ledger_id = CONCAT('star-transaction:v1:', SUBSTRING(ledger_id, CHAR_LENGTH('ledger:v1:') + 1))
WHERE ledger_id REGEXP '^ledger:v1:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$';

UPDATE kids_sync_commit
SET change_payload = CAST(
    REPLACE(CAST(change_payload AS CHAR CHARACTER SET utf8mb4), 'ledger:v1:', 'star-transaction:v1:')
    AS JSON
)
WHERE CAST(change_payload AS CHAR CHARACTER SET utf8mb4) LIKE '%ledger:v1:%';

UPDATE kids_request_deduplication
SET response_body = CAST(
    REPLACE(CAST(response_body AS CHAR CHARACTER SET utf8mb4), 'ledger:v1:', 'star-transaction:v1:')
    AS JSON
)
WHERE CAST(response_body AS CHAR CHARACTER SET utf8mb4) LIKE '%ledger:v1:%';

UPDATE kids_star_ledger
SET ledger_id = CONCAT('star-transaction:v1:', SUBSTRING(ledger_id, CHAR_LENGTH('ledger:v1:') + 1))
WHERE ledger_id REGEXP '^ledger:v1:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$';

COMMIT;

-- 为成员流水的时间倒序分页提供覆盖排序的组合索引，避免读取全量记录后在应用层分页。
SET @kids_ledger_pagination_index_exists = (
    SELECT COUNT(*)
    FROM information_schema.statistics
    WHERE table_schema = DATABASE()
        AND table_name = 'kids_star_ledger'
        AND index_name = 'idx_circle_member_created_id'
);
SET @kids_ledger_pagination_index_sql = IF(
    @kids_ledger_pagination_index_exists = 0,
    'ALTER TABLE kids_star_ledger ADD KEY idx_circle_member_created_id (circle_id, member_id, created_at, id)',
    'SELECT 1'
);
PREPARE kids_migration_statement FROM @kids_ledger_pagination_index_sql;
EXECUTE kids_migration_statement;
DEALLOCATE PREPARE kids_migration_statement;
