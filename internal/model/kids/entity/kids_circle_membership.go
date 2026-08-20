// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleMembership is the golang structure for table kids_circle_membership.
type KidsCircleMembership struct {
	Id           uint64      `json:"id"           orm:"id"            description:"主键"`       // 主键
	MembershipId string      `json:"membershipId" orm:"membership_id" description:"接口成员身份标识"` // 接口成员身份标识
	CircleId     string      `json:"circleId"     orm:"circle_id"     description:"接口圈子标识"`   // 接口圈子标识
	AccountId    string      `json:"accountId"    orm:"account_id"    description:"接口账号标识"`   // 接口账号标识
	ActorType    string      `json:"actorType"    orm:"actor_type"    description:"授权主体类型"`   // 授权主体类型
	ActorId      string      `json:"actorId"      orm:"actor_id"      description:"授权主体标识"`   // 授权主体标识
	Role         string      `json:"role"         orm:"role"          description:"成员角色"`     // 成员角色
	Permissions  string      `json:"permissions"  orm:"permissions"   description:"权限集合"`     // 权限集合
	Status       string      `json:"status"       orm:"status"        description:"成员身份状态"`   // 成员身份状态
	Version      uint64      `json:"version"      orm:"version"       description:"版本号"`      // 版本号
	DeletedAt    *gtime.Time `json:"deletedAt"    orm:"deleted_at"    description:"删除时间"`     // 删除时间
	CreatedAt    *gtime.Time `json:"createdAt"    orm:"created_at"    description:"创建时间"`     // 创建时间
	UpdatedAt    *gtime.Time `json:"updatedAt"    orm:"updated_at"    description:"更新时间"`     // 更新时间
}
