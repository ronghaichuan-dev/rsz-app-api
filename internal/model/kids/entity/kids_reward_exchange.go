// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardExchange is the golang structure for table kids_reward_exchange.
type KidsRewardExchange struct {
	Id                             uint64      `json:"id"                             orm:"id"                               description:"主键"`       // 主键
	ExchangeId                     string      `json:"exchangeId"                     orm:"exchange_id"                      description:"接口兑换标识"`   // 接口兑换标识
	CircleId                       string      `json:"circleId"                       orm:"circle_id"                        description:"接口圈子标识"`   // 接口圈子标识
	MemberId                       string      `json:"memberId"                       orm:"member_id"                        description:"接口成员标识"`   // 接口成员标识
	MemberNameSnapshot             string      `json:"memberNameSnapshot"             orm:"member_name_snapshot"             description:"成员名称快照"`   // 成员名称快照
	MemberAvatarSnapshot           string      `json:"memberAvatarSnapshot"           orm:"member_avatar_snapshot"           description:"成员头像快照"`   // 成员头像快照
	RewardId                       string      `json:"rewardId"                       orm:"reward_id"                        description:"接口奖励标识"`   // 接口奖励标识
	RewardTitleSnapshot            string      `json:"rewardTitleSnapshot"            orm:"reward_title_snapshot"            description:"奖励标题快照"`   // 奖励标题快照
	RewardVisualSnapshot           string      `json:"rewardVisualSnapshot"           orm:"reward_visual_snapshot"           description:"奖励视觉快照"`   // 奖励视觉快照
	StarsDeductedSnapshot          uint        `json:"starsDeductedSnapshot"          orm:"stars_deducted_snapshot"          description:"扣减星星快照"`   // 扣减星星快照
	RewardRepeatRuleSnapshot       string      `json:"rewardRepeatRuleSnapshot"       orm:"reward_repeat_rule_snapshot"      description:"重复规则快照"`   // 重复规则快照
	RewardCooldownDaysSnapshot     uint        `json:"rewardCooldownDaysSnapshot"     orm:"reward_cooldown_days_snapshot"    description:"冷却天数快照"`   // 冷却天数快照
	CooldownUntilAtSnapshot        *gtime.Time `json:"cooldownUntilAtSnapshot"        orm:"cooldown_until_at_snapshot"       description:"冷却结束快照"`   // 冷却结束快照
	PermanentlyUnavailableSnapshot uint        `json:"permanentlyUnavailableSnapshot" orm:"permanently_unavailable_snapshot" description:"永久不可兑换快照"` // 永久不可兑换快照
	LedgerId                       string      `json:"ledgerId"                       orm:"ledger_id"                        description:"账本流水标识"`   // 账本流水标识
	CommitSequence                 uint64      `json:"commitSequence"                 orm:"commit_sequence"                  description:"提交序列"`     // 提交序列
	ExchangedAt                    *gtime.Time `json:"exchangedAt"                    orm:"exchanged_at"                     description:"兑换时间"`     // 兑换时间
	CreatedAt                      *gtime.Time `json:"createdAt"                      orm:"created_at"                       description:"创建时间"`     // 创建时间
}
