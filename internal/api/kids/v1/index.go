package v1

import "github.com/gogf/gf/v2/frame/g"

type IndexReq struct {
	g.Meta `path:"/index" method:"get" tags:"A接口说明文档" summary:"A接口说明文档" description:"接口调用注意事项：1. 所有业务接口统一使用 /v1 版本前缀；2. 登录接口 POST /v1/kids/users/login 和健康检查 GET /v1/health 不需要 JWT，其他 kids 接口必须在 Authorization 请求头中携带 Bearer JWT；\n3. 接口返回统一格式为 code、message、data，业务是否成功以 code 为准；\n4. 需要多语言描述时在请求头 Language 中传 zh-CN 或 en，缺省默认 zh-CN；\n"`
}

type IndexRes struct {
}
