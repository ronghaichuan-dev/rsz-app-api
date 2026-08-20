// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsReward is the golang structure of table kids_reward for DAO operations like Where/Data.
type KidsReward struct {
	g.Meta             `orm:"table:kids_reward, do:true"`
	Id                 any         // 奖励ID
	CircleId           any         // 所属圈子ID
	Title              any         // 奖励标题
	Icon               any         // 图标标识
	ImageUrl           any         // 奖励图片地址
	StarCost           any         // 所需星星数量
	Stock              any         // 可用库存，-1表示不限量
	Description        any         // 奖励描述
	RepeatRule         any         // 重复兑换规则：none/daily/weekly/monthly/custom
	RepeatIntervalDays any         // 自定义重复兑换间隔天数
	CreatedAt          *gtime.Time // 创建时间
	UpdatedAt          *gtime.Time // 更新时间
	DeletedAt          *gtime.Time // 删除时间
}
