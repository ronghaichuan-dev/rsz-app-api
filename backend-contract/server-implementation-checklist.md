# Server Implementation Checklist

以下清单用于从零实现 Clearwave Backend v1。勾选不替代 contract tests；任何实现选择都必须先满足 OpenAPI JSON，而不是反向修改客户端字段。

## 0. 合同导入与生成

- [ ] 将 `openapi/clearwave-backend-v1.json` 作为唯一代码生成/路由/schema 输入，并拒绝外部 `$ref`。
- [ ] 在 CI 执行 `ai/aw/validate_backend_contract.py` 与聚焦测试。
- [ ] 固定 wire 主版本 `1`；未知 `X-Contract-Version` 以 operation 允许的 `PROTOCOL_ERROR` 失败。
- [ ] 不添加未在 spec 中出现的 path、method、request/response 字段或“临时”错误 body。
- [ ] 不配置 placeholder `servers`；base URL 从环境配置注入。

## 1. 真实环境配置（发布前必须提供）

- [ ] HTTPS base URL、证书/hostname 校验与环境路由。
- [ ] access/refresh/Guest token 签名与轮换密钥。
- [ ] Google OAuth issuer、audience、JWKS 获取/缓存策略。
- [ ] Google Play Developer API 服务凭据、package/product policy。
- [ ] object store bucket/prefix、预签名器、大小/type/digest 校验与清理策略。
- [ ] cursor/receipt/idempotency 加密或签名密钥。
- [ ] Live 与 Demo 的数据库、密钥、object store、principal 和命名空间隔离。

以上四类外部互操作核心配置（HTTPS origin、OAuth、Play、asset store）当前没有真实值。缺少时相应 operation 返回明确 unavailable；不得用 fixture secret 上线。

## 2. 存储模型与约束

- [ ] Account、Session、AccountBinding 分表/aggregate；token secret 只保存安全 verifier/加密值，支持 revoke 与 rotation。
- [ ] Circle、Membership、Administrator、Member、CircleSelection 独立 version/status/tombstone；每个 Circle 至少一个 Owner。
- [ ] Membership 对 `(account_id, circle_id)` 与 actor binding 建立受控唯一性，授权时总是读取 active row。
- [ ] Invite 保存 target、generation、single-use、expiry、status、permission scope 与安全 code verifier；consume 使用唯一约束。
- [ ] TaskTag、TaskDefinition、assignment、TaskOccurrence 按 schema 存储；Occurrence 唯一键为 `(task_id, member_id, scheduled_date)`。
- [ ] TaskCompletion、TaskCancellation、LedgerEntry、ExchangeRecord 为 append-only；同 ID 同 payload可重复，同 ID 异 payload触发 `AUDIT_INCONSISTENT`。
- [ ] TaskCompletion 具有独立 `completion_id` 与 canonical `version`；取消用单独 Cancellation，不更新/删除完成 payload。
- [ ] StarBalance 是 canonical projection，具有 version 与 source commit；所有 delta 与 Ledger 在同一事务更新。
- [ ] RewardDefinition、assignment、RewardCooldown 具有独立 version；Exchange 只能由 redeem transaction 产生。
- [ ] NotificationOutboxEntry 与兑换同事务以 `version=1` 写入；worker 更新投递状态/attempt 时提升 version/updated time 并发布独立 commit，不改兑换事实。
- [ ] AssetUpload pending metadata 与 committed Asset 分离；业务外键只引用 committed asset_id。
- [ ] Commit/commit_sequence/change cursor 记录一次跨域事务的完整 changes；顺序单调且可重放。
- [ ] Idempotency record 保存 lookup scope、route/body fingerprint、状态、receipt 与加密响应 snapshot，跨进程恢复。
- [ ] Feedback、purchase verification receipt、EntitlementSnapshot 按 schema 与数据分类隔离。
- [ ] soft delete/tombstone 不联级物理删除 append-only 审计或其必要快照。

## 3. ID 与 legacy 兼容

- [ ] 对每个 OpenAPI ID schema执行完整 pattern/length 校验，同时把合法 ID 当 opaque string。
- [ ] 原样接受 `admin:v1:<uuid>` 与 `admin:v1:migrated:circle:v1:<uuid>`，禁止重写。
- [ ] `account_id`、`session_id`、`membership_id` 使用服务端独立命名空间；不接收/生成本地 IdentityScope stableId。
- [ ] 支持 `x-clearwave-id-mapping-policy` 冻结的 legacy Completion UUIDv5 结果；映射 ID 只作 identifier，不把 pre-v1 本地 audit 当服务端 canonical。
- [ ] 验证 path/body 中重复出现的 Circle/Member/Task 等 ID 一致且属于 canonical aggregate。

## 4. 通用 HTTP middleware

- [ ] 每个 operation 校验 `X-Request-Id`、`X-Contract-Version`、`X-Client-Version`；response 回显 request id 与服务端时间。
- [ ] 按 operation `security`/`x-auth-context` 校验 public、Guest、refresh、account、membership 与 purpose-aware asset owner credential。
- [ ] 所有 mutation 在领域代码前执行 durable idempotency lookup；JCS/SHA-256 与 route fingerprint 精确符合 OpenAPI extension。
- [ ] 每个 mutable mutation 在事务内验证 `x-version-field`；冲突回传 `current_version`。
- [ ] 请求/响应 JSON 严格拒绝未知字段、错误类型、非法 date/zone/ID/cursor/enum/size。
- [ ] 成功统一使用该 operation 的具名 success envelope；失败统一使用 `ErrorEnvelope`。
- [ ] 只返回 operation 声明的 status/code；内部或第三方异常不能产生自由文本 body。
- [ ] 访问日志只含 operationId/path template/status/耗时/request/trace id 等 allowlist 字段。

## 5. Auth 与 membership

- [ ] Google exchange 验证 signature/issuer/audience/expiry/nonce；不接受 provider subject 请求字段。
- [ ] Guest session 只签发 invite/bootstrap/feedback capability，并设置服务端 expiry/rate limit；create 响应同时签发绑定 session/environment/expiry 的 `guest_upgrade_grant`，current bootstrap 不回显；Feedback asset capability 不能访问 Circle asset。
- [ ] 新签发 token bundle 严格使用 account/invite_guest session variant；Google exchange、refresh、invite upgrade 不得返回 Guest principal，Guest create 不得返回 account principal。
- [ ] Guest upgrade/account binding 是单事务，失败不留下单边 account/session/membership。
- [ ] refresh 原子轮换 token；revoke 后 access/refresh 都 fail closed；两者均支持幂等重放。
- [ ] current account/bootstrap 按 principal 返回严格 oneOf：Account 分支返回 account/session/binding/memberships/selection/entitlement/cursor，Guest 分支只返回 session/capabilities/expiry/null cursor。
- [ ] 每个 Circle operation 在事务内重新解析 active membership actor 与 permission；ChooseRole 和请求 role 均不参与授权。

## 6. Circle、Invite 与普通 mutation

- [ ] onboarding 一次提交 Circle、Owner Administrator、initial Members、Membership、Selection 与 receipt。
- [ ] Circle/Member/Admin/selection 的 upsert/delete/leave 遵守 version、最后 Owner、soft delete、selection fallback 与 tombstone；`upsertMember` 只更新且 version 非空，新增必须走 `createCircleMember` 原子命令；Admin create intent 只接受 null version 并生成 `pending_invite`，update intent 只接受当前 version。
- [ ] outgoing invite code 由服务端生成；create 按 role 强制唯一 target id、非空 target version 与 Admin/Member permission 规则，refresh 受控返回 secret，revoke 不回显。
- [ ] Administrator/Member redeem 分别验证 role，单事务 consume invite 并返回完整 bootstrap bundle/cursor。
- [ ] `OPTIMISTIC_MUTATION` 的 pending 与 RemoteCommitted 可由未来客户端明确区分；服务端不返回虚假成功。

## 7. Task、Ledger 与 Reward

- [ ] Task repeat/end/reminder/calendar 规则按唯一且三处相等的 IANA zone 解释；none 只能配 never，after count 是成员展开前的系列日期数，reminder anchor/local minute/time limit 组合按 schema 校验；Task update/delete 只影响 future boundary 之后的 pending occurrence。
- [ ] Task reminder 只同步配置；AlarmManager/Notification delivery 不进入后端。
- [ ] complete command 单事务重验 assignment/occurrence/proof/version，提交 Completion + positive Ledger + Balance 完整 bundle。
- [ ] cancel command读取原 Completion/Ledger，保留原事实并追加 Cancellation + reversal + Balance 完整 bundle。
- [ ] 人工星星调整是唯一受权 adjustment command；没有通用 append ledger API。
- [ ] RewardDraft 不上行；`upsertReward` 表示 Save，custom repeat 强制正整数 cooldown days、固定 repeat 强制 null，并从 canonical definition/assignment 生成 eligibility。
- [ ] redeem command 单事务重验 permission/assignment/cooldown/balance/price，提交 Ledger + Balance + Cooldown + Exchange + notification outbox。
- [ ] Exchange 只提供 list；路由表不存在 create/update/delete。

## 8. Sync、分页与 read model

- [ ] `getCircleBootstrap` 的所有页固定同一 snapshot；客户端可 staging 后整体应用。
- [ ] `pullCircleBootstrapDelta` 只返回完整、有序 `SyncCommit`；一个 commit 永不跨页。
- [ ] mutable changes 带 version/tombstone；append-only changes 带稳定 ID、payload 与 commit sequence；outgoing invite 与 versioned notification outbox 状态不得漏出对应 commit。
- [ ] 普通 tombstone 的 entity type/typed ID 必须匹配；future pending occurrence 删除使用 `(task_id, member_id, scheduled_date)` 专用 tombstone，禁止生成 surrogate occurrence ID。
- [ ] cursor 是 opaque、防篡改并绑定 principal/Circle/query/snapshot；malformed/mismatched cursor 明确 validation failure。
- [ ] `next_cursor`、`has_more`、`snapshot_cursor` 在每个分页 operation 中语义一致，limit 上限按 spec。
- [ ] balance batch、ledger/exchange/completion list、stats/compare 从 canonical snapshot 读取，不由客户端余额求和替代。
- [ ] Statistics 严格按显式 UTC range、period/bucket 与 IANA zone 生成，包括零 bucket；compare 双方使用同一 snapshot。
- [ ] Main 无新 BFF；未来客户端完整投影后仍从 Room 读页面。

## 9. Asset、Support 与 Entitlement

- [ ] asset prepare 不接受 URL/path/bytes；Circle purpose 强制 non-null Circle membership，Feedback purpose 强制 `circle_id=null` 与 Guest/account capability；Live 只签发短期预签名 PUT，Demo 只用 `demo_transport`。
- [ ] asset commit 读取 prepare 时绑定的 principal/Circle/purpose 与实际 object metadata/digest，校验 owner/purpose/version 后签发持久 asset_id。
- [ ] Feedback 校验 consent/contact 组合、长度、attachment asset 权限；content/contact 不进入日志。
- [ ] Play verification 不接受 order/price/status/expiry claim；向 Google 验证后一次提交 receipt + entitlement。
- [ ] entitlement snapshot 只由 canonical verification/服务端 policy产生，客户端 cache/time 不能放行服务端权益。
- [ ] 不实现诊断日志上传、Device Push 或任何 external SDK 替身 API。

## 10. Conformance 与故障测试

- [ ] 每个 operation 至少执行一份成功 fixture，逐字段 schema 比较。
- [ ] 每个 mutation验证 first commit、同 key replay、同 key异 payload conflict、并发 expected version冲突。
- [ ] 每个 `x-required-fixture-errors` 都由真实领域分支产生，而非测试层硬编码 response。
- [ ] 模拟 commit 成功后断开连接，确认重放返回同 receipt/bundle且不重复写入。
- [ ] 模拟本地/下游投影失败，确认 cursor 不推进且完整 commit 可再次取得。
- [ ] 对 Completion/Ledger/Exchange 同 ID 异 payload进行 fail-closed 与审计告警测试。
- [ ] 对 invite consume、task complete/cancel、reward redeem进行并发测试，证明无双消费/双扣/单边事实。
- [ ] 对 token/invite/purchase/feedback/asset URL/user text 做日志与 error redaction 测试。
- [ ] 证明 Offline 模式自建 Backend 调用计数为 0，Demo/Live 数据无法互读。
- [ ] 真实 Google/Play/object-store E2E 只在真实配置存在时运行；否则结果必须是 `not_run`/unavailable，不得标记通过。

## 发布停止条件

出现以下任一项不得宣称 v1 可用：capability 无唯一归属、OpenAPI/fixture/实现字段漂移、关键事务拆成多个客户端写、幂等结果不耐久、sync 拆 commit、append-only 审计可被覆盖、Offline 发出请求、Demo 冒充外部成功、或真实配置缺失却报告互操作通过。
