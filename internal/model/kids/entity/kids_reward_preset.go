// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardPreset is the golang structure for table kids_reward_preset.
type KidsRewardPreset struct {
	Id          uint64      `json:"id"          orm:"id"          description:"奖励预设ID"`   // 奖励预设ID
	Title       string      `json:"title"       orm:"title"       description:"预设奖励标题"`   // 预设奖励标题
	Icon        string      `json:"icon"        orm:"icon"        description:"图标标识"`     // 图标标识
	ImageUrl    string      `json:"imageUrl"    orm:"image_url"   description:"奖励图片地址"`   // 奖励图片地址
	StarCost    int         `json:"starCost"    orm:"star_cost"   description:"默认所需星星数量"` // 默认所需星星数量
	Description string      `json:"description" orm:"description" description:"预设描述"`     // 预设描述
	Enabled     uint        `json:"enabled"     orm:"enabled"     description:"是否启用预设"`   // 是否启用预设
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"  description:"创建时间"`     // 创建时间
	UpdatedAt   *gtime.Time `json:"updatedAt"   orm:"updated_at"  description:"更新时间"`     // 更新时间
}
