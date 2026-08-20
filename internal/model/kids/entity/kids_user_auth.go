// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsUserAuth is the golang structure for table kids_user_auth.
type KidsUserAuth struct {
	Id        uint64      `json:"id"        orm:"id"         description:"授权记录ID"`             // 授权记录ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"kids用户ID"`           // kids用户ID
	Provider  string      `json:"provider"  orm:"provider"   description:"授权服务商：google/apple"` // 授权服务商：google/apple
	OpenId    string      `json:"openId"    orm:"open_id"    description:"服务商开放ID或主体标识"`       // 服务商开放ID或主体标识
	Email     string      `json:"email"     orm:"email"      description:"服务商邮箱"`              // 服务商邮箱
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`               // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`               // 更新时间
}
