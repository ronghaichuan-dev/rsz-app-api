// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardPreset is the golang structure of table kids_reward_preset for DAO operations like Where/Data.
type KidsRewardPreset struct {
	g.Meta      `orm:"table:kids_reward_preset, do:true"`
	Id          any         // 奖励预设ID
	Title       any         // 预设奖励标题
	Icon        any         // 图标标识
	ImageUrl    any         // 奖励图片地址
	StarCost    any         // 默认所需星星数量
	Description any         // 预设描述
	Enabled     any         // 是否启用预设
	CreatedAt   *gtime.Time // 创建时间
	UpdatedAt   *gtime.Time // 更新时间
}
