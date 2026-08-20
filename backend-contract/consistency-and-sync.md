# Consistency、Sync 与恢复

本合同保留 Room 作为客户端 UI 的唯一生产投影，但 Online canonical authority 位于服务端。客户端不能把“本地写入成功”解释为“远端已提交”，也不能把跨域命令拆成多个接口后自行拼事实。

## 一致性类别

| 类别 | 服务端/客户端语义 |
|---|---|
| `LOCAL_ONLY` | 事实只属于设备；Offline 调用自建后端的计数必须为 0 |
| `EXTERNAL_SDK` | Google/Android/Firebase/Ads 等外部 owner；自建 Backend/Demo 不冒充其结果 |
| `READ_SYNC` | 只读 canonical snapshot/cursor；客户端在校验完整后投影到 Room |
| `OPTIMISTIC_MUTATION` | 普通可变 aggregate；Room 与 outbox 可同事务形成 durable pending，服务端以 expected version 提交 |
| `REMOTE_ATOMIC_COMMAND` | 权限、余额或审计敏感命令；先落 durable pending command，只有服务端完整 commit bundle 投影成功后才产生 canonical 本地审计事实 |
| `NOT_REQUIRED` | capability 已由另一个唯一 operation 或明确本地/外部边界满足，不增加 API |

OpenAPI operation 只使用后三种远端类别中的 `READ_SYNC`、`OPTIMISTIC_MUTATION`、`REMOTE_ATOMIC_COMMAND`；完整 capability 归属见 coverage JSON。

## Mutation 与 outbox 状态机

未来客户端接线必须至少区分以下持久状态，命名可按实现调整，语义不可合并：

1. `LocalDraft`：尚未形成用户提交，不进入远端 outbox。RewardDraft 明确保留在此状态。
2. `PendingCommand`：body、path 参数、`X-Request-Id`、`Idempotency-Key`、expected version 与目标环境已原子持久化；允许进程重启后重试。
3. `RemoteCommittedProjectionPending`：服务端已返回 receipt/commit bundle，但本地原子投影尚未成功。必须保留同一 key、request 与完整响应或 receipt 以恢复。
4. `RemoteCommitted`：完整 bundle 已按一次本地事务投影，相关 change cursor 同事务推进。
5. `Rejected`：稳定非重试错误已记录为受控结果；不能伪装成功。

`PendingCommand` 不能让 UI 产生 canonical Completion、Ledger、Exchange 等审计行。普通 CRUD 可以显示明确 pending overlay，但不能覆盖 Room 中最后一次 `RemoteCommitted` 版本。

## 幂等与版本

- 所有 mutation 必须 durable 保存并复用同一个 `Idempotency-Key`；同一用户意图不得在 retry、进程恢复或响应丢失后换 key。
- 服务端按 OpenAPI 的 `x-clearwave-idempotency-policy` 执行：lookup scope 为服务端解析的 principal、当前 membership actor（无则 null）与 key；route fingerprint 包含 method、operationId、按名称解码后的 path 参数；body 使用 RFC 8785 JCS UTF-8 后做 SHA-256。
- 同 scope/key/fingerprint 返回第一次提交的同一 canonical response；若 operation 返回 `MutationReceipt`，只把 receipt 的 `result_kind` 标记为 `idempotent_replay`，其余 response 必须与首次结果一致；不返回 receipt 的 session/upload mutation 必须逐字段返回首次结果。不同 route 或 payload fingerprint 稳定返回 `IDEMPOTENCY_CONFLICT`，不执行第二次写入。
- 幂等记录必须跨服务进程重启与响应丢失持久化。含 token、invite secret 或 upload target 的响应需加密存储，不得因安全清理而把重放变成第二次提交。
- 可变 aggregate 使用 operation 冻结的 `expected_version`、`expected_*_version`；同一 upsert 创建自身 aggregate 时仅在 schema 明确允许时使用 null，update/delete 使用当前正整数。`createInvite` 虽创建 Invite，仍必须提交既有 target 的非空 version。冲突返回 `VERSION_CONFLICT` 和 `current_version`，客户端先 pull/rebase，不盲目覆盖。
- `upsertMember` 只更新既有 Member，`expected_version` 必须为当前正整数；新增 Member 只能调用 `createCircleMember`，由服务端原子执行额度、权限与初始关系检查。
- v1 不使用含糊的客户端时间或 last-write-wins。服务端时间、版本、余额、价格、stars delta、actor 与审计时间一律由 canonical state 推导。

## 不可拆原子命令

### Task completion

`completeTask` 在一个服务端事务内重新验证 session、membership、assignment、occurrence、proof asset 与 `expected_occurrence_version`，然后一次提交：

`Occurrence(completed) + immutable TaskCompletion + positive LedgerEntry + StarBalance + receipt + change_cursor`

`stars_snapshot`、delta、actor 与 `completed_at_ms` 从 Task/occurrence/token/server clock 推导。请求不得提供这些值。

### Task cancellation

`cancelTaskCompletion` 只接受 `cancellation_id`、`expected_completion_version` 与受控 reason。`TaskCompletion.version` 是客户端读到的 canonical 完成事实版本；完成 payload 本身仍 append-only，取消不会递增或覆盖它。服务端需检查尚无 cancellation，再一次提交：

`Occurrence(cancelled) + original TaskCompletion + TaskCancellation + original positive LedgerEntry + reversal LedgerEntry + StarBalance + receipt + change_cursor`

已取消返回 `COMPLETION_ALREADY_CANCELLED`；同一幂等键重放返回第一次 bundle。取消不删除原完成或原正向流水，也不接受客户端 stars、operator、时间或 proof。

### Reward redemption

Reward definition 的 `cooldown_days` 与当前领域规则同构：只有 `repeat_rule=custom` 接受正整数天数，`none/daily/weekly/monthly` 必须为 null；服务端不得静默忽略交叉组合。

`redeemReward` 在一个事务内重验 membership capability、Member/Reward 状态、assignment、canonical price、balance、cooldown 与 `expected_reward_version`，一次提交：

`negative LedgerEntry + StarBalance + RewardCooldown + immutable ExchangeRecord + NotificationOutboxEntry + receipt + change_cursor`

notification 只在数据库 commit 后由 outbox 投递；通知失败不能回滚兑换，通知重试不能再次扣星或新增 Exchange。

`NotificationOutboxEntry` 是可同步的 versioned 状态。兑换事务创建 `version=1`；worker 每次改变 status/attempt/next attempt 时提升 version 与 `updated_at_ms`，并把该变化作为独立完整 server commit 发布。客户端看到投递失败只能更新展示/重试状态，绝不能重放兑换命令来驱动通知。

### 其他敏感命令

Invite redeem、onboarding、member/admin 删除或离开、人工星星调整、Play purchase verification 等 operation 的 `x-atomic-commit-bundle` 是最小事务边界。没有 public append-ledger 或 create/update/delete Exchange API。

## Task future occurrence 边界

Task Save/删除请求必须携带 `future_effective_from_date` 与 IANA `zone_id`：

- 只替换该日及之后、仍为 `pending` 的 occurrence，并提升 Task `series_revision`/version。
- 该日前 occurrence 与所有 `completed`/`cancelled` occurrence 永不改写。
- response 的 `AffectedOccurrenceBoundary` 返回实际替换的 pending 数与保留的 resolved 数。
- 重复规则按请求 `zone_id` 解释；DST 变化不能把 calendar date 转成另一日期。`zone_id` 不属于 occurrence 身份。
- `weekdays` 只在 `frequency=custom_weekdays` 时为非空集合；`none/daily/weekly/monthly/yearly` 必须传空数组，固定规则从 Task 起始日推导日历锚点。
- `frequency=none` 只能使用 `end_rule.kind=never`；`on_date` 不得早于首个 occurrence。
- `after_occurrences.occurrence_count` 统计成员展开前的系列日历日期数；多个 assigned Member 不放大该计数。
- `time_limit` reminder 必须存在 `time_limit_minute_of_day` 且自身 `local_minute_of_day=null`；`start_date` reminder 必须携带 0..1439 的本地分钟。
- Task 顶层、repeat rule 与非空 reminder config 的 `zone_id` 必须完全相同；该唯一 zone 同时解释 future boundary、重复日期与本地提醒。
- occurrence 的永久身份是 `(task_id, member_id, scheduled_date)`；`circle_id` 是授权/分区上下文，不进入复合 key。

## ID 与 legacy mapping

- CircleId、MemberId、AdminId、TaskId、TaskTagId、RewardId、ExchangeId 等客户端稳定 ID 按 OpenAPI pattern 原样接受，只作为 opaque identifier，不解析权限。
- AdminId 永久接受 `admin:v1:<uuid>` 与历史 `admin:v1:migrated:circle:v1:<uuid>`；服务端禁止重写或折叠两种命名空间。
- Backend `account_id`、`session_id`、`membership_id` 与本地 IdentityScope 永久分离；scope stableId 不上 wire。
- 新 Completion 使用客户端生成的 `completion:v1:<uuid>`。当前 Room v8 既有 Completion 没有独立 ID，只以 `(task_id, member_id, scheduled_date)` 标识。
- legacy 补 ID 必须使用 RFC 9562 UUIDv5/SHA-1，namespace 为 `2f7de0ac-72f2-5e64-8bd2-8dadc3fb92f6`，name 是以下 UTF-8 字节（换行是单个 LF）：

```text
clearwave/backend-contract/v1/legacy-task-completion
{circle_id}
{task_id}
{member_id}
{scheduled_date}
```

其中日期严格为零填充 `yyyy-MM-dd`，结果加 `completion:v1:` 前缀。映射必须在任何远端 command 入队前持久化且永不重新分配。该映射只补 identifier；v1 没有“信任并导入旧本地审计”的接口，历史本地 Completion/Ledger 不会因此自动成为服务端 canonical 事实。

## Bootstrap、delta 与分页

- `getCircleBootstrap` 在固定 `snapshot_cursor` 下分页。客户端先把所有页放入 staging，只有 `has_more=false` 且 schema/引用完整时才一次事务替换该 Circle bootstrap；中途失败不暴露半个 snapshot。
- `pullCircleBootstrapDelta` 输入最后一次已投影的 `change_cursor`，输出有序 `SyncCommit[]`。每个 commit 含 `commit_id`、单调 `commit_sequence`、服务端时间、全域 changes 与其 change cursor。全域 changes 必须包括无 secret 的 outgoing invite 状态、versioned notification outbox 状态以及其余 schema 声明的 collection。
- 服务端不得把一个 `SyncCommit` 拆到不同页；如果下一个 commit 超过当前 page budget，应降低本页 item 数并完整返回该 commit，而不是截断 changes。
- 客户端按 commit 顺序投影；mutable entity 只接受更高 version upsert/tombstone，重复同 version 同 payload为 no-op。
- 普通 `EntityTombstone` 的 `entity_type` 与 typed `entity_id` 必须一一匹配。Task occurrence 永久没有 surrogate ID；删除 future pending occurrence 只能通过 `TaskOccurrenceTombstone(circle_id, task_id, member_id, scheduled_date, version, audit)` 投影。
- Completion、Cancellation、Ledger、Exchange 是 append-only。同 ID 同 payload为 no-op；同 ID 异 payload返回/记录 `AUDIT_INCONSISTENT`，整个本地事务回滚且 cursor 不推进。
- `next_cursor` 只用于取得同一 snapshot 的下一页；`snapshot_cursor` 锁定页集合。业务页面不能解析、拼接或排序 opaque cursor。
- 服务端已提交但响应丢失或本地投影失败时，客户端重放原 command；receipt/commit id 命中后恢复第一次 bundle，再原子投影。不得退化为 pull 后猜测命令是否成功。

## Read model

Main、Member Tasks/Rewards/Ranking 继续读 Room，不建立平行远端 BFF。balance batch、ledger/exchange/completion list 与 Statistics 都来自同一 canonical commit/snapshot 语义。Statistics 的 period、bucket、UTC range 与 IANA zone 必须显式提交；取消动作只调用 `cancelTaskCompletion`。
