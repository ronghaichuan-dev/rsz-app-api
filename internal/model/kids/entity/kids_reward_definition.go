// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardDefinition is the golang structure for table kids_reward_definition.
type KidsRewardDefinition struct {
	Id            uint64      `json:"id"            orm:"id"             description:"主键"`      // 主键
	RewardId      string      `json:"rewardId"      orm:"reward_id"      description:"接口奖励标识"`  // 接口奖励标识
	CircleId      string      `json:"circleId"      orm:"circle_id"      description:"接口圈子标识"`  // 接口圈子标识
	Title         string      `json:"title"         orm:"title"          description:"奖励标题"`    // 奖励标题
	Description   string      `json:"description"   orm:"description"    description:"奖励描述"`    // 奖励描述
	Visual        string      `json:"visual"        orm:"visual"         description:"奖励视觉引用"`  // 奖励视觉引用
	StarsRequired uint        `json:"starsRequired" orm:"stars_required" description:"兑换所需星星"`  // 兑换所需星星
	RepeatRule    string      `json:"repeatRule"    orm:"repeat_rule"    description:"重复兑换规则"`  // 重复兑换规则
	CooldownDays  uint        `json:"cooldownDays"  orm:"cooldown_days"  description:"自定义冷却天数"` // 自定义冷却天数
	ZoneId        string      `json:"zoneId"        orm:"zone_id"        description:"时区标识"`    // 时区标识
	Status        string      `json:"status"        orm:"status"         description:"奖励状态"`    // 奖励状态
	Version       uint64      `json:"version"       orm:"version"        description:"版本号"`     // 版本号
	DeletedAt     *gtime.Time `json:"deletedAt"     orm:"deleted_at"     description:"删除时间"`    // 删除时间
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`    // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"更新时间"`    // 更新时间
}
