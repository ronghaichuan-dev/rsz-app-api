// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsFamilyMember is the golang structure for table kids_family_member.
type KidsFamilyMember struct {
	Id          uint64      `json:"id"          orm:"id"           description:"家庭成员ID"`         // 家庭成员ID
	CircleId    uint64      `json:"circleId"    orm:"circle_id"    description:"所属圈子ID"`         // 所属圈子ID
	Name        string      `json:"name"        orm:"name"         description:"显示名称"`           // 显示名称
	Gender      string      `json:"gender"      orm:"gender"       description:"性别：male/female"` // 性别：male/female
	Avatar      string      `json:"avatar"      orm:"avatar"       description:"头像地址或预设标识"`      // 头像地址或预设标识
	AvatarStyle string      `json:"avatarStyle" orm:"avatar_style" description:"虚拟形象风格标识"`       // 虚拟形象风格标识
	Relation    string      `json:"relation"    orm:"relation"     description:"家庭关系"`           // 家庭关系
	Owner       uint        `json:"owner"       orm:"owner"        description:"是否家庭拥有者"`        // 是否家庭拥有者
	BindUserId  uint64      `json:"bindUserId"  orm:"bind_user_id" description:"绑定用户ID"`         // 绑定用户ID
	BoundAt     *gtime.Time `json:"boundAt"     orm:"bound_at"     description:"绑定时间"`           // 绑定时间
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`           // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"   description:"更新时间"`           // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"   description:"删除时间"`           // 删除时间
}
