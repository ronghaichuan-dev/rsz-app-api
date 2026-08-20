// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskTag is the golang structure for table kids_task_tag.
type KidsTaskTag struct {
	Id        uint64      `json:"id"        orm:"id"         description:"任务标签ID"` // 任务标签ID
	Name      string      `json:"name"      orm:"name"       description:"标签名称"`   // 标签名称
	Color     string      `json:"color"     orm:"color"      description:"标签颜色"`   // 标签颜色
	SortOrder int         `json:"sortOrder" orm:"sort_order" description:"排序值"`    // 排序值
	DeletedAt *gtime.Time `json:"deletedAt" orm:"deleted_at" description:"删除时间"`   // 删除时间
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
}
