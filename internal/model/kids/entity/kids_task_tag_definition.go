// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskTagDefinition is the golang structure for table kids_task_tag_definition.
type KidsTaskTagDefinition struct {
	Id        uint64      `json:"id"        orm:"id"          description:"主键"`       // 主键
	TaskTagId string      `json:"taskTagId" orm:"task_tag_id" description:"接口任务标签标识"` // 接口任务标签标识
	CircleId  string      `json:"circleId"  orm:"circle_id"   description:"接口圈子标识"`   // 接口圈子标识
	Name      string      `json:"name"      orm:"name"        description:"标签名称"`     // 标签名称
	Status    string      `json:"status"    orm:"status"      description:"标签状态"`     // 标签状态
	Version   uint64      `json:"version"   orm:"version"     description:"版本号"`      // 版本号
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at"  description:"删除时间"`     // 删除时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at"  description:"创建时间"`     // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at"  description:"更新时间"`     // 更新时间
}
