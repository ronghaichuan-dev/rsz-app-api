// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUserToken is the golang structure for table kids_user_token.
type KidsUserToken struct {
	Id        uint64      `json:"id"        orm:"id"         description:"令牌ID"`     // 令牌ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"kids用户ID"` // kids用户ID
	Token     string      `json:"token"     orm:"token"      description:"访问令牌"`     // 访问令牌
	ExpiredAt *gtime.Time `json:"expiredAt" orm:"expired_at" description:"过期时间"`     // 过期时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`     // 创建时间
}
