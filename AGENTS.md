# 多应用程序微服务平台开发规范

本文档约束本项目中所有 GoFrame 微服务的开发方式。项目采用“一个应用程序对应一个微服务”的架构，并尽量保持 GoFrame CLI 默认目录规范，便于使用 `gf gen service` 等工具。

## 1. 架构原则

- 一个应用程序对应一个微服务，例如 `admin`、`kids`。
- 每个微服务必须独立启动、独立配置、独立端口、独立部署。
- 每个微服务必须使用独立数据库，不允许多个微服务直接共用同一个业务库。
- 每个微服务只提供自己应用程序需要的接口，不能把多个应用的接口混在同一个服务进程中。
- 微服务之间不能直接引用彼此的 `logic`，需要协作时通过 HTTP、gRPC、MQ 或专用 client 调用。
- 其他应用微服务需要使用 `admin` 的能力时，应调用 `admin` 暴露的 RPC/API 接口，不允许直连 `admin` 数据库。
- 公共能力只能放在 `internal/common`、`internal/middleware` 等公共目录，不能把具体业务逻辑放入公共层。
- 目录结构必须兼容 GoFrame CLI 默认约定，尤其是 `internal/logic` 和 `internal/service`。

## 2. 当前目录规范

```text
.
├── cmd/
│   ├── admin/main.go            # admin 微服务入口
│   └── kids/main.go             # kids 微服务入口
├── config/
│   ├── admin/
│   │   ├── config.dev.yaml
│   │   ├── config.test.yaml
│   │   └── config.prod.yaml
│   └── kids/
│       ├── config.dev.yaml
│       ├── config.test.yaml
│       └── config.prod.yaml
├── internal/
│   ├── api/                     # 请求/响应 DTO，按应用和版本划分
│   ├── controller/              # HTTP 控制器和路由注册
│   ├── logic/                   # 业务实现层，gf gen service 默认扫描源目录
│   ├── service/                 # 服务接口层，gf gen service 默认生成目标目录
│   ├── server/                  # 每个微服务的 Server 组装和启动逻辑
│   ├── dao/                     # gf gen dao 生成的 DAO
│   ├── model/                   # gf gen dao 生成的 DO/Entity
│   ├── common/                  # 公共能力
│   └── middleware/              # 公共中间件
├── manifest/                    # 部署相关文件
├── resource/                    # 静态资源或模板，纯 API 服务默认不依赖
├── go.mod
└── go.sum
```

新增应用程序微服务时，必须按以下目录补齐：

```text
cmd/<app>/main.go
config/<app>/config.dev.yaml
config/<app>/config.test.yaml
config/<app>/config.prod.yaml
internal/api/<app>/v1/
internal/controller/<app>/
internal/logic/<app>/
internal/server/<app>/cmd.go
internal/model/<app>/dao/        # gf gen dao 生成
internal/model/<app>/do/         # gf gen dao 生成
internal/model/<app>/entity/     # gf gen dao 生成
```

`internal/service` 由 `gf gen service` 维护，也可以在遵循生成风格的前提下手动调整。

## 3. 微服务启动规范

每个微服务入口只做一件事：调用对应 `internal/server/<app>` 中的 `Main.Run`。

示例：

```go
package main

import (
    "github.com/gogf/gf/v2/os/gctx"

    "rslytics-app-api/internal/server/admin"
)

func main() {
    admin.Main.Run(gctx.GetInitCtx())
}
```

服务组装逻辑必须放在：

```text
internal/server/<app>/cmd.go
```

服务组装职责包括：

- 解析运行环境。
- 加载对应环境配置文件。
- 创建 `g.Server("<app>")`。
- 注册公共中间件。
- 挂载当前应用的 controller 路由。
- 启动服务。

不得在 `cmd/<app>/main.go` 中编写业务逻辑、路由逻辑或配置解析逻辑。

## 4. 多环境配置规范

每个微服务必须提供三套配置：

```text
config/<app>/config.dev.yaml
config/<app>/config.test.yaml
config/<app>/config.prod.yaml
```

运行环境只允许：

```text
dev
test
prod
```

环境选择优先级：

```text
--env / -e > APP_ENV > GF_ENV > GO_ENV > dev
```

启动示例：

```bash
go run ./cmd/admin

go run ./cmd/admin --env test

go run ./cmd/kids -e prod

APP_ENV=prod go run ./cmd/kids
```

配置文件中不要配置不存在的静态目录，例如纯 API 服务不要配置：

```yaml
serverRoot: "resource/public"
```

除非确实创建并使用该静态资源目录。

## 5. 数据库隔离规范

本项目按微服务划分数据库连接和数据边界：

- 每个微服务必须配置自己的数据库连接。
- 每个微服务只读写自己的数据库。
- 不允许一个微服务直接访问另一个微服务的业务数据库。
- 不允许为了复用数据而跨应用 import 对方 `logic` 或直接调用对方 DAO。
- 跨微服务数据访问必须通过对方提供的 RPC/API/MQ 事件完成。
- 数据库配置必须写在当前微服务自己的配置文件中，例如 `config/admin/config.dev.yaml`。

推荐数据库配置结构：

```yaml
database:
  default:
    link: "mysql:root:password@tcp(127.0.0.1:3306)/rslytics_admin?loc=Local&parseTime=true"
    debug: true
```

不同微服务必须使用不同数据库或至少不同 schema，例如：

```text
admin -> rslytics_admin
kids  -> rslytics_kids
```

GoFrame ORM 默认使用 `g.DB()` 读取 `database.default`。如果一个微服务确实需要多个数据库连接，必须显式命名分组，并在代码中使用对应分组：

```yaml
database:
  default:
    link: "mysql:.../rslytics_kids"
  report:
    link: "mysql:.../rslytics_kids_report"
```

```go
g.DB("report")
```

多数据库分组只能用于当前微服务自己的数据边界，不能配置其他微服务的业务库。

## 6. RPC/API 服务间调用规范

`admin` 服务可能需要向其他应用微服务提供公共能力，例如用户、权限、组织、字典等。此类能力必须通过 RPC/API 暴露，不能通过共享数据库或共享 logic 实现。

规范：

- `admin` 对外提供的能力应定义清晰的 RPC/API 契约。
- 调用方只依赖契约和 client，不依赖 `internal/logic/admin`。
- `admin` 的数据库只允许 `admin` 服务自身访问。
- 跨服务调用必须设置超时、错误处理和日志。
- 高频调用应考虑缓存、批量接口或事件同步，避免服务间强耦合。
- 需要最终一致性的场景优先考虑 MQ 事件，而不是同步 RPC 强依赖。

建议目录预留：

```text
internal/rpc/                  # RPC 服务定义、注册或 client 封装
internal/client/               # 访问其他微服务的 client
internal/common/protocol/      # 多服务共享的协议常量或轻量 DTO
```

注意：如果协议 DTO 只属于某一个服务，应放在该服务自己的 `internal/api/<app>` 或 RPC 定义目录，不要随意放进 common。


## 7. GoFrame 分层职责

### 7.1 api 层

目录：

```text
internal/api/<app>/v1
```

职责：

- 定义请求 DTO。
- 定义响应 DTO。
- 使用 `g.Meta` 声明路由 path、method、tags、summary。
- 定义参数绑定、校验规则和文档描述。

示例：

```go
type UserGetReq struct {
    g.Meta `path:"/users/{id}" method:"get" tags:"Admin User" summary:"Get admin user"`
    Id     uint64 `p:"id" v:"required|min:1" dc:"User ID"`
}

type UserGetRes struct {
    Id       uint64 `json:"id" dc:"User ID"`
    Username string `json:"username" dc:"Username"`
}
```

规范：

- `api` 层不得调用 `service`、`logic`、数据库或外部服务。
- 接口请求/响应 DTO，以及业务入参/出参 Input/Output，都必须定义在 `internal/api/<app>/v1`。
- 请求结构体命名使用 `<Resource><Action>Req`。
- 响应结构体命名使用 `<Resource><Action>Res`。
- 所有对外字段必须写 `json` tag。
- 推荐补充 `dc` 描述，便于生成 OpenAPI 文档。

### 7.2 controller 层

目录：

```text
internal/controller/<app>
```

职责：

- 注册当前应用路由。
- 接收 `api` 层 DTO。
- 调用 `service` 层接口。
- 组装响应 DTO。

规范：

- `controller` 不写复杂业务逻辑。
- `controller` 不直接访问数据库、Redis、MQ。
- `controller` 不直接引用其他应用的 `logic`。
- 每个应用必须提供 `router.go`，并暴露：

```go
func Register(ctx context.Context, group *ghttp.RouterGroup)
```

- 路由版本默认按 `/v1` 分组：

```go
func Register(ctx context.Context, group *ghttp.RouterGroup) {
    group.Group("/v1", func(group *ghttp.RouterGroup) {
        group.Bind(
            health.New(),
            user.New(),
        )
    })
}
```

### 7.3 logic 层

目录：

```text
internal/logic/<app>
```

职责：

- 实现业务逻辑。
- 注册 service 接口实现。
- 编排 DAO、第三方 client、缓存、消息等基础能力。

规范：

- `logic` 是 `gf gen service` 的默认扫描源目录，不要移动。
- 业务实现应通过 `init` 注册到 `service`。
- 文件和类型命名要能表达业务含义。
- 不同应用的 logic 不能互相直接调用。
- 需要跨服务调用时，在公共 client 或 infra 层封装协议调用。

### 7.4 service 层

目录：

```text
internal/service
```

职责：

- 定义业务接口。
- 暴露获取接口实现的方法。
- 暴露注册接口实现的方法。

规范：

- `service` 是 `gf gen service` 的默认生成目标目录，不要移动。
- `service` 目录中禁止定义业务 Input/Output 结构体；此类结构体必须放在 `internal/api/<app>/v1`。
- 修改 `logic` 后优先运行：

```bash
gf gen service
```

- `controller` 只能依赖 `service`，不要直接依赖 `logic`。
- 手动调整 `service` 文件时，保持 GoFrame 生成风格。

### 7.5 server 层

目录：

```text
internal/server/<app>
```

职责：

- 定义 GoFrame `gcmd.Command`。
- 选择环境配置。
- 组装 HTTP Server。
- 注册全局中间件。
- 挂载当前应用路由。

规范：

- server 层不写业务逻辑。
- 每个微服务必须使用自己的 server package。
- `g.Server` 的名称应与应用名一致，例如 `g.Server("admin")`。

### 7.6 common 和 middleware 层

`internal/common` 放通用能力，例如：

- 统一响应结构。
- 通用错误处理。
- 通用 client 封装。
- 通用工具，不包含业务语义。

`internal/middleware` 放 HTTP 中间件，例如：

- 请求上下文。
- Trace ID。
- CORS。
- 鉴权入口。
- 日志。

规范：

- 公共层不能依赖具体应用的 `controller` 或 `logic`。
- 公共层可以被多个微服务复用。

## 8. 路由和接口规范

- 每个微服务根路径从 `/` 开始，不需要再加应用名前缀。
- API 版本统一放在 controller router 中，例如 `/v1`。
- RESTful 资源路径使用复数名词，例如 `/users/{id}`。
- 健康检查统一使用：

```text
GET /v1/health
```

- Swagger 地址由配置控制：

```yaml
server:
  openapiPath: "/api.json"
  swaggerPath: "/swagger"
```

## 9. 响应和错误规范

项目统一使用公共响应中间件：

```text
internal/common/response
```

统一响应格式：

```json
{
  "code": 0,
  "message": "OK",
  "data": {}
}
```

规范：

- controller 返回 `(res, err)`，不要手动 `WriteJson`，除非是文件下载、流式响应等特殊场景。
- 业务错误应返回 `error`，由统一响应中间件处理。
- 需要错误码时优先使用 GoFrame `gerror`、`gcode` 体系。

## 10. DAO/Model 生成规范

数据库 DAO、DO、Entity 必须使用 GoFrame CLI 生成，不允许手写基础 DAO 和数据库实体模型。

统一使用命令：

```bash
gf gen dao -c hack/config.yaml
```

生成配置文件：

```text
hack/config.yaml
```

本项目按微服务名作为 GoFrame DAO group，并把生成物集中放在 `internal/model/<service>` 下：

```text
internal/model/<service>/dao
internal/model/<service>/do
internal/model/<service>/entity
```

例如：

```text
internal/model/kids/dao
internal/model/kids/do
internal/model/kids/entity
```

说明：你期望的核心结构是 `dao` 和 `entity`，但 GoFrame 还会生成 `do` 包，建议保留：

- `dao`：数据访问对象。
- `do`：数据库操作对象。
- `entity`：表实体结构。

`hack/config.yaml` 中必须同时满足：

```yaml
path: "."
group: "kids" # 或 admin、merchant 等微服务名
daoPath: "internal/model/kids/dao"
doPath: "internal/model/kids/do"
entityPath: "internal/model/kids/entity"
```

关键规则：

- `path` 必须指向项目根目录 `.`。
- `group` 必须使用微服务名，例如 `admin`、`kids`。
- `daoPath`、`doPath`、`entityPath` 必须写成相对项目根目录的路径。
- 不要把 `path` 配成 `internal/dao/<service>` 或 `internal/model/<service>`。
- 如果把 `path` 配错，会生成类似 `internal/dao/kids/internal/dao/kids` 的错误深层目录。
- 每个微服务只能从自己的数据库生成 DAO/Model。
- `admin` 数据库只生成到 `internal/model/admin`。
- `kids` 数据库只生成到 `internal/model/kids`。
- 不允许 `kids` 代码使用 `admin` 分组下的 DAO/Model 访问 admin 数据。
- 不允许 `admin` 代码使用 `kids` 分组下的 DAO/Model 访问 kids 数据。
- 生成前必须确认 `hack/config.yaml` 中的数据库连接指向正确环境。
- 生成前必须确认 `tables` 和 `removePrefix` 配置是否符合预期，避免生成无关表。
- 生成后的 DAO/DO/Entity 不要手工改业务逻辑；业务逻辑写在 `internal/logic/<app>`。
- 如果确实需要调整生成结果，应优先调整数据库表结构、字段注释或 `gf gen dao` 配置后重新生成。
- 数据库表注释和字段注释必须统一使用中文，禁止新增英文 COMMENT。
- SQL 文件中的说明性注释也应使用中文，便于团队和生成后的模型描述保持一致。

推荐表命名：

```text
admin 数据库：admin_user、admin_role、admin_permission
kids 数据库：kids_user、kids_task、kids_reward
```

如果使用表前缀，建议在 `hack/config.yaml` 中配置 `removePrefix`，使生成后的 Go 类型更简洁。

## 11. GoFrame CLI 使用规范

常用命令：

```bash
gf gen service

gf gen dao -c hack/config.yaml
```

执行前确认目录存在：

```text
internal/logic
internal/service
```

执行后必须检查：

- `internal/service` 是否符合预期。
- controller 对 service 的调用是否需要同步调整。
- 是否有未注册的 logic 实现。

使用 `gf gen dao` 必须通过 `hack/config.yaml`，并确认生成路径不会覆盖非本服务代码。使用 `gf gen ctrl` 等其他命令前，也必须先确认生成路径不会覆盖已有人工代码。

## 12. 编码规范

- Go 代码必须执行 `gofmt`。
- 包名使用小写，不使用下划线。
- 文件名使用小写，可用下划线分隔，例如 `admin_user.go`。
- 导出类型、方法、接口使用清晰业务命名。
- `context.Context` 必须作为业务方法第一个参数。
- 登录、注册、业务操作等方法的 Input/Output 结构体必须定义在 `internal/api/<app>/v1`，不能定义在 `internal/service`。
- 不要在业务逻辑中使用全局可变状态，service 注册对象除外。
- 不要在 controller 中写复杂业务编排。
- 不要在 logic 中直接读取 HTTP request 对象。
- 不要跨应用直接 import 对方 `logic`。
- 新增依赖前确认是否已有 GoFrame 内置能力可用。
- 项目中不能出现硬编码

## 13. 新增接口流程

以给 `kids` 增加登录接口为例：

1. 在 `internal/api/kids/v1` 新增请求/响应 DTO 以及业务 Input/Output。
2. 在 `internal/logic/kids` 实现登录业务。
3. 运行 `gf gen service` 更新 `internal/service`。
4. 在 `internal/controller/kids` 新增 controller。
5. 在 `internal/controller/kids/router.go` 绑定 controller。
6. 运行 `gofmt`。
7. 运行 `go test ./...`。
8. 启动对应服务验证接口。

## 14. 新增微服务流程

以新增 `merchant` 应用为例：

1. 创建 `cmd/merchant/main.go`。
2. 创建 `config/merchant/config.dev.yaml`。
3. 创建 `config/merchant/config.test.yaml`。
4. 创建 `config/merchant/config.prod.yaml`。
5. 创建 `internal/api/merchant/v1`。
6. 创建 `internal/controller/merchant/router.go`。
7. 创建 `internal/logic/merchant`。
8. 创建 `internal/server/merchant/cmd.go`。
9. 在 `hack/config.yaml` 增加 `merchant` 的 `gfcli.gen.dao` 配置。
10. 运行 `gf gen dao -c hack/config.yaml` 生成 DAO/Model。
11. 运行 `gf gen service`。
12. 运行 `gofmt -w cmd internal`。
13. 运行 `go test ./...`。

## 15. 提交前检查

提交或交付前必须至少执行：

```bash
gofmt -w cmd internal

go test ./...
```

如果修改了 `logic`，需要确认是否应该运行：

```bash
gf gen service
```

如果修改了数据库表结构或 DAO 配置，需要运行：

```bash
gf gen dao -c hack/config.yaml
```

## 16. kids 登录规范

kids 微服务登录必须支持三种方式：

```text
guest
google
apple
```

接口入口：

```text
POST /v1/kids/users/login
```

游客登录规则：

- 游客登录必须提交 `deviceId`。
- 同一个设备重复游客登录，应返回同一个游客账号。
- 游客账号只能存在于 kids 微服务自己的数据库中。

授权登录规则：

- Google 和 Apple 登录必须提交服务端可校验的 `identityToken`。
- 服务端必须根据 Google/Apple 后台公钥校验 `identityToken` 的签名、签发方、客户端 ID 和过期时间。
- 服务端禁止信任客户端直接传入的 `openId`，只能使用校验成功后的 token `sub` 作为 openId。
- 校验成功后，以 provider + sub 作为唯一授权身份。

游客绑定授权账号规则：

- 如果同一个设备第一次以游客模式登录，系统创建游客账号。
- 如果同一个设备第二次以 Google/Apple 授权登录，且该设备存在未绑定的游客账号，则应把该游客账号绑定到授权身份。
- 绑定后继续使用原游客账号的 userId，不新建用户。
- 绑定成功后游客账号应升级为正式授权账号。
- 如果授权身份已存在，则直接返回授权身份对应账号，不再绑定当前设备游客账号。

数据库表建议见：

```text
manifest/sql/kids_auth.sql
```

生成 DAO/Model 前先创建表，然后执行：

```bash
gf gen dao -c hack/config.yaml
```


## 17. 持久化和注释规范

- 所有业务数据必须持久化到当前微服务自己的数据库中，禁止使用内存 `map`、全局变量、进程缓存作为业务数据的最终存储。
- 内存缓存只能用于可丢弃的性能优化，不能作为任务、奖励、家庭成员、星星流水、登录身份等核心业务数据的事实来源。
- 所有新增、修改、删除业务数据的 logic 方法必须通过 DAO 或 GoFrame ORM 操作数据库，并根据一致性要求使用事务。
- 多表写入必须使用事务，例如任务和任务分配、任务完成和星星流水、奖励兑换和星星扣减、授权登录和用户绑定。
- 查询接口也必须从数据库读取实时数据，不能从启动时初始化的内存假数据返回。
- 原型数据、演示数据必须通过 SQL migration/seed 脚本写入数据库，不能在业务逻辑中硬编码初始化到内存。
- 每个函数都必须增加中文注释，注释应说明函数职责、关键入参或关键业务规则。
- Go 代码注释统一使用中文；禁止新增英文函数注释，除非引用第三方生成代码且不手工维护。
- 导出函数、导出类型、导出方法必须符合 Go 注释习惯，以标识符开头并使用中文说明。
- 非导出函数也必须写中文注释，重点说明业务意图，避免无意义注释。
- 修改已有函数时，如果缺少注释，必须同步补齐中文注释。

## 18. 工具函数和常量目录规范

- 可复用常量必须统一定义在 `internal/consts` 目录，禁止在业务文件中散落重复字符串常量。
- 常用工具函数必须统一定义在 `internal/utils` 目录，禁止在具体业务 logic 文件中堆放通用函数。
- `internal/consts` 只能定义跨文件复用的常量，不得包含业务流程逻辑、数据库访问或外部服务调用。
- `internal/utils` 只能放无业务状态的通用工具函数，不得持有业务数据，不得绕过当前微服务数据边界访问其他微服务数据库。
- 与某个具体业务强绑定且不复用的小型私有辅助函数，可以保留在当前 logic 文件中，但仍必须添加中文注释。
- 新增表名、数据库分组名、时间格式、固定枚举等常量时，优先放入 `internal/consts`。
- 日期格式化、布尔转换、字符串默认值、token 随机串等通用函数时，优先放入 `internal/utils`。

## 19. JWT 鉴权和邀请码规范

- kids 登录接口必须签发 JWT Token，禁止再使用普通随机字符串作为访问令牌。
- JWT Token 必须写入访问令牌表，鉴权中间件必须同时校验 JWT 签名和数据库令牌有效性。
- JWT 必须包含 `userId`、`iat`、`exp` 等基础声明，并使用服务端配置的 `auth.jwt.secret` 进行签名。
- 除 `GET /v1/health` 和 `POST /v1/kids/users/login` 外，kids 微服务其他接口默认都必须校验 JWT。
- 鉴权中间件校验通过后，必须把 `userId` 写入请求上下文，业务逻辑优先从上下文获取当前用户身份。
- 创建圈子、生成邀请码、使用邀请码加入圈子等接口都需要登录后携带 JWT。
- 邀请码必须持久化到 kids 独立数据库，必须有过期时间、邀请角色、邀请人、使用人和使用时间。
- 邀请加入角色等请求枚举参数必须优先在 API DTO 中使用 GoFrame 自带 `v:"required|in:..."` 标签校验；controller 绑定已完成校验后，logic 不要重复校验同一枚举范围，更禁止手写多个字符串比较判断枚举范围。
- 邀请码统一为 6 位数字码；过期或已使用的邀请码不能再次加入。
- 只有圈子管理员可以生成管理员邀请码或成员邀请码。
- 管理员邀请码用于邀请其他管理员加入圈子；成员邀请码用于邀请成员加入圈子，可选择绑定到指定家庭成员。

## 20. kids 多语言规范

- kids 微服务接口返回描述必须支持多语言，优先从请求头 `Language` 读取语言标识。
- `Language` 缺失时可兼容读取 `Accept-Language`，仍缺失时默认使用 `zh-CN`。
- 当前支持语言统一定义在 `internal/consts/i18n.go`，新增语言时必须同步补充常量和语言文件。
- 多语言文件统一放在 `manifest/i18n/<language>/` 目录，例如 `manifest/i18n/zh-CN/common.yaml`。
- 请求语言解析、翻译调用和错误文案处理必须通过 `internal/common/i18n`，语言格式标准化等通用函数必须放在 `internal/utils`。
- kids 服务必须注册多语言中间件，在 JWT 鉴权和统一响应之前把语言写入请求上下文。
- 统一响应中的 `message` 必须根据请求上下文语言翻译，包括成功、鉴权失败、业务错误和框架错误码。
- 业务逻辑中需要持久化展示给用户的通知标题、通知内容、星星流水标题等，也必须按当前请求语言生成后写入数据库。
- GoFrame 内置参数校验必须通过 `manifest/i18n/<language>/validation.yaml` 维护 `gf.gvalid.rule.*` 多语言模板，新增或使用新的校验规则时必须同步补充 `zh-CN` 和 `en` 文案。
- API DTO 的 `v` 标签可使用稳定英文错误 key 作为自定义错误消息，key 必须同步写入对应语言文件；controller 禁止手写校验多语言分支。
- 统一响应处理中必须识别 GoFrame `gvalid.Error` 并优先返回首个校验错误，避免多个校验错误拼接后无法稳定翻译。
- 新增业务错误时，错误原文应作为稳定翻译 key，并同步补充 `zh-CN` 和 `en` 语言文件。
- 禁止在 controller 中手写多语言分支；controller 只调用 service 并返回 DTO。
