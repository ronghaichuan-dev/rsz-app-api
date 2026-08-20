// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardAssignee is the golang structure for table kids_reward_assignee.
type KidsRewardAssignee struct {
	Id        uint64      `json:"id"        orm:"id"         description:"奖励指派ID"` // 奖励指派ID
	RewardId  uint64      `json:"rewardId"  orm:"reward_id"  description:"奖励ID"`   // 奖励ID
	KidId     uint64      `json:"kidId"     orm:"kid_id"     description:"儿童成员ID"` // 儿童成员ID
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
}
