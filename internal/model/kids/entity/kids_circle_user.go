// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleUser is the golang structure for table kids_circle_user.
type KidsCircleUser struct {
	Id        uint64      `json:"id"        orm:"id"         description:"圈子用户关系ID"`          // 圈子用户关系ID
	CircleId  uint64      `json:"circleId"  orm:"circle_id"  description:"圈子ID"`              // 圈子ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"用户ID"`              // 用户ID
	Role      string      `json:"role"      orm:"role"       description:"圈子角色：admin/member"` // 圈子角色：admin/member
	MemberId  uint64      `json:"memberId"  orm:"member_id"  description:"绑定的家庭成员ID"`         // 绑定的家庭成员ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`              // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`              // 更新时间
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`              // 删除时间
	LeftAt    *gtime.Time `json:"leftAt"    orm:"left_at"    description:"退出时间"`              // 退出时间
}
