// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardAssignment is the golang structure of table kids_reward_assignment for DAO operations like Where/Data.
type KidsRewardAssignment struct {
	g.Meta    `orm:"table:kids_reward_assignment, do:true"`
	Id        any         // 主键
	RewardId  any         // 接口奖励标识
	MemberId  any         // 接口成员标识
	CreatedAt *gtime.Time // 创建时间
}
