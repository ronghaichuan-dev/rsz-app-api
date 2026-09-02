# rslytics-app-api

基于 GoFrame v2 的多应用程序微服务基础架构。这里的多个 APP 指多个独立应用程序，例如管理后台、儿童端、小程序、商家端等；每个应用程序由一个独立微服务提供接口。

## 架构定位

- 一个应用程序对应一个微服务：例如 `admin` 应用对应 `admin` 服务，`kids` 应用对应 `kids` 服务。
- 每个微服务独立启动、独立端口、独立配置、独立部署。
- 每个微服务只暴露该应用程序需要的接口，不把多个应用的接口混在同一个进程里。
- 多个微服务可以共享通用能力，例如统一响应、中间件、基础设施 client、工具包等。
- 微服务之间不直接引用彼此的 `logic`，如需协作通过 HTTP/gRPC/MQ/RPC client 等方式通信。
- 目录保持 GoFrame CLI 默认规范，`gf gen service` 可以从 `internal/logic` 扫描并生成/更新 `internal/service`。

## 示例映射

| 应用程序 | 微服务 | 启动入口 | 默认端口 | 业务目录 |
| --- | --- | --- | --- | --- |
| 管理后台 Admin App | admin | `cmd/admin` | `:8001` | `internal/controller/admin`, `internal/logic/admin` |
| 儿童端 Kids App | kids | `cmd/kids` | `:8002` | `internal/controller/kids`, `internal/logic/kids` |

## 目录结构

```text
.
├── cmd/
│   ├── admin/main.go            # 管理后台应用的微服务入口
│   └── kids/main.go             # 儿童端应用的微服务入口
├── config/
│   ├── admin/                   # admin 服务配置
│   │   ├── config.dev.yaml
│   │   ├── config.test.yaml
│   │   └── config.prod.yaml
│   └── kids/                    # kids 服务配置
│       ├── config.dev.yaml
│       ├── config.test.yaml
│       └── config.prod.yaml
├── internal/
│   ├── api/
│   │   ├── admin/v1/            # admin 请求/响应 DTO
│   │   └── kids/v1/             # kids 请求/响应 DTO
│   ├── controller/
│   │   ├── admin/               # admin HTTP 控制器与路由注册
│   │   └── kids/                # kids HTTP 控制器与路由注册
│   ├── logic/
│   │   ├── admin/               # admin 业务实现
│   │   └── kids/                # kids 业务实现
│   ├── service/                 # GoFrame service 接口层
│   ├── common/response/         # 多个微服务共享的统一响应
│   ├── middleware/              # 多个微服务共享的中间件
│   └── server/
│       ├── admin/cmd.go         # admin 的 GoFrame Server 组装
│       └── kids/cmd.go          # kids 的 GoFrame Server 组装
├── manifest/                    # 部署、镜像、脚本等清单
├── resource/                    # 静态资源、模板等
├── go.mod
└── go.sum
```

## 运行

安装依赖：

```bash
go mod tidy
```

启动管理后台应用微服务：

```bash
go run ./cmd/admin
```

启动儿童端应用微服务：

```bash
go run ./cmd/kids
```

默认使用 `dev` 配置，也可以通过命令行参数切换环境：

```bash
go run ./cmd/admin --env dev
go run ./cmd/admin --env test
go run ./cmd/admin --env prod
```

短参数也支持：

```bash
go run ./cmd/kids -e test
```

也可以通过环境变量切换，优先级为 `--env/-e` > `APP_ENV` > `GF_ENV` > `GO_ENV` > `dev`：

```bash
APP_ENV=prod go run ./cmd/kids
```

## 示例接口

`admin` 默认监听 `:8001`：

- `GET http://127.0.0.1:8001/v1/health`
- `GET http://127.0.0.1:8001/v1/users/1`
- `GET http://127.0.0.1:8001/swagger`

`kids` 默认监听 `:8002`：

- `GET http://127.0.0.1:8002/v1/health`
- `GET http://127.0.0.1:8002/v1/kids/profile`
- `POST http://127.0.0.1:8002/v1/kids/users/login`
- `GET http://127.0.0.1:8002/v1/kids/tasks`
- `GET http://127.0.0.1:8002/v1/kids/rewards`
- `GET http://127.0.0.1:8002/swagger`
- `GET http://127.0.0.1:8002/api.json`（Swagger/OpenAPI JSON 文档）

统一响应格式：

```json
{
  "code": 0,
  "message": "OK",
  "data": {}
}
```

## 使用 GoFrame CLI

生成或更新 service 接口：

```bash
gf gen service
```

本项目已保留 GoFrame 默认扫描目录：

- 扫描源目录：`internal/logic`
- 生成目标目录：`internal/service`

## 新增一个应用程序的微服务

以新增商家端应用 `merchant` 为例：

1. 创建启动入口：`cmd/merchant/main.go`。
2. 创建三套服务配置：`config/merchant/config.dev.yaml`、`config/merchant/config.test.yaml`、`config/merchant/config.prod.yaml`。
3. 创建服务组装：`internal/server/merchant/cmd.go`。
4. 创建协议定义：`internal/api/merchant/v1`。
5. 创建控制器：`internal/controller/merchant`。
6. 创建业务实现：`internal/logic/merchant`。
7. 运行 `gf gen service` 更新 `internal/service`。
8. 如需访问其他微服务，封装独立 client，不跨应用直接引用对方 `logic`。

## 多环境配置

每个微服务都有独立的三套配置：

```text
config/<service>/config.dev.yaml
config/<service>/config.test.yaml
config/<service>/config.prod.yaml
```

环境选择规则：

1. 命令行参数：`--env` 或 `-e`
2. 环境变量：`APP_ENV`
3. 环境变量：`GF_ENV`
4. 环境变量：`GO_ENV`
5. 默认值：`dev`

示例：

```bash
go run ./cmd/admin -e dev
go run ./cmd/admin -e test
go run ./cmd/admin -e prod
APP_ENV=test go run ./cmd/kids
```
