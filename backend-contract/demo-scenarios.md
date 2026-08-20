# Demo Scenarios v1

Demo 是同一 OpenAPI wire 的有状态本地/测试实现，不是另一套宽松 DTO，也不是对 Google、Play、Firebase 或 object store 真实互操作的证明。Demo handler 的预期 owner 已记录在 coverage JSON；本任务只冻结行为，不接入运行时。

## Demo 不变量

- 使用与 Live 完全相同的 path、method、header、schema、错误、版本、幂等和 commit bundle。
- Demo account/session/membership、数据库、cursor、idempotency record、clock、asset 与 Live 完全隔离。
- 场景使用 fixture 中的合成 token、邀请码、ID、文本和摘要；任何 `fixture_*` proof 只能被 Demo verifier 接受，Live 必须拒绝。
- 所有 mutation 持久化，进程重启后以相同 key 重放仍返回同一 receipt/commit/bundle。
- `server_time_ms` 由可控 Demo clock 产生；业务规则不能读取设备当前时间绕过 expiry/cooldown/version。
- Offline transport guard 的调用计数始终为 0；Demo 被选择才允许调用 Demo backend。
- 外部 SDK 场景只报告 `simulated_contract_result`、`external_interop_not_run` 或明确 unavailable，绝不报告“Google/Play 已真实验证”。

## 冻结场景

| ID | 前置状态与动作 | 必须观察的结果 |
|---|---|---|
| `D01_AUTH_GUEST` | 无 session，调用 `createInviteGuestSession`，持久保存合成 Guest token/grant，再调用 current bootstrap | 只返回 invite/bootstrap/feedback capabilities；create 响应签发 `fixture_*` upgrade grant；current bootstrap 命中 Guest 分支且不重新回显 grant、不伪造 Account/membership；Guest 访问普通 Circle mutation 得 `FORBIDDEN` |
| `D02_AUTH_EXCHANGE` | Demo verifier 接受指定合成 proof；匿名升级必须使用 D01 签发且仍有效的 Guest upgrade grant | 一次提交 account、session、binding、memberships 与 cursor；不出现 provider subject 字段；明确标记外部 Google interop 未运行 |
| `D03_SESSION_RECOVERY` | refresh 首次响应丢失后同 key/body 重放 | token pair 与完整 response 逐字段等于第一次 commit；该 operation 不伪造 receipt/result_kind；异 body 得 `IDEMPOTENCY_CONFLICT` |
| `D04_ONBOARDING` | active account 无 Circle，提交 Circle + Owner + initial Members | 单 commit 返回完整 canonical bundle；Member limit/version/重复 ID 失败不留下单边行 |
| `D05_INVITE_REDEEM` | Owner create/refresh invite，Guest 通过对应角色 operation redeem | create/refresh 只回显受控合成 code；redeem 原子返回 session（需要时）、membership、Circle、绑定资料、receipt/cursor；used/expired/role mismatch/target deleted 均有稳定失败且无部分消费 |
| `D06_TASK_SERIES` | 创建 weekly Task，再从显式 future date 更新/删除 | 只替换边界日及之后 pending occurrence；过去和 completed/cancelled 保持；response 计数与 snapshot 一致 |
| `D07_PROOF_ASSET` | prepare `task_proof`、通过 Demo transport 放入合成 bytes、commit，再 complete Task | prepare 返回 `demo_transport`、`upload_url=null`；commit 得持久 `asset_id`；digest/purpose/size 不符得 `ASSET_REJECTED` |
| `D08_COMPLETE_CANCEL` | assigned Member 完成 occurrence，响应丢失后重放，再由授权 Admin 取消 | complete bundle 同时含 occurrence/completion/正向 ledger/balance；cancel 保留原事实并追加 cancellation/reversal/balance；重放不重复加减星 |
| `D09_LEDGER_ADJUST` | Owner/Admin 对 active Member 执行受权人工调整 | 只产生 adjustment ledger 与 canonical balance；越界失败；不存在任意 append ledger endpoint |
| `D10_REWARD_REDEEM` | assigned Member 余额足够且无 cooldown，兑换 Reward | ledger/balance/cooldown/Exchange/notification outbox 同一 commit；重放不重复扣星；余额不足、未分配、冷却、永久不可用分别返回稳定 code |
| `D11_SYNC_RECOVERY` | 制造跨域 commit，按小 limit 拉 delta；本地投影在 commit 中途失败后重新拉 | 一个 commit 不跨页；失败事务不推进 cursor；重试同 commit 后完整投影；同审计 ID 异 payload触发 `AUDIT_INCONSISTENT` |
| `D12_READ_MODELS` | 同一 snapshot 查询 balance batch、ledger/exchange/completion list、stats/compare | opaque cursor 和 `snapshot_cursor` 一致；bucket 使用明确 UTC range/IANA zone；Main 无独立 BFF |
| `D13_FEEDBACK` | Guest/account 以 `circle_id=null` prepare/commit Feedback attachment 后提交合成 feedback，首次响应丢失后重放 | asset 与 receipt 持久且 owner/purpose 相容；日志不含 content/contact；无 consent 的 contact 或跨 purpose asset 得 validation failure/`ASSET_REJECTED` |
| `D14_ENTITLEMENT` | Demo Play verifier 被配置为 accepted/rejected/unavailable | accepted 原子产生 verification receipt + entitlement；其余稳定失败；结果必须标记为 Demo 模拟而非真实 Play interop |
| `D15_OFFLINE_ZERO_CALLS` | 切换 Offline 后遍历 Main、Task、Ledger、Reward、Statistics、Settings | 自建 Backend/Demo handler 调用计数严格为 0，页面只读写 Offline 本地事实 |

## Fixture 使用

`examples/*.json` 是场景输入/输出的 wire 基线：

- 每个 operation 至少一份成功 fixture；
- 每个 mutation 恰有 `success_first_committed` 与 `success_idempotent_replay`；
- operation 的 `x-required-fixture-errors` 必须逐一出现；
- invite code 仅允许 fixture allowlist，credential/user content 使用 `fixture_*` 合成值；
- fixture 中没有真实 token、购买订单、用户文本、路径、外部 URI 或生产 ID。

Demo 实现可以从 fixture 构建 contract tests，但不能把 fixture response 写死成无状态 stub。需要先执行领域校验和状态转移，再生成与 fixture schema 相容的结果；否则无法验证 replay、version、cooldown、invite consume、cursor 与 projection recovery。

## Asset 的持久 Demo 语义

`prepareAssetUpload` 为 `upload_id` 建立 pending metadata，并持久绑定 token principal、nullable Circle 与 purpose；Demo transport 以隔离 blob store 接收 bytes；`commitAssetUpload` 校验同一 owner、实际 size/SHA-256 后生成稳定 `asset_id`。重启后已 commit asset 仍可被 proof/profile/reward/feedback 引用和校验；重新 prepare 或 replay commit 不得生成第二个 asset。

Demo 不应返回假的 HTTPS upload origin。`upload_mode=demo_transport` 时 `upload_url` 必须为 null，由测试/应用内 Demo transport 使用 upload_id 定位隔离 blob。

## 场景重置与并发

- 每个测试使用独立 scenario namespace；reset 只清理该 namespace，不能删除其他测试或 Live 数据。
- 并发提交同 expected version 时只能一个 first commit，其他得到 `VERSION_CONFLICT` 或命中同一 idempotency record。
- 并发兑换、完成、取消、invite redeem 必须通过数据库约束与事务证明不会产生双扣、双加、双消费或单边 audit。
- Demo clock、ID seed 与初始 state 必须可注入以获得确定性 fixture；生产随机/签名实现不受 Demo seed 影响。
