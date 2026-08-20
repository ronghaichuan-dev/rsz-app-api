# Security 与 Redaction

本合同采用“服务端重新推导事实、最小 wire、默认不记录 payload”的边界。OpenAPI schema 允许某字段只表示业务必须传输，不表示该字段可以进入普通日志、Analytics、Crashlytics、错误消息或长期明文缓存。

## 信任矩阵

| 输入 | 服务端处理 |
|---|---|
| `Authorization` token | 验证 session、principal、expiry、revocation；不记录原文 |
| Google ID token | 只作短生命周期 proof；验证 issuer/audience/signature/nonce/expiry 后解析绑定，不接受独立 provider subject 字段 |
| Guest token/grant | `createInviteGuestSession` 分别签发 token 与 `guest_upgrade_grant`；只允许声明 capability，不能升级成任意 account/membership 权限 |
| Circle/Member/Admin/Task/Reward/Exchange 等 ID | 按 pattern 接受为 opaque 定位符，再验证 canonical ownership/membership；不能从前缀推导权限 |
| 客户端 role、operator、scope stableId | 不在可信请求 schema 中；出现未知字段即拒绝 |
| price、balance、stars delta、server time、audit actor | 永远从 canonical entity/token/server clock 推导 |
| `expected_*_version` | 只作并发 precondition，不作内容 authority |
| asset metadata | prepare/commit 两次校验；业务 operation 只接受已 commit 且 purpose/owner 匹配的 `asset_id` |
| Play purchase token | 向 Google Play Developer API 验证；不信任客户端 order、price、active/expiry claim |

所有 Circle mutation 必须在数据库事务内重新读取 session、membership、actor、target entity、assignment/entitlement 与 version。网关或页面加载时的权限快照不能替代事务内检查。

## 敏感数据分级

### Secret / credential

- access token、refresh token、Google ID token、guest upgrade grant；
- purchase token；
- invite code；
- 预签名 upload URL 及其签名/query；
- 服务端 Google/Play/object-store 凭据。

这些值不得记录、埋点、回显到错误或写入普通业务表。需要支持幂等重放的 secret response 必须加密静态存储、限制服务账号访问，并绑定环境/principal/key。

`guest_upgrade_grant` 与 Guest token 分开签发但同级加密保存；不得由 session id、本地 scope 或 token 文本推导，`getCurrentAccount` 也不回显。Live/Demo/environment、session、expiry 不匹配时必须 fail closed。

### 用户内容与高敏引用

- Feedback `content`、`contact`、attachment；
- Circle/Member/Admin/Task/Reward 用户可编辑名称、说明与快照；
- proof/profile/reward/feedback asset；
- actor display name snapshot。

允许在业务事实中按合同保存，但日志只记录受控枚举、数量、Boolean、operationId、HTTP status、耗时、trace/request id 与脱敏 fingerprint。不得记录全文、资源 URL、asset storage key 或可反推出用户的完整稳定 ID。

### 受控业务事实

版本、commit sequence、数量、状态枚举、错误 code、retryable、bucket 类型可用于排障。即使属于受控事实，也应遵循最小保留和按需访问。

## 日志与错误

允许的网络日志字段：`operation_id`、path template（不是带实际 ID 的 URL）、status、耗时、`request_id`、`trace_id`、重试次数、幂等命中 Boolean、稳定错误 code。

禁止：Authorization/header secret、完整请求/响应 body、query 中的 cursor、invite/purchase token、Google subject、用户输入、完整 ID、asset URL/path、第三方 response body、异常堆栈直接返回客户端。

服务端内部错误先映射为 [`error-catalog.md`](error-catalog.md) 的稳定 code。`field` 只能给 schema 字段位置，`trace_id` 只作受控关联；两者都不能带原值。

## Asset prepare/upload/commit

1. `prepareAssetUpload` 接受客户端生成 `upload_id`、purpose、nullable `circle_id`、content type、byte size 与 SHA-256；不接受文件名、URI、绝对/相对路径或图片 bytes。`task_proof`/`profile_avatar`/`reward_image` 必须带 Circle，`feedback_attachment` 必须令 `circle_id=null`。
2. Live 返回短期 `presigned_put` target；Demo 返回 `demo_transport` 且 `upload_url=null`。预签名 URL 只能保存在内存中的当前上传流程，禁止日志与业务持久化。
3. object store 只接受冻结的 method/content length/type；上传隔离在未提交命名空间，不能被业务读取。
4. prepare/commit 都使用 `asset_owner` 鉴权：Circle purpose 从 token 解析 active membership；Feedback purpose 从 token 解析具备 `feedback_asset_upload` 的 Guest/account。`commitAssetUpload` 必须回查 prepare 时绑定的 principal/Circle，再校验实际对象大小、digest、type、expiry、purpose 与 expected upload version，然后签发持久 `asset_id`。
5. Task proof、profile avatar、Reward image/Feedback attachment 只能引用 purpose 与 Circle/principal 相容的 committed asset。Feedback asset 不能被 Circle 业务引用，Circle asset 也不能附到 Feedback；失败返回 `ASSET_REJECTED`，不能自动降级为外部 URL/path。
6. 被拒绝或过期 upload 可以异步清理；已进入不可变审计快照的 asset 引用必须遵守审计保留策略，不能因资料删除而悬空。

## Invite、反馈与购买

- Invite code 仅在 create/refresh 的成功 response 与 redeem request 出现。服务端保存不可逆 verifier 或等价安全表示；如为幂等重放保存原 code，必须加密。revoke/used/expired 后仍保留必要审计，不回显 code。
- Feedback 必须校验 `privacy_consent_version`；只有 `contact_consent=true` 且 `contact_type=email` 时接受 contact。正文/contact 不进通用日志或 Analytics，attachment 复用 committed asset 权限。
- `client_created_at_ms` 仅为反馈上下文，不能替代 server receipt time。
- Play verification 只接受 package/product/type/token。order id、价格、entitlement active/expiry 不在请求 schema；服务端从 Google canonical response 与本地 product policy 决定 entitlement。

## 外部与本地边界

- Google Credential Manager、Play Billing client、Firebase/Remote Config/Crashlytics、Ads/UMP、TikTok、AlarmManager、本地通知、Picker/分享不属于自建 API。
- v1 诊断日志上传已安全退役：当前 production caller 为 0，且没有 consent/redaction/retention 产品合同。不得复活遗留 `serverUrl`/multipart path。
- Device Push 为 `NOT_REQUIRED`：当前没有 FCM dependency、token lifecycle、安装 actor 或 notification acknowledgement 产品入口。产品若新增入口，必须重新盘点并独立冻结，不能暗中塞入现有 schema。
- Offline 模式必须通过 transport guard 保证自建后端调用为 0；不能为了“同步”读取/写入 Online 或 Demo 命名空间。

## 部署最低要求

- 仅 HTTPS，明确证书/hostname 校验、现代 TLS 与受控代理策略；真实 origin 由环境配置注入。
- token、idempotency secret response、provider binding 与服务凭据加密静态存储；密钥轮换不改变 wire ID。
- Live/Demo 数据库、token signing key、object store bucket/prefix、cursor signing key 与幂等命名空间隔离。
- 对 public/Guest/auth/invite/feedback/asset/Play verification 设置独立限流和 abuse controls。
- 数据删除执行 soft-delete/tombstone 与审计保留边界；append-only Completion/Ledger/Exchange 不做用户资料联级物理删除。
- fixture 全部是 `fixture_*`/allowlist 合成值；禁止用生产导出替换测试样本。
