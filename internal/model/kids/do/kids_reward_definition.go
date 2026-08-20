// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardDefinition is the golang structure of table kids_reward_definition for DAO operations like Where/Data.
type KidsRewardDefinition struct {
	g.Meta        `orm:"table:kids_reward_definition, do:true"`
	Id            any         // 主键
	RewardId      any         // 接口奖励标识
	CircleId      any         // 接口圈子标识
	Title         any         // 奖励标题
	Description   any         // 奖励描述
	Visual        any         // 奖励视觉引用
	StarsRequired any         // 兑换所需星星
	RepeatRule    any         // 重复兑换规则
	CooldownDays  any         // 自定义冷却天数
	ZoneId        any         // 时区标识
	Status        any         // 奖励状态
	Version       any         // 版本号
	DeletedAt     *gtime.Time // 删除时间
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
