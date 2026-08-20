// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardRedeemRecord is the golang structure of table kids_reward_redeem_record for DAO operations like Where/Data.
type KidsRewardRedeemRecord struct {
	g.Meta    `orm:"table:kids_reward_redeem_record, do:true"`
	Id        any         // 奖励兑换记录ID
	CircleId  any         // 圈子ID
	RewardId  any         // 奖励ID
	KidId     any         // 儿童成员ID
	UserId    any         // 兑换用户ID
	Title     any         // 奖励标题
	Icon      any         // 图标标识
	ImageUrl  any         // 奖励图片地址
	StarCost  any         // 消耗星星数量
	Remark    any         // 兑换备注
	CreatedAt *gtime.Time // 创建时间
}
