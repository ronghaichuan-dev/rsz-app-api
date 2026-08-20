# Clearwave Backend Contract v1

本目录冻结 Clearwave 自建 Backend Contract v1。服务端、Live adapter、Demo handler 与 contract test 必须共同以
[`openapi/clearwave-backend-v1.json`](openapi/clearwave-backend-v1.json) 为唯一 machine-readable 权威合同；其余文档只解释该 JSON，不能覆盖其中的 path、method、字段、错误或 schema。

## 冻结状态

- OpenAPI：3.1.0，JSON Schema dialect：2020-12，合同版本：`1.0.0`，wire 主版本：`1`。
- 46 个唯一 operation、210 个具名 schema、30 个稳定错误码。
- 73 个任务 1 capability 被精确归属：47 个映射记录指向 46 个 operation，另有 11 个 `LOCAL_ONLY`、9 个 `EXTERNAL_SDK`、6 个明确 `NOT_REQUIRED`。
- 13 个 fixture 文件包含 365 个 case：90 个成功 case、275 个风险失败 case；每个 mutation 同时冻结首次提交与幂等重放，30 个稳定错误码均至少有一个 fixture owner。
- 本合同没有 `servers` 节点。真实 HTTPS base URL 是部署配置，不得用仓库 placeholder 代替。
- 当前客户端、HTTP adapter 与 Demo adapter 已按本合同预接线；没有可信真实环境时仍只能报告 `conditional_ready`，不能把 fixture/Demo 结果表述为真实互操作。

## 产物职责

| 产物                                                     | 职责                                                                           |
|--------------------------------------------------------|------------------------------------------------------------------------------|
| `openapi/clearwave-backend-v1.json`                    | 唯一 wire 权威：operation、header、schema、错误、分页、ID 映射与幂等规则                          |
| `postman/Clearwave-Backend-v1.postman_collection.json` | 可直接导入 Postman 的 46-operation 请求清单、合成示例与轻量断言；由 OpenAPI/fixture 机械生成           |
| `postman/API_FUNCTION_USAGE_GUIDE.md`                  | 逐接口解释用途、当前产品功能入口、触发时机与主要客户端 owner                                            |
| `coverage/backend-capabilities.json`                   | 任务 1 capability 到 operation/本地/外部/不需要 API 的全覆盖账本，以及未来 owner                  |
| `coverage/contract-operation-gates.json`               | 46 个 operation 到 success fixture、实际 Demo handler 与领域 contract test 的精确 owner |
| `coverage/demo-operation-coverage.json`                | 46 个 operation 的有状态 Demo handler/test owner 清单                               |
| `examples/*.json`                                      | 请求与响应 fixture；只使用合成值，不是可连接的环境凭据                                              |
| `authentication.md`                                    | account/session/Guest/membership 的信任与生命周期                                    |
| `consistency-and-sync.md`                              | outbox、版本、原子提交、投影、分页与恢复规则                                                    |
| `error-catalog.md`                                     | 稳定错误语义与客户端处理                                                                 |
| `security-and-redaction.md`                            | 敏感边界、日志、asset 与外部 proof 约束                                                   |
| `demo-scenarios.md`                                    | 与 Live 同 wire 的有状态 Demo 行为，不伪造外部 SDK 成功                                      |
| `server-implementation-checklist.md`                   | 服务端可以逐项执行的实现与验收清单                                                            |

## Wire 基线

- JSON 字段一律 `lower_snake_case`，未知字段因 `additionalProperties: false` 被拒绝。
- 时间点一律 UTC epoch milliseconds；日历日一律 `yyyy-MM-dd`；日历统计、重复与边界一律携带 IANA `zone_id`。
- 每次请求必带 `X-Request-Id`、`X-Contract-Version: 1`、`X-Client-Version`。除明确 public operation 外，通过 `Authorization: Bearer ...` 提交对应凭据。
- 所有 `POST`、`PUT`、`PATCH`、`DELETE` 必带 `Idempotency-Key`；可变 aggregate 同时提交 operation 冻结的 `expected_*_version` 字段。
- 成功 envelope 必含 `contract_version`、`request_id`、`server_time_ms`、强类型 `data`；领域 commit/审计 mutation 返回稳定 receipt，session/upload 等 mutation 也必须保证首次结果与幂等重放的 canonical parity，必要时包含 `change_cursor`。
- 错误 envelope 必含 `contract_version`、`request_id`、`server_time_ms`、`error.code` 与 `error.retryable`；其余错误上下文为显式 nullable 字段。
- 列表只使用 opaque `cursor`、`limit`、`next_cursor`、`has_more`、`snapshot_cursor`。`pullCircleBootstrapDelta` 以 `change_cursor` 作为输入；一个 server commit 永不跨页拆分。
- 业务图片与证明只引用已 commit 的 `asset_id`；绝不通过业务 operation 传 URI、路径或图片内容。
- `createInviteGuestSession` 必须签发独立 `guest_upgrade_grant`，供 `exchangeGoogleProof` 完成匿名升级；该 grant 与 token 同级保护。
- Circle delta 的每个 commit 必须完整携带 outgoing invite、notification outbox 与复合 occurrence tombstone；不能以伪造 occurrence ID 代替复合 key。

## API 族

| Tag | 数量 | operationId |
|---|---:|---|
| Auth | 5 | `exchangeGoogleProof`、`createInviteGuestSession`、`refreshSession`、`revokeSession`、`getCurrentAccount` |
| Circle | 12 | onboarding、Circle/selection、Member、Administrator、leave 的冻结 operation |
| Invite | 5 | create/refresh/revoke 与 Administrator/Member 两种 atomic redeem |
| Sync | 1 | `pullCircleBootstrapDelta` |
| Asset | 2 | `prepareAssetUpload`、`commitAssetUpload` |
| Task | 5 | Tag/Task upsert/delete 与 occurrence list |
| TaskLedger | 2 | `completeTask`、`cancelTaskCompletion` |
| Ledger | 3 | balance batch、ledger list、受权人工调整 |
| Reward | 3 | definition save/delete 与 eligibility |
| Redemption | 1 | `redeemReward` |
| ReadModels | 4 | Exchange list、Statistics、Compare、Completion details |
| Support | 1 | `submitFeedback` |
| Entitlement | 2 | Play purchase verification 与 entitlement snapshot |

Main 不新增 BFF，未来 UI 仍只读 Room projection；Exchange 没有 public create/update/delete；Statistics 取消入口只能复用 `cancelTaskCompletion`。诊断日志上传在 v1 安全退役，Device Push 因当前没有产品入口而不建接口。Google Credential Manager、Play Billing client、Firebase、Ads/UMP、TikTok、AlarmManager、本地通知、系统 Picker/分享仍是外部或本地边界。

## Postman

Postman 可直接导入 [`postman/Clearwave-Backend-v1.postman_collection.json`](postman/Clearwave-Backend-v1.postman_collection.json)。集合默认使用不可连接的 `https://staging.invalid`
，所有凭据变量为空；导入和幂等重放说明见 [`postman/README.md`](postman/README.md)，逐接口产品用途见 [`postman/API_FUNCTION_USAGE_GUIDE.md`](postman/API_FUNCTION_USAGE_GUIDE.md)。OpenAPI 仍是唯一接口事实源，Postman
文件不得反向覆盖合同。

## 校验

从仓库根目录执行：

```shell
PYTHONDONTWRITEBYTECODE=1 python3 ai/aw/validate_backend_contract.py
PYTHONDONTWRITEBYTECODE=1 python3 -m unittest -v ai.aw.tests.test_validate_backend_contract
PYTHONDONTWRITEBYTECODE=1 python3 -m unittest -v ai.aw.tests.test_validate_backend_release_artifact
PYTHONDONTWRITEBYTECODE=1 python3 backend-contract/postman/generate_collection.py --check
```

validator 会拒绝非法 JSON/重复 key、错误 OpenAPI 版本、重复 operationId、断裂或外部 `$ref`、占位 object、非 `lower_snake_case` 字段、缺失通用协议、单/双向错误映射缺口、幂等重放漂移、Guest grant 无签发者、sync change family 缺失、判别字段交叉混用、Task 时间规则漂移、fixture/schema 不一致、inventory/coverage 漏项及孤儿 operation；还会逐 operation 对比 Kotlin endpoint/stable error registry、Demo handler、success fixture、领域测试与 HTTP/Demo 参数化 suite，并扫描 supported Online 调用图的 raw path、页面 Demo 分支、硬编码 unavailable 与 production fake。

最终 `assembleLiveRelease` 后，使用 `validate_backend_release_artifact.py` 显式传入该 variant 的 APK、`mapping.txt` 与 `resources.txt`。Gate 只输出 basename、SHA-256、计数和去敏 marker id；它跳过 `META-INF` 签名 entry 内容，并拒绝最终 DEX/class/resource graph 中的 Demo transport/store/persona/seed/control、固定邀请码、Billing fake 开关/token 与 Backend example/localhost/loopback origin。

## 兼容性与变更控制

v1 冻结后，path、method、operationId、字段名、ID 规则、错误语义与原子 bundle 都是兼容性边界。发现真实不可满足项时，必须在同一变更中同步修改 OpenAPI、fixtures、coverage、validator、聚焦测试与兼容性说明，并重新执行自动 review；不得只改说明或某一侧实现。

初始冻结提交后的 standalone 恢复审查在 runtime 接线前补齐了 Guest upgrade grant 签发、完整 Circle commit projection、复合 occurrence tombstone、通知 outbox 版本与高风险判别/Task 时间约束。该闭环未改变任何 path、method 或 operationId；新增 required 字段/collection 属于 v1 尚未实现前的合同纠正，后续实现不得兼容旧的不完整形状。

当前仍需由真实环境提供但不属于合同猜测范围的配置只有：HTTPS base URL 与证书策略、Google OAuth issuer/audience、Google Play Developer API 服务凭据、实际 asset object store/签名器。缺少这些配置时只能报告未运行或 unavailable，不能声称真实互操作成功。
