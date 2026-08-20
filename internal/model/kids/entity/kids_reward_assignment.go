// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardAssignment is the golang structure for table kids_reward_assignment.
type KidsRewardAssignment struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键"`     // 主键
	RewardId  string      `json:"rewardId"  orm:"reward_id"  description:"接口奖励标识"` // 接口奖励标识
	MemberId  string      `json:"memberId"  orm:"member_id"  description:"接口成员标识"` // 接口成员标识
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
}
