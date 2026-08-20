// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsReward is the golang structure for table kids_reward.
type KidsReward struct {
	Id                 uint64      `json:"id"                 orm:"id"                   description:"奖励ID"`                                    // 奖励ID
	CircleId           uint64      `json:"circleId"           orm:"circle_id"            description:"所属圈子ID"`                                  // 所属圈子ID
	Title              string      `json:"title"              orm:"title"                description:"奖励标题"`                                    // 奖励标题
	Icon               string      `json:"icon"               orm:"icon"                 description:"图标标识"`                                    // 图标标识
	ImageUrl           string      `json:"imageUrl"           orm:"image_url"            description:"奖励图片地址"`                                  // 奖励图片地址
	StarCost           int         `json:"starCost"           orm:"star_cost"            description:"所需星星数量"`                                  // 所需星星数量
	Stock              int         `json:"stock"              orm:"stock"                description:"可用库存，-1表示不限量"`                            // 可用库存，-1表示不限量
	Description        string      `json:"description"        orm:"description"          description:"奖励描述"`                                    // 奖励描述
	RepeatRule         string      `json:"repeatRule"         orm:"repeat_rule"          description:"重复兑换规则：none/daily/weekly/monthly/custom"` // 重复兑换规则：none/daily/weekly/monthly/custom
	RepeatIntervalDays int         `json:"repeatIntervalDays" orm:"repeat_interval_days" description:"自定义重复兑换间隔天数"`                             // 自定义重复兑换间隔天数
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"           description:"创建时间"`                                    // 创建时间
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:"更新时间"`                                    // 更新时间
	DeletedAt          *gtime.Time `json:"deletedAt"          orm:"deleted_at"           description:"删除时间"`                                    // 删除时间
}
