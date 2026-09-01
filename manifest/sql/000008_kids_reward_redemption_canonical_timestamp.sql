-- 为奖励兑换的事务事实保留毫秒精度，并以同步提交时间修复历史审计投影。
-- 本迁移仅作用于 kids 微服务自己的数据库；执行前应完成 000006 时间精度迁移。

ALTER TABLE kids_reward_cooldown
    MODIFY COLUMN cooldown_until_at DATETIME(3) NULL COMMENT '冷却结束时间',
    MODIFY COLUMN last_redeemed_at DATETIME(3) NOT NULL COMMENT '最近兑换时间',
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间';

ALTER TABLE kids_reward_exchange
    MODIFY COLUMN cooldown_until_at_snapshot DATETIME(3) NULL COMMENT '冷却结束快照',
    MODIFY COLUMN exchanged_at DATETIME(3) NOT NULL COMMENT '兑换时间';

ALTER TABLE kids_notification_outbox
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间';

START TRANSACTION;

-- 历史兑换使用对应 sync commit 的持久化时间，恢复 Exchange、Ledger、Balance 与通知的同一审计时刻。
UPDATE kids_reward_exchange AS exchange_record
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.circle_id = exchange_record.circle_id
    AND sync_commit.commit_sequence = exchange_record.commit_sequence
SET exchange_record.exchanged_at = sync_commit.created_at;

UPDATE kids_star_ledger AS ledger
INNER JOIN kids_reward_exchange AS exchange_record
    ON exchange_record.ledger_id = ledger.ledger_id
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.circle_id = exchange_record.circle_id
    AND sync_commit.commit_sequence = exchange_record.commit_sequence
SET ledger.created_at = sync_commit.created_at;

UPDATE kids_star_balance AS balance
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.commit_id = balance.source_commit_id
    AND sync_commit.commit_sequence = balance.source_commit_sequence
INNER JOIN kids_reward_exchange AS exchange_record
    ON exchange_record.circle_id = balance.circle_id
    AND exchange_record.member_id = balance.member_id
    AND exchange_record.commit_sequence = balance.source_commit_sequence
SET balance.updated_at = sync_commit.created_at;

UPDATE kids_notification_outbox AS notification
INNER JOIN kids_reward_exchange AS exchange_record
    ON exchange_record.exchange_id = notification.exchange_id
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.circle_id = exchange_record.circle_id
    AND sync_commit.commit_sequence = exchange_record.commit_sequence
SET notification.created_at = sync_commit.created_at,
    notification.updated_at = sync_commit.created_at;

-- 冷却表仅保存每个奖励成员组合的最近一次兑换，因此只修复该组合的最大提交序列。
UPDATE kids_reward_cooldown AS cooldown
INNER JOIN (
    SELECT reward_id, member_id, MAX(commit_sequence) AS commit_sequence
    FROM kids_reward_exchange
    GROUP BY reward_id, member_id
) AS latest_exchange
    ON latest_exchange.reward_id = cooldown.reward_id
    AND latest_exchange.member_id = cooldown.member_id
INNER JOIN kids_reward_exchange AS exchange_record
    ON exchange_record.reward_id = latest_exchange.reward_id
    AND exchange_record.member_id = latest_exchange.member_id
    AND exchange_record.commit_sequence = latest_exchange.commit_sequence
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.circle_id = exchange_record.circle_id
    AND sync_commit.commit_sequence = exchange_record.commit_sequence
SET cooldown.last_redeemed_at = sync_commit.created_at,
    cooldown.updated_at = sync_commit.created_at;

UPDATE kids_mutation_receipt AS receipt
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.commit_id = receipt.commit_id
SET receipt.committed_at = sync_commit.created_at
WHERE receipt.operation_id = 'redeemReward';

-- 旧版兑换 commit 只保存了标识摘要；以不可变审计行重建完整同步投影，避免重启后遗漏兑换事实。
UPDATE kids_sync_commit AS sync_commit
INNER JOIN kids_reward_exchange AS exchange_record
    ON exchange_record.circle_id = sync_commit.circle_id
    AND exchange_record.commit_sequence = sync_commit.commit_sequence
INNER JOIN (
    SELECT reward_id, member_id, MAX(commit_sequence) AS commit_sequence
    FROM kids_reward_exchange
    GROUP BY reward_id, member_id
) AS latest_exchange
    ON latest_exchange.reward_id = exchange_record.reward_id
    AND latest_exchange.member_id = exchange_record.member_id
    AND latest_exchange.commit_sequence = exchange_record.commit_sequence
INNER JOIN kids_star_ledger AS ledger
    ON ledger.ledger_id = exchange_record.ledger_id
INNER JOIN kids_star_balance AS balance
    ON balance.circle_id = exchange_record.circle_id
    AND balance.member_id = exchange_record.member_id
    AND balance.source_commit_sequence = exchange_record.commit_sequence
INNER JOIN kids_reward_cooldown AS cooldown
    ON cooldown.reward_id = exchange_record.reward_id
    AND cooldown.member_id = exchange_record.member_id
INNER JOIN kids_notification_outbox AS notification
    ON notification.exchange_id = exchange_record.exchange_id
SET sync_commit.change_payload = JSON_OBJECT(
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
    'reward_cooldowns', JSON_ARRAY(JSON_OBJECT(
        'reward_id', cooldown.reward_id,
        'member_id', cooldown.member_id,
        'cooldown_until_ms', IF(cooldown.cooldown_until_at IS NULL, NULL, CAST(UNIX_TIMESTAMP(cooldown.cooldown_until_at) * 1000 AS UNSIGNED)),
        'permanently_unavailable', IF(cooldown.permanently_unavailable = 1, TRUE, FALSE),
        'last_redeemed_at_ms', CAST(UNIX_TIMESTAMP(cooldown.last_redeemed_at) * 1000 AS UNSIGNED),
        'version', cooldown.version
    )),
    'ledger_entries', JSON_ARRAY(JSON_OBJECT(
        'ledger_id', ledger.ledger_id,
        'circle_id', ledger.circle_id,
        'member_id', ledger.member_id,
        'source', JSON_EXTRACT(ledger.source, '$'),
        'delta', ledger.delta,
        'reason', ledger.reason,
        'actor', JSON_EXTRACT(ledger.actor, '$'),
        'reversal_of_ledger_id', ledger.reversal_of_ledger_id,
        'created_at_ms', CAST(UNIX_TIMESTAMP(ledger.created_at) * 1000 AS UNSIGNED),
        'commit_sequence', ledger.commit_sequence
    )),
    'star_balances', JSON_ARRAY(JSON_OBJECT(
        'circle_id', balance.circle_id,
        'member_id', balance.member_id,
        'balance', balance.balance,
        'version', balance.version,
        'source_commit_id', balance.source_commit_id,
        'source_commit_sequence', balance.source_commit_sequence,
        'updated_at_ms', CAST(UNIX_TIMESTAMP(balance.updated_at) * 1000 AS UNSIGNED)
    )),
    'exchanges', JSON_ARRAY(JSON_OBJECT(
        'exchange_id', exchange_record.exchange_id,
        'circle_id', exchange_record.circle_id,
        'member_id', exchange_record.member_id,
        'member_name_snapshot', exchange_record.member_name_snapshot,
        'member_avatar_snapshot', JSON_EXTRACT(exchange_record.member_avatar_snapshot, '$'),
        'reward_id', exchange_record.reward_id,
        'reward_title_snapshot', exchange_record.reward_title_snapshot,
        'reward_visual_snapshot', JSON_EXTRACT(exchange_record.reward_visual_snapshot, '$'),
        'stars_deducted_snapshot', exchange_record.stars_deducted_snapshot,
        'reward_repeat_rule_snapshot', exchange_record.reward_repeat_rule_snapshot,
        'reward_cooldown_days_snapshot', exchange_record.reward_cooldown_days_snapshot,
        'cooldown_until_ms_snapshot', IF(exchange_record.cooldown_until_at_snapshot IS NULL, NULL, CAST(UNIX_TIMESTAMP(exchange_record.cooldown_until_at_snapshot) * 1000 AS UNSIGNED)),
        'permanently_unavailable_snapshot', IF(exchange_record.permanently_unavailable_snapshot = 1, TRUE, FALSE),
        'ledger_id', exchange_record.ledger_id,
        'exchanged_at_ms', CAST(UNIX_TIMESTAMP(exchange_record.exchanged_at) * 1000 AS UNSIGNED),
        'commit_sequence', exchange_record.commit_sequence
    )),
    'notification_outbox_entries', JSON_ARRAY(JSON_OBJECT(
        'notification_id', notification.notification_id,
        'event_type', notification.event_type,
        'exchange_id', notification.exchange_id,
        'status', notification.status,
        'attempt_count', notification.attempt_count,
        'next_attempt_at_ms', IF(notification.next_attempt_at IS NULL, NULL, CAST(UNIX_TIMESTAMP(notification.next_attempt_at) * 1000 AS UNSIGNED)),
        'created_at_ms', CAST(UNIX_TIMESTAMP(notification.created_at) * 1000 AS UNSIGNED),
        'updated_at_ms', CAST(UNIX_TIMESTAMP(notification.updated_at) * 1000 AS UNSIGNED),
        'version', notification.version
    )),
    'tombstones', JSON_ARRAY()
);

-- 修复已保存的幂等成功快照，使旧 key 重放时仍得到同一 canonical 时间而不会重复扣星。
UPDATE kids_request_deduplication AS deduplication
INNER JOIN kids_mutation_receipt AS receipt
    ON receipt.commit_id = JSON_UNQUOTE(JSON_EXTRACT(deduplication.response_body, '$.receipt.commit_id'))
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.commit_id = receipt.commit_id
SET deduplication.response_body = JSON_SET(
    deduplication.response_body,
    '$.receipt.committed_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.exchange.exchanged_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.ledger_entry.created_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.balance.updated_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.cooldown.last_redeemed_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.notification_outbox.created_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED),
    '$.notification_outbox.updated_at_ms', CAST(UNIX_TIMESTAMP(sync_commit.created_at) * 1000 AS UNSIGNED)
)
WHERE deduplication.operation_id = 'redeemReward';

COMMIT;
