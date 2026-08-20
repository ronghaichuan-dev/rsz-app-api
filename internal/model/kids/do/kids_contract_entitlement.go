// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractEntitlement is the golang structure of table kids_contract_entitlement for DAO operations like Where/Data.
type KidsContractEntitlement struct {
	g.Meta        `orm:"table:kids_contract_entitlement, do:true"`
	Id            any         // 主键
	EntitlementId any         // 权益标识
	AccountId     any         // 账号标识
	ProductId     any         // 商品标识
	Status        any         // 权益状态
	ValidUntilAt  *gtime.Time // 有效截止时间
	VerifiedAt    *gtime.Time // 验证时间
	RevokedAt     *gtime.Time // 撤销时间
	Version       any         // 版本号
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
