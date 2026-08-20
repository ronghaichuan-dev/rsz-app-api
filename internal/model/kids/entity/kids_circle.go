// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircle is the golang structure for table kids_circle.
type KidsCircle struct {
	Id          uint64      `json:"id"          orm:"id"            description:"圈子ID"`    // 圈子ID
	Name        string      `json:"name"        orm:"name"          description:"圈子名称"`    // 圈子名称
	Icon        string      `json:"icon"        orm:"icon"          description:"圈子图标标识"`  // 圈子图标标识
	OwnerUserId uint64      `json:"ownerUserId" orm:"owner_user_id" description:"创建者用户ID"` // 创建者用户ID
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"    description:"创建时间"`    // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"    description:"更新时间"`    // 更新时间
	DeletedAt   *gtime.Time `json:"deletedAt"   orm:"deleted_at"    description:"删除时间"`    // 删除时间
}
