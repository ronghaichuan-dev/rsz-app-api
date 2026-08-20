# Error Catalog v1

错误码、HTTP status、`retryable` 与 operation 归属的 machine-readable 权威位于 OpenAPI 的 `x-clearwave-error-catalog`、各 operation 的 `x-error-codes` 和 response `x-error-codes`。本页只给服务端与客户端统一处理规则。

## Error envelope

每个失败响应使用 `ErrorEnvelope`：

- `contract_version`：固定 `1`；
- `request_id`：回显本次 `X-Request-Id`；
- `server_time_ms`：服务端产生的 UTC epoch milliseconds；
- `error.code`：下表中的稳定枚举；
- `error.retryable`：必须与 catalog 固定值一致；
- `error.retry_after_ms`：限流/暂时不可用时可给等待时长，否则为 null；
- `error.field`：可选 JSON Pointer 风格 wire 字段位置，不含用户值；
- `error.current_version`：版本冲突时返回 canonical version，否则为 null；
- `error.trace_id`：可选受控关联 ID，不是堆栈或内部错误文本。

服务端不得把异常类名、SQL、token、邀请码、资源 ID、用户输入或第三方 body 填入 error envelope。一个 HTTP status 可以承载多个稳定 code；客户端必须按 code 分支，不能只看 status。

## 稳定目录

| code | HTTP | retryable | 稳定语义 |
|---|---:|:---:|---|
| `UNAUTHENTICATED` | 401 | false | 缺少或无法验证 Bearer 凭据 |
| `TOKEN_EXPIRED` | 401 | true | access/refresh credential 已过期；允许先执行受控 refresh/重新登录 |
| `ACCOUNT_DISABLED` | 403 | false | canonical account 已禁用 |
| `FORBIDDEN` | 403 | false | 当前 membership actor 缺少 operation capability |
| `NOT_FOUND` | 404 | false | 资源不存在或对当前 actor 不可见；不得泄露两者差异 |
| `VALIDATION_FAILED` | 422 | false | wire/schema 或业务输入不合法 |
| `VERSION_CONFLICT` | 409 | false | expected version 与 canonical version 不同 |
| `IDEMPOTENCY_CONFLICT` | 409 | false | 同 lookup scope/key 对应不同 route/body fingerprint |
| `RATE_LIMITED` | 429 | true | 触发服务端限流 |
| `UNAVAILABLE` | 503 | true | 服务或真实外部依赖暂时不可用 |
| `TIMEOUT` | 504 | true | 服务端受控处理超时；不能据此推断 mutation 未提交 |
| `PROTOCOL_ERROR` | 502 | false | 服务端/上游无法满足冻结 wire 协议或主版本 |
| `INVITE_NOT_FOUND` | 404 | false | invite code 不存在 |
| `INVITE_EXPIRED` | 410 | false | invite 已过期 |
| `INVITE_USED` | 409 | false | single-use invite 已消费 |
| `INVITE_ROLE_MISMATCH` | 403 | false | invite target role 与 redeem operation 不匹配 |
| `INVITE_TARGET_DELETED` | 410 | false | invite 绑定的 Member/Admin 已删除 |
| `AUDIT_INCONSISTENT` | 409 | false | append-only 同 ID 出现不同 payload |
| `ASSET_REJECTED` | 422 | false | asset purpose、类型、大小、摘要或提交状态不满足策略 |
| `INSUFFICIENT_BALANCE` | 422 | false | canonical balance 不足以兑换 |
| `NOT_ASSIGNED` | 422 | false | Task/Reward 未分配给目标 Member |
| `COOLDOWN_ACTIVE` | 422 | false | Reward 仍处于该 Member 的独立冷却 |
| `BALANCE_OVERFLOW` | 422 | false | 任务完成或人工调整会越过服务端余额数值边界 |
| `COMPLETION_ALREADY_CANCELLED` | 409 | false | 另一个已提交命令已取消 Completion |
| `LAST_OWNER` | 409 | false | mutation 会移除 Circle 最后一个 Owner |
| `MEMBER_LIMIT_REACHED` | 422 | false | canonical entitlement 的成员额度已满 |
| `PERMANENTLY_UNAVAILABLE` | 422 | false | 一次性 Reward 已永久不可兑换 |
| `PHOTO_REQUIRED` | 422 | false | Task 要求有效、已 commit 的 proof asset |
| `PURCHASE_REJECTED` | 422 | false | Google Play canonical verification 拒绝 purchase proof |
| `TAG_IN_USE` | 409 | false | active Task 仍引用目标 Tag |

## 客户端处理规则

- `retryable=true` 只表示同一 operation 可以在前置条件满足后重试；mutation 必须复用原 `Idempotency-Key`、body、path、expected version 与 durable request intent。
- `TIMEOUT`、连接断开或响应解析失败都属于“提交结果未知”，必须幂等重放，不能先创建第二个 key。
- `TOKEN_EXPIRED` 先做一次受控 refresh；refresh 失败则清除远端 session 并重新认证。不得无限循环。
- `RATE_LIMITED` 优先遵守 `retry_after_ms`；没有值时使用有上限退避。`UNAVAILABLE`/`TIMEOUT` 同样需要退避和网络/生命周期约束。
- `VERSION_CONFLICT` 必须读取 `current_version`、pull 最新 aggregate 并让领域层决定重放/合并；不能只替换 expected version 后自动覆盖用户看不到的变化。
- `IDEMPOTENCY_CONFLICT`、`AUDIT_INCONSISTENT` 是一致性事故，禁止自动换 key 绕过。客户端保持 pending/rejected 证据并触发受控恢复或人工诊断。
- Invite、余额、assignment、cooldown、权限类非重试错误要呈现稳定脱敏 UI 结果；不得在客户端重新计算并声称服务端成功。
- malformed/不属于当前 snapshot 的 cursor 使用 operation 声明的 `VALIDATION_FAILED`，`field` 指向 `/query/cursor` 或 `/query/change_cursor`；客户端丢弃未投影 staging 页并从明确 bootstrap 入口重建，不能猜 cursor。所有 GET 均声明该错误并提供 fixture，用于 header/path/query 的 schema 或业务组合验证。

## 服务端映射规则

1. 在领域边界将内部异常转换成唯一稳定 code，并按 catalog 固定 status/retryable。
2. mutation 事务提交后即使 response serialization/transport 失败，也必须保留幂等记录；后续不能返回与首次 commit 不同的业务错误。
3. 外部 Google/Play/object-store 错误只映射到本 operation 允许的 code，不透传第三方 code/body。
4. validator 要求 catalog 中每个 code 至少有一个 operation owner 和一个 fixture owner，并要求高风险 operation 的声明错误均存在 fixture；缺 owner 或 fixture 即合同失败。
