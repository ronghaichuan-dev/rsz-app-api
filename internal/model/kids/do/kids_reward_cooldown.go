// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardCooldown is the golang structure of table kids_reward_cooldown for DAO operations like Where/Data.
type KidsRewardCooldown struct {
	g.Meta                 `orm:"table:kids_reward_cooldown, do:true"`
	Id                     any         // 主键
	RewardId               any         // 接口奖励标识
	MemberId               any         // 接口成员标识
	CooldownUntilAt        *gtime.Time // 冷却结束时间
	LastRedeemedAt         *gtime.Time // 最近兑换时间
	PermanentlyUnavailable any         // 是否永久不可兑换
	Version                any         // 版本号
	CreatedAt              *gtime.Time // 创建时间
	UpdatedAt              *gtime.Time // 更新时间
}
