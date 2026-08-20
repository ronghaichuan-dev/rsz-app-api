// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardAssignee is the golang structure of table kids_reward_assignee for DAO operations like Where/Data.
type KidsRewardAssignee struct {
	g.Meta    `orm:"table:kids_reward_assignee, do:true"`
	Id        any         // 奖励指派ID
	RewardId  any         // 奖励ID
	KidId     any         // 儿童成员ID
	CreatedAt *gtime.Time // 创建时间
}
