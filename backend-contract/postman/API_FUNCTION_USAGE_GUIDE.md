# Clearwave Backend v1 接口功能与当前使用场景

## 1. 文档范围

本文解释 `Clearwave Backend Contract v1` 的全部 46 个自建后端 operation，回答两个问题：

1. 接口负责做什么；
2. 当前 Android 产品中的哪个功能会使用它，以及在什么时机触发。

接口 method、path 与 operationId 以 `backend-contract/openapi/clearwave-backend-v1.json` 为唯一事实源；本文只解释产品语义，不覆盖 OpenAPI。Postman 请求位于同目录的
`Clearwave-Backend-v1.postman_collection.json`。

当前总体状态是 `conditional_ready`：Online/Demo 客户端接线、HTTP adapter、Demo handler、Room projector 和测试已经闭合，但真实 Backend/TLS/OAuth/Play 服务端互操作尚未执行。下文的“当前使用”表示客户端已经预留并接入对应调用位置；只有配置真实
HTTPS 服务后才会发出 Live HTTP 请求。

## 2. 共同行为

- Offline 模式不会调用以下任何自建后端接口，仍使用本地 Room 事实源。
- Online/Demo 页面不直接持有 transport，也不把网络 DTO 直接显示为 UI 事实；请求统一经过领域 Services/repository/adapter，canonical 响应在事务中投影到 Room，页面再读取 Room。
- `POST`、`PUT`、`PATCH`、`DELETE` 都要求稳定 `Idempotency-Key`。同一逻辑重试必须复用原 key。
- 所有请求都带 `X-Request-Id`、`X-Contract-Version: 1` 和 `X-Client-Version`；除公开 Auth 入口外使用 Bearer credential。
- 页面跳转、提醒、相机/Picker、Firebase、Ads、Google Credential Manager 和 Play Billing 客户端不是这些接口的一部分。
- Reward 卡片颜色是 Android 本机展示元数据，不进入 `upsertReward`、canonical Reward 或 Postman 请求。
- Statistics 的 24 小时日内图由客户端在完成 Online 同步与日级 parity 后读取 Room 中的 Completion/Ledger 原始事实聚合；当前没有 `hour` bucket 或独立日内接口。

## 3. 首次启动与后续启动会涉及哪些接口

### 3.1 Online + assign task + 邀请码加入

无现有账号凭据时先调用 `createInviteGuestSession` 建立受限 Guest，再由 `AdminInviteActivity` 调用 `redeemAdministratorInvite`。成功后通过 `getCurrentAccount`、`pullCircleBootstrapDelta` 恢复完整账号/Circle
投影，并通过 `getCurrentEntitlement` 进入试用/订阅/Main gate。

### 3.2 Online + assign task + 非邀请码登录

`LoginActivity` 取得 Google proof 后调用 `exchangeGoogleProof`，随后使用 `getCurrentAccount` 复核 binding；`CreateCircleActivity` 保存时调用 `completeOnboarding`，再通过 `pullCircleBootstrapDelta` 和
`getCurrentEntitlement` 进入 Main/权益流程。

### 3.3 Online + complete task + 任一登录方式

邀请码路径使用 `createInviteGuestSession`（需要时）和 `redeemMemberInvite`；非邀请码账号路径使用 `exchangeGoogleProof`。成员资料与头像流程按实际保存动作使用 `createCircleMember`/`upsertMember` 以及 Asset
两阶段接口，完成同步后再读取 entitlement。

### 3.4 后续启动

有可恢复 credential 时，系统按需使用 `refreshSession`，然后调用 `getCurrentAccount`、`pullCircleBootstrapDelta` 和 `getCurrentEntitlement`。用户主动退出时 best-effort 调用 `revokeSession`。

## 4. Auth：身份、账号与 Session（5）

| operationId                | Method / Path                                | 做什么                                                                                                  | 当前在哪个功能中使用                                                                                                 | 主要客户端 owner                                             |
|----------------------------|----------------------------------------------|------------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------------------|---------------------------------------------------------|
| `exchangeGoogleProof`      | `POST /v1/auth/google:exchange`              | 把短生命周期 Google ID token proof 交换为 Clearwave account、binding、membership 与 Session；Guest 升级也在同一原子命令内完成。 | Online 非邀请码登录；`LoginActivity` 取得 Credential Manager proof 后调用。邀请码 Guest 后续绑定 Google 账号时也复用。                | `BackendSessionCoordinator`、`RuntimeAuthBackendAdapter` |
| `createInviteGuestSession` | `POST /v1/auth/guest-sessions`               | 创建只能执行邀请码 bootstrap/redeem 的受限 Guest Session，并签发升级 grant。                                            | Admin/Member 邀请码登录且当前没有账号 credential 时，在进入 `AdminInviteActivity` 或 `MemberInviteActivity` 前使用；不能进入普通 Main。 | `BackendSessionCoordinator`、`RuntimeAuthBackendAdapter` |
| `refreshSession`           | `POST /v1/auth/sessions:refresh`             | 轮换 refresh credential，签发新的 access/refresh token，并保持同一账号 binding。                                     | 后续启动、token 临近过期，以及 HTTP 401 `TOKEN_EXPIRED` 的 single-flight 自动恢复；不是用户直接点击的页面接口。                            | `BackendSessionCoordinator`、credential store            |
| `revokeSession`            | `POST /v1/auth/sessions/{session_id}:revoke` | 撤销远端 Session；同一幂等请求重放首次 revoke receipt。                                                              | 用户退出账号、环境/账号切换或凭据清理时 best-effort 调用。远端失败不会恢复本地已退出状态。                                                       | `BackendSessionCoordinator`、`IdentityServices`          |
| `getCurrentAccount`        | `GET /v1/account/bootstrap`                  | 读取 current account、Session metadata、binding 和 memberships 的 bootstrap。                               | 登录/refresh 成功后的闭合校验、后续启动恢复、邀请码已提交但本地投影失败后的恢复，以及进入 Main 前的账号复核。                                             | `BackendSessionCoordinator`、`RuntimeAuthBackendAdapter` |

## 5. Circle：圈子、成员、管理员与选择（12）

| operationId           | Method / Path                                                      | 做什么                                                                                                   | 当前在哪个功能中使用                                                                                   | 主要客户端 owner                                                                   |
|-----------------------|--------------------------------------------------------------------|-------------------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `completeOnboarding`  | `POST /v1/circles:onboard`                                         | 一次性创建 Circle、Owner Administrator、初始 Members、Membership 与当前选择，返回完整 canonical bundle。                   | Online assign task + 非邀请码登录后的 `CreateCircleActivity` 首次建圈；只有该原子提交和本地投影完整成功后才能进入权益/Main gate。 | `CircleMemberServices`、`CircleMemberRepository`、`RuntimeCircleBackendAdapter` |
| `createCircleMember`  | `POST /v1/circles/{circle_id}/members`                             | 在指定 Circle 内创建新的 canonical Member，并校验 actor 权限与稳定 ID。                                                 | 首次成员资料创建、`CreateCircleMemberActivity`，以及管理端“新增成员”入口；编辑已有成员不能调用该接口。                           | `CircleMemberServices`、`CircleMemberRepository`                               |
| `listMyCircles`       | `GET /v1/circles`                                                  | 分页读取当前账号可见的 Circles 及对应 Membership。                                                                   | Main 的 Circle 切换器、登录恢复后的可用圈子列表，以及当前 Circle 失效后的 fallback 选择。                                 | `CircleMemberRepository`、`RuntimeCircleBackendAdapter`                        |
| `getCircleBootstrap`  | `GET /v1/circles/{circle_id}/bootstrap`                            | 在同一 snapshot 下分页读取 Circle identity aggregate，包括 Circle、Admin、Member、Membership、selection 与 invite 状态。 | 首次进入或切换 Circle、邀请兑换恢复、Main 进入前补齐身份投影；成员的 canonical `bound_account_id` 同时用于派生已绑定/未绑定展示状态。     | `CircleMemberRepository`、`CircleCanonicalProjector`                           |
| `selectCurrentCircle` | `PUT /v1/account/circle-selection`                                 | 同步账号当前选择的 Circle，并返回 canonical selection/version。                                                     | 用户在 Main 的 Circle 切换功能中选择另一个圈子时调用，使选择可跨设备恢复。                                                 | `CircleMemberServices`、`CircleMemberRepository`                               |
| `updateCircle`        | `PATCH /v1/circles/{circle_id}`                                    | 更新 Circle 名称、图标和版本，服务端复验管理权限与 expected version。                                                       | Circle 资料/设置页面的改名、修改图标操作。                                                                    | `CircleMemberServices`、`CircleMemberRepository`                               |
| `deleteCircle`        | `DELETE /v1/circles/{circle_id}`                                   | 软删除 Circle，并返回仍可用的 selection fallback；不删除历史审计事实。                                                      | Owner 在 Circle 管理页面确认删除圈子时使用。                                                                | `CircleMemberServices`、`CircleMemberRepository`                               |
| `leaveCircle`         | `POST /v1/circles/{circle_id}:leave`                               | 原子结束当前账号的 Membership，并选择可用 fallback Circle。                                                           | Member/Admin 在 Circle 设置中执行“退出圈子”。Owner/最后 Owner 等边界由服务端拒绝或仲裁。                               | `CircleMemberServices`、`CircleMemberRepository`                               |
| `upsertMember`        | `PUT /v1/circles/{circle_id}/members/{member_id}`                  | 更新既有 Member 的名称、性别、preset/asset avatar 与 canonical version。                                           | Member 资料编辑、`MemberVirtualAvatarActivity` 保存，以及管理端编辑现有成员；新增成员必须使用 `createCircleMember`。      | `CircleMemberServices`、`CircleMemberRepository`                               |
| `deleteMember`        | `DELETE /v1/circles/{circle_id}/members/{member_id}`               | 软删除 Member，并移除未来 Task/Reward assignment，同时保留历史 Completion/Ledger/Exchange 快照。                         | 管理端成员管理页面的删除成员操作。                                                                            | `CircleMemberServices`、`CircleMemberRepository`                               |
| `upsertAdministrator` | `PUT /v1/circles/{circle_id}/administrators/{administrator_id}`    | 创建 Administrator draft，或更新 Owner/Admin profile、权限和 canonical version。                                 | 新增管理员、管理员资料编辑、Owner/Admin 头像与名称保存。                                                           | `CircleMemberServices`、`CircleMemberRepository`                               |
| `deleteAdministrator` | `DELETE /v1/circles/{circle_id}/administrators/{administrator_id}` | 软删除 Administrator 及其 Membership，保留既有历史事实。                                                             | Owner 在管理员管理页面移除 Admin 时调用；last-owner 等限制由服务端复验。                                             | `CircleMemberServices`、`CircleMemberRepository`                               |

## 6. Invite：邀请码生成、更新、撤销与兑换（5）

| operationId                 | Method / Path                                              | 做什么                                                                          | 当前在哪个功能中使用                                                                                                                  | 主要客户端 owner                                                        |
|-----------------------------|------------------------------------------------------------|------------------------------------------------------------------------------|-----------------------------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------|
| `createInvite`              | `POST /v1/circles/{circle_id}/invites`                     | 由服务端生成一次性 secret、TTL、目标角色/成员/管理员和 generation，返回受控展示码。                        | Main 或 Member/Admin 邀请展示页首次生成跨设备邀请码。Online 不再使用本地六码冒充远端邀请码。                                                                 | `InviteBackendPort`、`RuntimeCircleBackendAdapter`                  |
| `refreshInvite`             | `POST /v1/circles/{circle_id}/invites/{invite_id}:refresh` | 原子撤销旧 secret，并为同一目标返回新的 generation、有效期和展示码。                                  | 邀请页面点击“刷新/更换邀请码”时使用。重复请求返回首次 replacement receipt。                                                                           | `InviteBackendPort`、`RuntimeCircleBackendAdapter`                  |
| `revokeInvite`              | `POST /v1/circles/{circle_id}/invites/{invite_id}:revoke`  | 撤销尚未使用的 outgoing invite，冻结其后续兑换。                                             | 邀请页面主动撤销、目标关系被删除或管理员希望立即使邀请码失效时使用。                                                                                          | `InviteBackendPort`、`RuntimeCircleBackendAdapter`                  |
| `redeemAdministratorInvite` | `POST /v1/invites/administrator:redeem`                    | 一次性校验并消费管理员邀请码，原子返回 Admin、Circle、Membership、selection、Session/receipt。       | `online + assign task + 邀请码加入` 的 `AdminInviteActivity`。投影闭合后才返回 `ADMIN_JOINED`。                                             | `CanonicalAdminInviteGateway`、`CircleInviteRedemptionCoordinator`  |
| `redeemMemberInvite`        | `POST /v1/invites/member:redeem`                           | 一次性校验并消费成员邀请码，原子返回 bound Member、Circle、Membership、selection、Session/receipt。 | `online + complete task + 邀请码加入` 的 `MemberInviteActivity`；成功后进入 `CreateCircleMemberActivity → MemberVirtualAvatarActivity`。 | `CanonicalMemberInviteGateway`、`CircleInviteRedemptionCoordinator` |

## 7. Sync：Main 与领域增量同步（1）

| operationId                | Method / Path                      | 做什么                                                                             | 当前在哪个功能中使用                                                                                 | 主要客户端 owner                                                          |
|----------------------------|------------------------------------|---------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|----------------------------------------------------------------------|
| `pullCircleBootstrapDelta` | `GET /v1/circles/{circle_id}/sync` | 按不可拆分的 server commit 拉取完整 bootstrap 或 change cursor 之后的增量；一页中的一个 commit 必须整体投影。 | Startup、登录/Invite 恢复、Main 刷新和 mutation ACK 后统一同步；投影的 Completion/Ledger 原始事实也供本机 24 小时统计使用。 | `BackendSyncServices`、`BackendSyncCoordinator`、`AtomicSyncProjector` |

## 8. Asset：业务图片和证明文件（2）

Asset 上传是两阶段流程：先 prepare，客户端把 bytes 上传到返回的 target，最后 commit。业务接口只保存 committed `asset_id`，绝不传 URI、绝对路径或图片内容。

| operationId          | Method / Path                                | 做什么                                                                        | 当前在哪个功能中使用                                                                  | 主要客户端 owner                                            |
|----------------------|----------------------------------------------|----------------------------------------------------------------------------|-----------------------------------------------------------------------------|--------------------------------------------------------|
| `prepareAssetUpload` | `POST /v1/assets/uploads:prepare`            | 按 purpose、owner、content type、大小与 SHA-256 校验上传意图，返回预签名或 Demo upload target。 | Task 照片证明、Member/Admin 头像、Reward 图片和 Feedback attachment 在最终保存前需要上传本地文件时使用。 | `BackendAssetServices`、`RuntimeBackendAssetRemotePort` |
| `commitAssetUpload`  | `POST /v1/assets/uploads/{upload_id}:commit` | 校验已上传 bytes 与 prepare metadata 一致，签发可长期引用的 committed `asset_id`。           | 上述文件实际上传完成后自动调用；只有成功返回的 asset ID 才能继续提交 Task、Profile、Reward 或 Feedback。     | `BackendAssetServices`、`RuntimeBackendAssetRemotePort` |

## 9. Task：标签、任务定义与 occurrence（5）

| operationId           | Method / Path                                            | 做什么                                                          | 当前在哪个功能中使用                                                                       | 主要客户端 owner                                                   |
|-----------------------|----------------------------------------------------------|--------------------------------------------------------------|----------------------------------------------------------------------------------|---------------------------------------------------------------|
| `upsertTaskTag`       | `PUT /v1/circles/{circle_id}/task-tags/{task_tag_id}`    | 创建或更新 Task Tag，并返回 canonical Tag/version。                    | Task 标签管理，以及 `TaskEditorActivity` 保存新标签或重命名标签。                                   | `TaskServices`、`TaskMutationCoordinator`、`TaskBackendAdapter` |
| `deleteTaskTag`       | `DELETE /v1/circles/{circle_id}/task-tags/{task_tag_id}` | 软删除 Task Tag，并按合同处理未来 Task definition 的引用。                   | 标签管理页面确认删除标签时使用；不会改写历史 occurrence 快照。                                            | `TaskServices`、`TaskMutationCoordinator`                      |
| `upsertTask`          | `PUT /v1/circles/{circle_id}/tasks/{task_id}`            | 保存 Task definition、成员分配、星星、重复/结束规则、照片要求、Tag 和本地 reminder 配置。 | `TaskEditorActivity` 用户明确点击保存；输入过程和草稿变化不会持续上传。                                   | `ModeAwareTaskEditorRepositoryPort`、`TaskMutationCoordinator` |
| `deleteTask`          | `DELETE /v1/circles/{circle_id}/tasks/{task_id}`         | 软删除 Task，只取消边界后的 pending occurrences，并保留完成/取消历史。             | Task 管理/编辑页面确认删除任务时使用。                                                           | `TaskServices`、`TaskMutationCoordinator`                      |
| `listTaskOccurrences` | `GET /v1/circles/{circle_id}/task-occurrences`           | 按日期范围、成员和 opaque cursor 分页读取 canonical occurrences。          | Member 读取自身任务，以及 Owner/Admin 在 Main 选择成员后读取其任务；投影后也用于重建本机 AlarmManager reminder。 | `TaskBackendAdapter`、`TaskCanonicalProjector`                 |

## 10. TaskLedger：任务完成/取消的跨域原子命令（2）

| operationId            | Method / Path                                                          | 做什么                                                                            | 当前在哪个功能中使用                                                                                              | 主要客户端 owner                                                                          |
|------------------------|------------------------------------------------------------------------|--------------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------|
| `completeTask`         | `POST /v1/circles/{circle_id}/task-completions`                        | 原子完成 occurrence，并在同一 commit 中生成 Completion、正向 Ledger、Balance、receipt 与 cursor。 | Member 完成自身任务，或 Owner/Admin 在 Main 替选中成员完成；照片证明先取得 `asset_id`。请求只传目标 `member_id`，actor 由服务端 Session 推导。 | `StarLedgerServices.taskCompletionTransactions`、`OnlineTaskLedgerService`            |
| `cancelTaskCompletion` | `POST /v1/circles/{circle_id}/task-completions/{completion_id}:cancel` | 保留原 Completion/Ledger，原子追加 Cancellation、负向 reversal 和新 Balance。                | Task completion detail/Statistics 明细页面执行取消；P11 只复用该命令，不另建取消接口。                                          | `StatisticsCompletionCancellation` → `StarLedgerServices.taskCompletionTransactions` |

## 11. Ledger：余额、流水与人工调整（3）

| operationId            | Method / Path                                   | 做什么                                               | 当前在哪个功能中使用                                                       | 主要客户端 owner                                                           |
|------------------------|-------------------------------------------------|---------------------------------------------------|------------------------------------------------------------------|-----------------------------------------------------------------------|
| `adjustMemberStars`    | `POST /v1/circles/{circle_id}/star-adjustments` | 由授权 Owner/Admin 人工增减成员星星，并原子追加 Ledger/Balance。    | `StarBalanceActivity` 或成员余额管理中的手工加星/扣星操作。                        | `StarLedgerServices.balanceAdjustments`、`OnlineStarAdjustmentService` |
| `getMemberBalances`    | `GET /v1/circles/{circle_id}/member-balances`   | 在同一 snapshot 中批量读取成员余额，未产生流水的成员返回零。               | `MemberRankingFragment` 排名刷新，以及需要批量余额的 Main/成员管理投影；避免逐成员 N+1 请求。 | `StarLedgerBackendIntegration`、`LedgerCanonicalProjector`             |
| `listStarTransactions` | `GET /v1/circles/{circle_id}/star-transactions` | 按稳定 snapshot/cursor 分页读取 append-only Star Ledger。 | Star Balance/流水历史、成员账本列表及相关统计同步。页面最终读取投影后的 Room Ledger。          | `RuntimeStarLedgerBackendPort`、`BackendReadModelAdapter`              |

## 12. Reward：奖励定义、删除与资格（3）

| operationId            | Method / Path                                                 | 做什么                                                                                      | 当前在哪个功能中使用                                                                       | 主要客户端 owner                                      |
|------------------------|---------------------------------------------------------------|------------------------------------------------------------------------------------------|----------------------------------------------------------------------------------|--------------------------------------------------|
| `upsertReward`         | `PUT /v1/circles/{circle_id}/rewards/{reward_id}`             | 保存 Reward definition、成员分配、星星价格、重复/cooldown policy 与 committed asset reference；不包含本机卡片颜色。 | `CustomRewardActivity` 或奖励管理页面明确点击保存；Reward Draft 和 `RewardVisualColor` 均只在本机保存。 | `RewardServices`、`RewardMutationCoordinator`     |
| `deleteReward`         | `DELETE /v1/circles/{circle_id}/rewards/{reward_id}`          | 软删除 Reward，关闭未来 eligibility，但保留 Ledger、Cooldown 与 Exchange 历史快照。                         | Owner/Admin 在奖励管理页面确认删除奖励。                                                       | `RewardServices`、`RewardMutationCoordinator`     |
| `getRewardEligibility` | `GET /v1/circles/{circle_id}/rewards/{reward_id}/eligibility` | 读取成员维度的 assignment、余额条件、cooldown、一次性状态和 canonical eligibility。                           | Member 查看自身资格，或 Owner/Admin 在奖励管理页查看选中成员资格；兑换前复验也使用该读取，且不会写 cooldown。            | `RuntimeRewardEligibilityPort`、Reward repository |

## 13. Redemption：奖励兑换原子事务（1）

| operationId    | Method / Path                                     | 做什么                                                                                             | 当前在哪个功能中使用                                                                                     | 主要客户端 owner                                                            |
|----------------|---------------------------------------------------|-------------------------------------------------------------------------------------------------|------------------------------------------------------------------------------------------------|------------------------------------------------------------------------|
| `redeemReward` | `POST /v1/circles/{circle_id}/reward-redemptions` | 在一个服务端 commit 中扣星并生成 Ledger、Balance、Cooldown、不可变 Exchange、notification outbox、receipt 与 cursor。 | Member 兑换自身奖励，或 Owner/Admin 在奖励管理页替选中成员兑换。请求只传目标 `member_id`，actor 由服务端 Session 推导；通知失败不会回滚兑换。 | `RewardRedemptionServices.redemptions`、`OnlineRewardRedemptionService` |

## 14. ReadModels：历史、统计与完成明细（4）

| operationId                 | Method / Path                                    | 做什么                                                                       | 当前在哪个功能中使用                                                                                      | 主要客户端 owner                                                    |
|-----------------------------|--------------------------------------------------|---------------------------------------------------------------------------|-------------------------------------------------------------------------------------------------|----------------------------------------------------------------|
| `listExchangeHistory`       | `GET /v1/circles/{circle_id}/exchanges`          | 按 member/filter/zone 和 snapshot cursor 分页读取 append-only Exchange history。 | `ExchangeHistoryActivity` 的全部/月份筛选、分页与进程恢复。                                                     | `ExchangeServices.historyRepository`、`BackendReadModelAdapter` |
| `getStatistics`             | `GET /v1/circles/{circle_id}/statistics`         | 按明确 range、bucket、metric、week start 和 IANA zone 聚合单成员统计。                   | Statistics 趋势/汇总页面当前使用自然日 bucket 并与 Room 逐 bucket 校验；选中日期后的 24 小时图仅聚合已同步的本地原始事实，不调用小时接口。        | `StatisticsUiDataSource`、`BackendReadModelAdapter`             |
| `compareStatistics`         | `GET /v1/circles/{circle_id}/statistics:compare` | 在同一 snapshot、range 和 bucket 下比较两个成员。                                      | Member Comparison/Statistics 对比页面；双方数据必须使用相同 cursor 与时间口径。                                      | `StatisticsUiDataSource`、`BackendReadModelAdapter`             |
| `listTaskCompletionDetails` | `GET /v1/circles/{circle_id}/task-completions`   | 分页读取 append-only Completion、Cancellation、reversal 和 proof asset 明细。       | `TaskCompletionDetailsActivity` 的筛选，以及 Statistics 选中日内区间后的任务明细；取消动作另行调用 `cancelTaskCompletion`。 | `StatisticsUiDataSource`、`BackendReadModelAdapter`             |

## 15. Support：结构化反馈（1）

| operationId      | Method / Path               | 做什么                                                            | 当前在哪个功能中使用                                                                | 主要客户端 owner                                         |
|------------------|-----------------------------|----------------------------------------------------------------|---------------------------------------------------------------------------|-----------------------------------------------------|
| `submitFeedback` | `POST /v1/support/feedback` | 提交经过隐私同意的结构化反馈、可选联系方式和最多三个 committed attachments，返回幂等 receipt。 | `FeedbackActivity` 的 Backend 提交链；只有 Backend 明确不可用时才允许转到系统邮件，邮件发送不能冒充接口成功。 | `FeedbackSubmissionAdapter`、`SupportBackendAdapter` |

## 16. Entitlement：Play 验签与会员/试用资格（2）

| operationId             | Method / Path                              | 做什么                                                                      | 当前在哪个功能中使用                                                                                        | 主要客户端 owner                                                                   |
|-------------------------|--------------------------------------------|--------------------------------------------------------------------------|---------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| `verifyPlayPurchase`    | `POST /v1/entitlements/google-play:verify` | 由服务端向 Google Play 验证 purchase token，并原子提交 canonical entitlement/receipt。 | Play Billing 购买或恢复得到候选 Purchase 后、acknowledge/consume 和成功权益回调之前。Demo Backend 不会自动启用 Billing fake。 | `BillingService` 的 `ServerPurchaseVerifier` → `RuntimeBackendEntitlementPort` |
| `getCurrentEntitlement` | `GET /v1/entitlements/current`             | 读取服务端 verified entitlement、状态、有效期和 trial eligibility。                    | Startup、登录/邀请完成后的权益 gate、`FreeTrialGuideActivity`、`GuideSubscriptionsActivity`、订阅页和后续启动。          | `BackendMergedEntitlementSource`、`EntitlementProvider`                        |

## 17. 不属于这 46 个接口的当前能力

以下能力已有产品入口，但 owner 不是 Clearwave 自建后端，因此不应在 Postman 清单中寻找对应接口：

- Google 登录 UI 与 proof 获取：Android Credential Manager / Google Identity；
- 商品目录、价格、购买、恢复、acknowledge/consume：Google Play Billing；
- Analytics、Remote Config、Crashlytics：Firebase；
- Ads/UMP、TikTok attribution：对应第三方 SDK；
- Task/每日提醒及首次通知提示：AlarmManager、NotificationManager 与本地 registry；
- Reward 卡片颜色：Room 中的本机展示元数据，不进入 canonical Reward；
- 相机、Photo Picker、FileProvider、分享、邮件和应用市场：Android 平台/外部应用；
- 语言、主题、week start、FAQ/About 等设置：本地 Persistent/resources；
- Diagnostic log upload 与 Device Push：当前合同明确退役或 `NOT_REQUIRED`，没有伪造接口。

## 18. 联调阅读顺序

1. 先用本文确认产品入口和触发时机；
2. 在 Postman Collection 中按同名 operationId 找到请求；
3. 用 OpenAPI 确认所有字段、错误码和鉴权；
4. 用 `backend-contract/examples/*.json` 验证成功、失败、幂等重放和权限 fixture；
5. 用 `backend-contract/server-implementation-checklist.md` 完成服务端事务与 staging 验收。
