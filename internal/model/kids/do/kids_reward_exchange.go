// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardExchange is the golang structure of table kids_reward_exchange for DAO operations like Where/Data.
type KidsRewardExchange struct {
	g.Meta                         `orm:"table:kids_reward_exchange, do:true"`
	Id                             any         // 主键
	ExchangeId                     any         // 接口兑换标识
	CircleId                       any         // 接口圈子标识
	MemberId                       any         // 接口成员标识
	MemberNameSnapshot             any         // 成员名称快照
	MemberAvatarSnapshot           any         // 成员头像快照
	RewardId                       any         // 接口奖励标识
	RewardTitleSnapshot            any         // 奖励标题快照
	RewardVisualSnapshot           any         // 奖励视觉快照
	StarsDeductedSnapshot          any         // 扣减星星快照
	RewardRepeatRuleSnapshot       any         // 重复规则快照
	RewardCooldownDaysSnapshot     any         // 冷却天数快照
	CooldownUntilAtSnapshot        *gtime.Time // 冷却结束快照
	PermanentlyUnavailableSnapshot any         // 永久不可兑换快照
	LedgerId                       any         // 账本流水标识
	CommitSequence                 any         // 提交序列
	ExchangedAt                    *gtime.Time // 兑换时间
	CreatedAt                      *gtime.Time // 创建时间
}
