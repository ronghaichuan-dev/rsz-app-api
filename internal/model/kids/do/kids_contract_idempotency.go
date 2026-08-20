// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractIdempotency is the golang structure of table kids_contract_idempotency for DAO operations like Where/Data.
type KidsContractIdempotency struct {
	g.Meta               `orm:"table:kids_contract_idempotency, do:true"`
	Id                   any         // 主键
	PrincipalScope       any         // 主体和成员作用域
	IdempotencyKey       any         // 幂等键
	OperationId          any         // 合同操作标识
	RouteFingerprint     any         // 路由摘要
	BodyFingerprint      any         // 请求体摘要
	ResponseStatus       any         // 首次响应状态码
	ResponseBody         any         // 首次规范响应
	ResponseChangeCursor any         // 首次响应变更游标
	ResponseEtag         any         // 首次响应实体标签
	CreatedAt            *gtime.Time // 创建时间
	UpdatedAt            *gtime.Time // 更新时间
}
