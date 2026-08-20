# Clearwave Backend Postman Collection

`Clearwave-Backend-v1.postman_collection.json` 是可直接导入 Postman 的 Collection v2.1，按 OpenAPI 的 13 个 Tag 分组并覆盖全部 46 个
operation。集合由冻结 OpenAPI 与 `examples/*.json` 的权威成功 fixture 机械生成，不是另一份接口事实源。

逐接口的产品用途、当前功能入口与调用时机见 [`API_FUNCTION_USAGE_GUIDE.md`](API_FUNCTION_USAGE_GUIDE.md)。

## 导入与配置

1. 在 Postman 选择 **Import**，导入 `Clearwave-Backend-v1.postman_collection.json`。
2. 新建本机 Postman Environment，并覆盖集合中的 `baseUrl`。默认值为 `https://staging.invalid`，用于防止误连真实服务。
3. 按场景只在本机 Environment 的 Current Value 中填写 `accessToken`、`refreshToken`、`inviteAccessToken`、`googleIdToken`、
   `guestUpgradeGrant`、`inviteCode` 或 `purchaseToken`；不得把真实值写回集合或提交到仓库。
4. path 参数默认使用合成 fixture ID；连接真实 staging 时按当前 seed/响应更新对应变量。

每次请求自动生成 `X-Request-Id`。mutation 的 `Idempotency-Key` 按 operation 保存，因此连续点击 **Send** 会重放同一逻辑请求；修改请求体并准备发起新命令时，先删除集合变量 `_idempotencyKey_<operationId>`。

集合附带每个 operation 的一个成功响应示例和最小 envelope 断言。完整的成功、错误、权限、版本冲突和幂等 fixture 仍以
`backend-contract/examples/*.json` 为准。

## 生成与校验

从仓库根目录执行：

```shell
PYTHONDONTWRITEBYTECODE=1 python3 backend-contract/postman/generate_collection.py
PYTHONDONTWRITEBYTECODE=1 python3 backend-contract/postman/generate_collection.py --check
```

OpenAPI 或 fixture 发生受控变更时必须重新生成；`--check` 会拒绝缺失 operation、孤儿 fixture、method/path/parameter 漂移及过期集合。
