-- 修复历史余额回填产生的非规范来源提交，并补齐遗漏的有效成员余额快照。
-- 本迁移为追加式修复：每个受影响圈子创建一个新的规范提交，余额读取和后续 expected_balance_version 均可使用该快照。
-- 执行期间应停止旧版 kids 进程，避免与在线提交序列并发分配。

START TRANSACTION;

INSERT INTO kids_sync_sequence (id, next_commit_sequence)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE next_commit_sequence = next_commit_sequence;

SELECT next_commit_sequence
INTO @kids_balance_repair_next_sequence
FROM kids_sync_sequence
WHERE id = 1
FOR UPDATE;

CREATE TEMPORARY TABLE kids_member_balance_snapshot_repair (
    circle_id VARCHAR(128) NOT NULL COMMENT '圈子标识',
    commit_id VARCHAR(128) NOT NULL COMMENT '修复提交标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '修复提交序列',
    PRIMARY KEY (circle_id),
    UNIQUE KEY uk_commit_id (commit_id),
    UNIQUE KEY uk_commit_sequence (commit_sequence)
) ENGINE=InnoDB COMMENT='余额快照修复暂存表';

INSERT INTO kids_member_balance_snapshot_repair (circle_id, commit_id, commit_sequence)
SELECT
    affected.circle_id,
    CONCAT('commit:v1:', UUID()),
    (@kids_balance_repair_next_sequence := @kids_balance_repair_next_sequence + 1) - 1
FROM (
    SELECT DISTINCT member.circle_id
    FROM kids_member AS member
    LEFT JOIN kids_star_balance AS balance
        ON balance.circle_id = member.circle_id
        AND balance.member_id = member.member_id
    LEFT JOIN kids_sync_commit AS source_commit
        ON source_commit.commit_id = balance.source_commit_id
        AND source_commit.commit_sequence = balance.source_commit_sequence
    WHERE member.status = 'active'
        AND (
            balance.id IS NULL
            OR balance.version < 1
            OR balance.source_commit_sequence < 1
            OR balance.source_commit_id NOT REGEXP '^commit:v1:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
            OR source_commit.id IS NULL
        )
) AS affected
ORDER BY affected.circle_id;

INSERT INTO kids_sync_commit (commit_id, circle_id, commit_sequence, change_payload, created_at)
SELECT
    repair.commit_id,
    repair.circle_id,
    repair.commit_sequence,
    JSON_OBJECT(
        'circles', JSON_ARRAY(),
        'memberships', JSON_ARRAY(),
        'administrators', JSON_ARRAY(),
        'members', JSON_ARRAY(),
        'circle_selections', JSON_ARRAY(),
        'outgoing_invites', JSON_ARRAY(),
        'task_tags', JSON_ARRAY(),
        'tasks', JSON_ARRAY(),
        'task_occurrences', JSON_ARRAY(),
        'task_occurrence_tombstones', JSON_ARRAY(),
        'task_completions', JSON_ARRAY(),
        'task_cancellations', JSON_ARRAY(),
        'rewards', JSON_ARRAY(),
        'reward_cooldowns', JSON_ARRAY(),
        'ledger_entries', JSON_ARRAY(),
        'star_balances', JSON_ARRAY(),
        'exchanges', JSON_ARRAY(),
        'notification_outbox_entries', JSON_ARRAY(),
        'tombstones', JSON_ARRAY()
    ),
    UTC_TIMESTAMP(3)
FROM kids_member_balance_snapshot_repair AS repair;

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
    repair.commit_id,
    repair.commit_sequence,
    UTC_TIMESTAMP(3)
FROM kids_member AS member
INNER JOIN kids_member_balance_snapshot_repair AS repair
    ON repair.circle_id = member.circle_id
LEFT JOIN kids_star_balance AS balance
    ON balance.circle_id = member.circle_id
    AND balance.member_id = member.member_id
WHERE member.status = 'active'
    AND balance.id IS NULL;

UPDATE kids_star_balance AS balance
INNER JOIN kids_member AS member
    ON member.circle_id = balance.circle_id
    AND member.member_id = balance.member_id
    AND member.status = 'active'
INNER JOIN kids_member_balance_snapshot_repair AS repair
    ON repair.circle_id = balance.circle_id
LEFT JOIN kids_sync_commit AS source_commit
    ON source_commit.commit_id = balance.source_commit_id
    AND source_commit.commit_sequence = balance.source_commit_sequence
SET
    balance.balance = GREATEST(balance.balance, 0),
    balance.version = GREATEST(balance.version, 1),
    balance.source_commit_id = repair.commit_id,
    balance.source_commit_sequence = repair.commit_sequence,
    balance.updated_at = UTC_TIMESTAMP(3)
WHERE balance.version < 1
    OR balance.source_commit_sequence < 1
    OR balance.source_commit_id NOT REGEXP '^commit:v1:[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$'
    OR source_commit.id IS NULL;

UPDATE kids_sync_sequence
SET next_commit_sequence = @kids_balance_repair_next_sequence
WHERE id = 1;

DROP TEMPORARY TABLE kids_member_balance_snapshot_repair;

COMMIT;
