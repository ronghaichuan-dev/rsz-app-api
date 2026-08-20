// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskPreset is the golang structure for table kids_task_preset.
type KidsTaskPreset struct {
	Id          uint64      `json:"id"          orm:"id"          description:"任务预设ID"`   // 任务预设ID
	Title       string      `json:"title"       orm:"title"       description:"预设任务标题"`   // 预设任务标题
	Icon        string      `json:"icon"        orm:"icon"        description:"图标标识"`     // 图标标识
	Star        int         `json:"star"        orm:"star"        description:"默认星星数量"`   // 默认星星数量
	NeedPhoto   uint        `json:"needPhoto"   orm:"need_photo"  description:"是否建议照片凭证"` // 是否建议照片凭证
	Description string      `json:"description" orm:"description" description:"预设描述"`     // 预设描述
	Enabled     uint        `json:"enabled"     orm:"enabled"     description:"是否启用预设"`   // 是否启用预设
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`     // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`     // 更新时间
}
