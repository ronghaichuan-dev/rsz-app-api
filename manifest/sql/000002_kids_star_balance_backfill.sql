-- 为既有有效成员补齐规范星星余额事实行。
-- 本迁移应在 kids_task_ledger_migration.sql 与 kids_circle_domain_migration.sql 执行后运行。
-- 使用唯一键保证可重复执行时不重复创建余额记录，也不覆盖既有余额或版本。

INSERT INTO kids_star_balance (
    circle_id,
    member_id,
    balance,
    version,
    source_commit_id,
    source_commit_sequence,
    updated_at
)
SELECT
    member.circle_id,
    member.member_id,
    0,
    1,
    'migration_star_balance_backfill',
    0,
    UTC_TIMESTAMP(3)
FROM kids_member AS member
LEFT JOIN kids_star_balance AS balance
    ON balance.circle_id = member.circle_id
    AND balance.member_id = member.member_id
WHERE member.status = 'active'
    AND balance.id IS NULL;
