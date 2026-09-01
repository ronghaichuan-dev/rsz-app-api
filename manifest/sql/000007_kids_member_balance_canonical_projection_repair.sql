-- 修复历史调整响应、余额行和同步提交投影不一致的余额事实。
-- 本迁移只追加新的修复提交并提升受影响余额版本，绝不覆盖既有版本或既有提交。
-- 执行期间应停止旧版 kids 进程，避免与在线提交序列并发分配。

START TRANSACTION;

INSERT INTO kids_sync_sequence (id, next_commit_sequence)
VALUES (1, 1)
ON DUPLICATE KEY UPDATE next_commit_sequence = next_commit_sequence;

SELECT next_commit_sequence
INTO @kids_balance_canonical_repair_next_sequence
FROM kids_sync_sequence
WHERE id = 1
FOR UPDATE;

CREATE TEMPORARY TABLE kids_member_balance_canonical_repair (
    circle_id VARCHAR(128) NOT NULL COMMENT '圈子标识',
    member_id VARCHAR(128) NOT NULL COMMENT '成员标识',
    balance INT NOT NULL COMMENT '修复前后保持的星星余额',
    previous_version BIGINT UNSIGNED NOT NULL COMMENT '修复前余额版本',
    commit_id VARCHAR(128) NOT NULL COMMENT '修复提交标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '修复提交序列',
    updated_at DATETIME(3) NOT NULL COMMENT '规范投影更新时间',
    PRIMARY KEY (circle_id, member_id)
) ENGINE=InnoDB COMMENT='成员余额规范投影修复暂存表';

CREATE TEMPORARY TABLE kids_member_balance_canonical_repair_commit (
    circle_id VARCHAR(128) NOT NULL COMMENT '圈子标识',
    commit_id VARCHAR(128) NOT NULL COMMENT '修复提交标识',
    commit_sequence BIGINT UNSIGNED NOT NULL COMMENT '修复提交序列',
    PRIMARY KEY (circle_id),
    UNIQUE KEY uk_commit_id (commit_id),
    UNIQUE KEY uk_commit_sequence (commit_sequence)
) ENGINE=InnoDB COMMENT='成员余额规范投影修复提交暂存表';

INSERT INTO kids_member_balance_canonical_repair_commit (circle_id, commit_id, commit_sequence)
SELECT
    affected.circle_id,
    CONCAT('commit:v1:', UUID()),
    (@kids_balance_canonical_repair_next_sequence := @kids_balance_canonical_repair_next_sequence + 1) - 1
FROM (
    SELECT DISTINCT balance.circle_id
    FROM kids_star_balance AS balance
    LEFT JOIN kids_sync_commit AS source_commit
        ON source_commit.commit_id = balance.source_commit_id
        AND source_commit.commit_sequence = balance.source_commit_sequence
    WHERE source_commit.id IS NULL
        OR JSON_CONTAINS(
            source_commit.change_payload,
            JSON_OBJECT(
                'circle_id', balance.circle_id,
                'member_id', balance.member_id,
                'balance', balance.balance,
                'version', balance.version,
                'source_commit_id', balance.source_commit_id,
                'source_commit_sequence', balance.source_commit_sequence,
                'updated_at_ms', CAST(UNIX_TIMESTAMP(balance.updated_at) * 1000 AS UNSIGNED)
            ),
            '$.star_balances'
        ) = 0
) AS affected
ORDER BY affected.circle_id;

INSERT INTO kids_member_balance_canonical_repair (
    circle_id,
    member_id,
    balance,
    previous_version,
    commit_id,
    commit_sequence,
    updated_at
)
SELECT
    balance.circle_id,
    balance.member_id,
    balance.balance,
    balance.version,
    repair.commit_id,
    repair.commit_sequence,
    UTC_TIMESTAMP(3)
FROM kids_star_balance AS balance
INNER JOIN kids_member_balance_canonical_repair_commit AS repair
    ON repair.circle_id = balance.circle_id
LEFT JOIN kids_sync_commit AS source_commit
    ON source_commit.commit_id = balance.source_commit_id
    AND source_commit.commit_sequence = balance.source_commit_sequence
WHERE source_commit.id IS NULL
    OR JSON_CONTAINS(
        source_commit.change_payload,
        JSON_OBJECT(
            'circle_id', balance.circle_id,
            'member_id', balance.member_id,
            'balance', balance.balance,
            'version', balance.version,
            'source_commit_id', balance.source_commit_id,
            'source_commit_sequence', balance.source_commit_sequence,
            'updated_at_ms', CAST(UNIX_TIMESTAMP(balance.updated_at) * 1000 AS UNSIGNED)
        ),
        '$.star_balances'
    ) = 0;

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
        'star_balances', (
            SELECT JSON_ARRAYAGG(
                JSON_OBJECT(
                    'circle_id', balance.circle_id,
                    'member_id', balance.member_id,
                    'balance', balance.balance,
                    'version', balance.previous_version + 1,
                    'source_commit_id', balance.commit_id,
                    'source_commit_sequence', balance.commit_sequence,
                    'updated_at_ms', CAST(UNIX_TIMESTAMP(balance.updated_at) * 1000 AS UNSIGNED)
                )
            )
            FROM kids_member_balance_canonical_repair AS balance
            WHERE balance.circle_id = repair.circle_id
        ),
        'exchanges', JSON_ARRAY(),
        'notification_outbox_entries', JSON_ARRAY(),
        'tombstones', JSON_ARRAY()
    ),
    UTC_TIMESTAMP(3)
FROM kids_member_balance_canonical_repair_commit AS repair;

UPDATE kids_star_balance AS balance
INNER JOIN kids_member_balance_canonical_repair AS repair
    ON repair.circle_id = balance.circle_id
    AND repair.member_id = balance.member_id
SET
    balance.version = repair.previous_version + 1,
    balance.source_commit_id = repair.commit_id,
    balance.source_commit_sequence = repair.commit_sequence,
    balance.updated_at = repair.updated_at;

UPDATE kids_sync_sequence
SET next_commit_sequence = @kids_balance_canonical_repair_next_sequence
WHERE id = 1;

DROP TEMPORARY TABLE kids_member_balance_canonical_repair;
DROP TEMPORARY TABLE kids_member_balance_canonical_repair_commit;

COMMIT;
