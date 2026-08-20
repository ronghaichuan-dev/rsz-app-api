// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleInfo is the golang structure for table kids_circle_info.
type KidsCircleInfo struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:"主键"`       // 主键
	CircleId             string      `json:"circleId"             orm:"circle_id"              description:"接口圈子标识"`   // 接口圈子标识
	Name                 string      `json:"name"                 orm:"name"                   description:"圈子名称"`     // 圈子名称
	Icon                 string      `json:"icon"                 orm:"icon"                   description:"圈子视觉引用"`   // 圈子视觉引用
	OwnerAdministratorId string      `json:"ownerAdministratorId" orm:"owner_administrator_id" description:"所有者管理员标识"` // 所有者管理员标识
	Status               string      `json:"status"               orm:"status"                 description:"圈子状态"`     // 圈子状态
	Version              uint64      `json:"version"              orm:"version"                description:"版本号"`      // 版本号
	DeletedAt            *gtime.Time `json:"deletedAt"            orm:"deleted_at"             description:"删除时间"`     // 删除时间
	CreatedAt            *gtime.Time `json:"createdAt"            orm:"created_at"             description:"创建时间"`     // 创建时间
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:"更新时间"`     // 更新时间
}
