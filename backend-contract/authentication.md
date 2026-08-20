# Authentication、Session 与 Membership

本合同把设备上的 `IdentityScope`、Backend account session 和每个 Circle 的 membership actor 明确分离。服务端授权只能来自已验证 token、当前 session、active membership 与 canonical entity；客户端的 scope stableId、ChooseRole、operator role、Member/Admin 快照都不是授权事实。

## Principal 层级

1. `invite_guest` 是服务端签发、短期、受限 principal，只能执行 session 声明的 invite redeem/current bootstrap/feedback capability。
2. `account` 是 Google proof exchange 后的服务端账号 principal；`account_id`、`session_id` 与 `membership_id` 均由服务端生成。
3. `membership` 将一个 account 绑定到一个 Circle 内的唯一 actor。`actor_type`、`actor_id`、`role`、`permissions`、`status` 每个 Circle 独立解析。
4. audit `ActorSnapshot` 在事务内从 token 与 canonical membership 生成。请求中的任何用户/角色字段都不能覆盖它。

`ChooseRole` 只表达客户端入口意图。它不能创建 membership，也不能把一个 Circle 的 Owner/Admin/Member 权限复制到另一个 Circle。

## 请求鉴权类别

| `x-auth-context` | Bearer 规则 | 典型 operation |
|---|---|---|
| `public` | 不带 Authorization | `createInviteGuestSession` |
| `public_or_guest` | 可匿名，也可带受限 Guest token | `exchangeGoogleProof` |
| `refresh` | Bearer 值是 refresh credential，不是 access token | `refreshSession` |
| `guest_or_account` | 必须是 active Guest 或 account access token | invite redeem、feedback、current bootstrap |
| `account` | 必须是 active account access token | Circle list/onboarding、selection、entitlement、session revoke |
| `membership` | account/Guest token 后还必须解析目标 Circle 的 active membership 与 capability | Circle/Task/Ledger/Reward/Statistics 等 Circle operation |
| `asset_owner` | 按 purpose 解析 owner：Circle asset 必须有 active membership；Feedback attachment 必须是获授权 Guest/account 且 `circle_id=null` | `prepareAssetUpload`、`commitAssetUpload` |

所有类别仍必须提交三项通用版本/关联 header；所有 mutation 仍必须提交 `Idempotency-Key`。OpenAPI operation 的 `security` 与 `x-auth-context` 是最终判定依据。

## Google proof exchange

`exchangeGoogleProof` 只接收：

- 短生命周期 `google_id_token` proof；
- 防重放 `proof_nonce`；
- 可空的 `guest_upgrade_grant`。

客户端不上传 provider subject、Google user id、email 推导账号、IdentityScope stableId 或客户端角色。服务端必须验证 token 签名、issuer、audience、expiry、nonce 与账号禁用状态，再从 proof 内部解析 provider binding。provider subject 不进入业务 wire 或日志。

无 `guest_upgrade_grant` 时，服务端创建或恢复 account/session/binding；有 grant 时，exchange 在同一事务中校验 Guest、绑定/升级 account，并返回 `AuthSessionBundleData` 所需的 account、session、account binding、memberships 与 bootstrap cursor。升级失败不得留下单边账号或丢失已兑换 membership。grant 只能来自 `createInviteGuestSession` 的 required 响应字段，必须绑定 Guest session、环境与到期时间，并在成功升级、session revoke 或到期后失效。

真实 Google issuer/audience 尚是部署配置；配置缺失应返回 `UNAVAILABLE`，不能让任意合成 proof 在 Live 环境成功。

## 受限 Guest bootstrap

`createInviteGuestSession` 以 `purpose=invite_redeem` 和 `client_nonce` 创建短期 session。成功响应同时签发 `guest_upgrade_grant`；客户端必须把它与 Guest access token 一起写入后续任务定义的凭据安全存储，不能从 token、session id 或本地 scope 猜测生成。`getCurrentAccount` 不重新回显该 secret。服务端返回的 capability 集仅允许：

- `invite_redeem_administrator`
- `invite_redeem_member`
- `current_account_bootstrap`
- `submit_feedback`
- `feedback_asset_upload`

Guest 不能创建 Circle、修改资料、执行 Task/Ledger/Reward 写入、验证 Play purchase 或读取其他 Circle。Feedback asset 只能使用 `purpose=feedback_attachment`、`circle_id=null`，且 commit 必须回到 prepare 时绑定的同一 principal；该能力不能借道上传 Circle asset。Guest refresh token 可以为 `null`；到期或撤销后 fail closed。

## Session 生命周期

- 新签发完整 token bundle 必须使用判别后的 `AccountAuthSession` 或 `InviteGuestAuthSession`；Google exchange、refresh 与 Guest invite upgrade 只能返回 account variant，Guest create 只能返回 invite_guest variant。通用 `AuthSession` 只作为两者的共同字段基座，不能直接替代 operation 响应类型。
- `AuthSession` 的 access/refresh token 是 opaque secret；客户端只存入后续任务定义的凭据安全存储，不得进入普通 `Persistent`、Room 业务表、日志、Crashlytics 或 fixture。
- `refreshSession` 使用 refresh credential 与 body 中的 `session_id`，原子轮换 token pair；旧 refresh credential 在成功提交后失效。幂等重放返回第一次轮换的同一安全结果。
- `revokeSession` 接受稳定 revoke reason；服务端提交 revocation 后，旧 access/refresh token 均不可恢复使用。
- `getCurrentAccount` 按 token principal 返回显式 `oneOf`：account 分支含 `principal_kind=account`、account/binding/memberships/selection/entitlement；Guest 分支含 `principal_kind=invite_guest`、session/capabilities/expiry，且不得伪造 Account、binding 或 membership。发现 `ACCOUNT_DISABLED`、`TOKEN_EXPIRED` 或已 revoke 时客户端必须清除远端会话并 fail closed。
- token 过期时间和 `server_time_ms` 由服务端产生；客户端时间只用于 UX，不参与授权。

## Invite redeem 与 membership

Administrator 与 Member 使用两个明确 operation，以防入口角色含糊。每个 redeem 都是单一 validate+consume 命令：邀请码查找、目标角色/目标状态、single-use、expiry、account/Guest 状态与 membership 冲突必须在同一事务校验。

Outgoing invite 必须绑定同 Circle 的 canonical target 与非空 `expected_target_version`：Administrator invite 只带 `target_administrator_id` 和非空 permission scope，Member invite 只带 `target_member_id` 且 permission scope 为空。Administrator target 由 `upsertAdministrator` 先建立 `pending_invite` profile；Member target 必须是 active。role、两种 target id 或 permission scope 交叉混用统一拒绝，不能依赖页面入口猜测。

成功 bundle 必须返回：必要时的 `AuthSession`、membership、canonical Circle、绑定的 Administrator 或 Member、redemption receipt 与 change cursor。服务端不能先返回“验证成功”再由客户端拼第二个 consume/write 请求；这会产生 TOCTOU 与单边 membership。

## Account 与本地 scope 的永久边界

- `AccountBinding` 只描述 Backend account 与某个本地 projection 环境的绑定策略，不携带 scope stableId。
- Offline 模式永远不创建 Guest/account session，也不调用任何自建后端 operation。
- Google subject 的本地不可逆摘要域不会转换为 `account_id`；服务端 `account_id` 也不能反向替换本地 scope 的隔离职责。
- Demo 与 Live 的 session、account、membership、幂等记录、cursor 和 asset 命名空间必须隔离；切换环境必须撤销/清空对应会话投影，不能跨环境复用 token。
