// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardCooldown is the golang structure for table kids_reward_cooldown.
type KidsRewardCooldown struct {
	Id                     uint64      `json:"id"                     orm:"id"                      description:"主键"`       // 主键
	RewardId               string      `json:"rewardId"               orm:"reward_id"               description:"接口奖励标识"`   // 接口奖励标识
	MemberId               string      `json:"memberId"               orm:"member_id"               description:"接口成员标识"`   // 接口成员标识
	CooldownUntilAt        *gtime.Time `json:"cooldownUntilAt"        orm:"cooldown_until_at"       description:"冷却结束时间"`   // 冷却结束时间
	LastRedeemedAt         *gtime.Time `json:"lastRedeemedAt"         orm:"last_redeemed_at"        description:"最近兑换时间"`   // 最近兑换时间
	PermanentlyUnavailable uint        `json:"permanentlyUnavailable" orm:"permanently_unavailable" description:"是否永久不可兑换"` // 是否永久不可兑换
	Version                uint64      `json:"version"                orm:"version"                 description:"版本号"`      // 版本号
	CreatedAt              *gtime.Time `json:"createdAt"              orm:"created_at"              description:"创建时间"`     // 创建时间
	UpdatedAt              *gtime.Time `json:"updatedAt"              orm:"updated_at"              description:"更新时间"`     // 更新时间
}
