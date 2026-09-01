-- 为任务完成的 occurrence 与 completion 保留毫秒精度，并修复已受影响的 canonical bundle。
-- 本迁移不创建或删除 Completion、LedgerEntry、余额、commit、receipt 或幂等键；仅以既有 receipt
-- 的 committed_at 统一同一 completeTask 提交已持久化的时间投影及其同步、幂等载荷。
-- 执行期间必须停止旧版 kids 进程，避免在线请求读取到修复中的历史 bundle。

ALTER TABLE kids_task_occurrence
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间',
    MODIFY COLUMN updated_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3) COMMENT '更新时间';

ALTER TABLE kids_task_completion
    MODIFY COLUMN completed_at DATETIME(3) NOT NULL COMMENT '完成时间',
    MODIFY COLUMN created_at DATETIME(3) NOT NULL DEFAULT CURRENT_TIMESTAMP(3) COMMENT '创建时间';

START TRANSACTION;

-- 仅修复尚未被后续取消或余额写入覆盖的任务完成，receipt 的时间是该提交唯一的规范时刻。
UPDATE kids_task_completion AS completion
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.circle_id = completion.circle_id
    AND sync_commit.commit_sequence = completion.commit_sequence
INNER JOIN kids_mutation_receipt AS receipt
    ON receipt.commit_id = sync_commit.commit_id
    AND receipt.operation_id = 'completeTask'
INNER JOIN kids_task_occurrence AS occurrence
    ON occurrence.circle_id = completion.circle_id
    AND occurrence.completion_id = completion.completion_id
    AND occurrence.state = 'completed'
INNER JOIN kids_star_ledger AS ledger
    ON ledger.circle_id = completion.circle_id
    AND ledger.member_id = completion.member_id
    AND ledger.commit_sequence = completion.commit_sequence
    AND ledger.delta = completion.stars_snapshot
INNER JOIN kids_star_balance AS balance
    ON balance.circle_id = completion.circle_id
    AND balance.member_id = completion.member_id
    AND balance.source_commit_id = sync_commit.commit_id
    AND balance.source_commit_sequence = completion.commit_sequence
SET completion.completed_at = receipt.committed_at,
    occurrence.updated_at = receipt.committed_at,
    ledger.created_at = receipt.committed_at,
    balance.updated_at = receipt.committed_at,
    sync_commit.created_at = receipt.committed_at
WHERE completion.completed_at <> receipt.committed_at
    OR occurrence.updated_at <> receipt.committed_at
    OR ledger.created_at <> receipt.committed_at
    OR balance.updated_at <> receipt.committed_at
    OR sync_commit.created_at <> receipt.committed_at;

-- 已修复的基础事实重新生成完整同步 payload，使增量同步与命令响应读取同一 canonical bundle。
UPDATE kids_sync_commit AS sync_commit
INNER JOIN kids_mutation_receipt AS receipt
    ON receipt.commit_id = sync_commit.commit_id
    AND receipt.operation_id = 'completeTask'
INNER JOIN kids_task_completion AS completion
    ON completion.circle_id = sync_commit.circle_id
    AND completion.commit_sequence = sync_commit.commit_sequence
INNER JOIN kids_task_occurrence AS occurrence
    ON occurrence.circle_id = completion.circle_id
    AND occurrence.completion_id = completion.completion_id
    AND occurrence.state = 'completed'
INNER JOIN kids_star_ledger AS ledger
    ON ledger.circle_id = completion.circle_id
    AND ledger.member_id = completion.member_id
    AND ledger.commit_sequence = completion.commit_sequence
    AND ledger.delta = completion.stars_snapshot
INNER JOIN kids_star_balance AS balance
    ON balance.circle_id = completion.circle_id
    AND balance.member_id = completion.member_id
    AND balance.source_commit_id = sync_commit.commit_id
    AND balance.source_commit_sequence = completion.commit_sequence
SET sync_commit.change_payload = JSON_OBJECT(
    'circles', JSON_ARRAY(),
    'memberships', JSON_ARRAY(),
    'administrators', JSON_ARRAY(),
    'members', JSON_ARRAY(),
    'circle_selections', JSON_ARRAY(),
    'outgoing_invites', JSON_ARRAY(),
    'task_tags', JSON_ARRAY(),
    'tasks', JSON_ARRAY(),
    'task_occurrences', JSON_ARRAY(JSON_OBJECT(
        'circle_id', occurrence.circle_id,
        'task_id', occurrence.task_id,
        'member_id', occurrence.member_id,
        'scheduled_date', DATE_FORMAT(occurrence.scheduled_date, '%Y-%m-%d'),
        'zone_id', occurrence.zone_id,
        'definition_revision', occurrence.definition_revision,
        'title_snapshot', occurrence.title_snapshot,
        'notes_snapshot', occurrence.notes_snapshot,
        'emoji_snapshot', occurrence.emoji_snapshot,
        'stars_snapshot', occurrence.stars_snapshot,
        'photo_required_snapshot', IF(occurrence.photo_required_snapshot = 1, TRUE, FALSE),
        'task_tag_id_snapshot', occurrence.task_tag_id_snapshot,
        'state', occurrence.state,
        'completion_id', occurrence.completion_id,
        'version', occurrence.version,
        'created_at_ms', CAST(UNIX_TIMESTAMP(occurrence.created_at) * 1000 AS UNSIGNED),
        'updated_at_ms', CAST(UNIX_TIMESTAMP(occurrence.updated_at) * 1000 AS UNSIGNED)
    )),
    'task_occurrence_tombstones', JSON_ARRAY(),
    'task_completions', JSON_ARRAY(JSON_OBJECT(
        'completion_id', completion.completion_id,
        'circle_id', completion.circle_id,
        'task_id', completion.task_id,
        'member_id', completion.member_id,
        'scheduled_date', DATE_FORMAT(completion.scheduled_date, '%Y-%m-%d'),
        'zone_id', completion.zone_id,
        'completed_at_ms', CAST(UNIX_TIMESTAMP(completion.completed_at) * 1000 AS UNSIGNED),
        'proof_asset_id', completion.proof_asset_id,
        'title_snapshot', completion.title_snapshot,
        'stars_snapshot', completion.stars_snapshot,
        'completed_by', JSON_EXTRACT(completion.completed_by, '$'),
        'version', completion.version,
        'commit_sequence', completion.commit_sequence
    )),
    'task_cancellations', JSON_ARRAY(),
    'rewards', JSON_ARRAY(),
    'reward_cooldowns', JSON_ARRAY(),
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
    'exchanges', JSON_ARRAY(),
    'notification_outbox_entries', JSON_ARRAY(),
    'tombstones', JSON_ARRAY()
);

-- 使用同一基础事实回写保存的首次响应，幂等重放只改 result_kind，不改变任何 ID 或 cursor。
UPDATE kids_request_deduplication AS deduplication
INNER JOIN kids_mutation_receipt AS receipt
    ON receipt.commit_id = JSON_UNQUOTE(JSON_EXTRACT(deduplication.response_body, '$.receipt.commit_id'))
    AND receipt.operation_id = 'completeTask'
INNER JOIN kids_sync_commit AS sync_commit
    ON sync_commit.commit_id = receipt.commit_id
INNER JOIN kids_task_completion AS completion
    ON completion.circle_id = sync_commit.circle_id
    AND completion.commit_sequence = sync_commit.commit_sequence
INNER JOIN kids_task_occurrence AS occurrence
    ON occurrence.circle_id = completion.circle_id
    AND occurrence.completion_id = completion.completion_id
    AND occurrence.state = 'completed'
INNER JOIN kids_star_ledger AS ledger
    ON ledger.circle_id = completion.circle_id
    AND ledger.member_id = completion.member_id
    AND ledger.commit_sequence = completion.commit_sequence
    AND ledger.delta = completion.stars_snapshot
INNER JOIN kids_star_balance AS balance
    ON balance.circle_id = completion.circle_id
    AND balance.member_id = completion.member_id
    AND balance.source_commit_id = sync_commit.commit_id
    AND balance.source_commit_sequence = completion.commit_sequence
SET deduplication.response_body = JSON_SET(
    deduplication.response_body,
    '$.receipt.committed_at_ms', CAST(UNIX_TIMESTAMP(receipt.committed_at) * 1000 AS UNSIGNED),
    '$.occurrence.updated_at_ms', CAST(UNIX_TIMESTAMP(occurrence.updated_at) * 1000 AS UNSIGNED),
    '$.completion.completed_at_ms', CAST(UNIX_TIMESTAMP(completion.completed_at) * 1000 AS UNSIGNED),
    '$.ledger_entry.created_at_ms', CAST(UNIX_TIMESTAMP(ledger.created_at) * 1000 AS UNSIGNED),
    '$.balance.updated_at_ms', CAST(UNIX_TIMESTAMP(balance.updated_at) * 1000 AS UNSIGNED)
)
WHERE deduplication.operation_id = 'completeTask';

COMMIT;
