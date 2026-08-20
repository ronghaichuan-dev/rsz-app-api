// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsEntitlement is the golang structure for table kids_entitlement.
type KidsEntitlement struct {
	Id            uint64      `json:"id"            orm:"id"             description:"主键"`     // 主键
	EntitlementId string      `json:"entitlementId" orm:"entitlement_id" description:"权益标识"`   // 权益标识
	AccountId     string      `json:"accountId"     orm:"account_id"     description:"账号标识"`   // 账号标识
	ProductId     string      `json:"productId"     orm:"product_id"     description:"商品标识"`   // 商品标识
	Status        string      `json:"status"        orm:"status"         description:"权益状态"`   // 权益状态
	ValidUntilAt  *gtime.Time `json:"validUntilAt"  orm:"valid_until_at" description:"有效截止时间"` // 有效截止时间
	VerifiedAt    *gtime.Time `json:"verifiedAt"    orm:"verified_at"    description:"验证时间"`   // 验证时间
	RevokedAt     *gtime.Time `json:"revokedAt"     orm:"revoked_at"     description:"撤销时间"`   // 撤销时间
	Version       uint64      `json:"version"       orm:"version"        description:"版本号"`    // 版本号
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`   // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"更新时间"`   // 更新时间
}
