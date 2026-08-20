// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRequestDeduplication is the golang structure for table kids_request_deduplication.
type KidsRequestDeduplication struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:"主键"`       // 主键
	PrincipalScope       string      `json:"principalScope"       orm:"principal_scope"        description:"主体和成员作用域"` // 主体和成员作用域
	IdempotencyKey       string      `json:"idempotencyKey"       orm:"idempotency_key"        description:"幂等键"`      // 幂等键
	OperationId          string      `json:"operationId"          orm:"operation_id"           description:"接口操作标识"`   // 接口操作标识
	RouteFingerprint     string      `json:"routeFingerprint"     orm:"route_fingerprint"      description:"路由摘要"`     // 路由摘要
	BodyFingerprint      string      `json:"bodyFingerprint"      orm:"body_fingerprint"       description:"请求体摘要"`    // 请求体摘要
	ResponseStatus       uint        `json:"responseStatus"       orm:"response_status"        description:"首次响应状态码"`  // 首次响应状态码
	ResponseBody         string      `json:"responseBody"         orm:"response_body"          description:"首次规范响应"`   // 首次规范响应
	ResponseChangeCursor string      `json:"responseChangeCursor" orm:"response_change_cursor" description:"首次响应变更游标"` // 首次响应变更游标
	ResponseEtag         string      `json:"responseEtag"         orm:"response_etag"          description:"首次响应实体标签"` // 首次响应实体标签
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:"创建时间"`     // 创建时间
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:"更新时间"`     // 更新时间
}
