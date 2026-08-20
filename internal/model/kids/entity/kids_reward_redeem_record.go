// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsRewardRedeemRecord is the golang structure for table kids_reward_redeem_record.
type KidsRewardRedeemRecord struct {
	Id        uint64      `json:"id"        orm:"id"         description:"奖励兑换记录ID"` // 奖励兑换记录ID
	CircleId  uint64      `json:"circleId"  orm:"circle_id"  description:"圈子ID"`     // 圈子ID
	RewardId  uint64      `json:"rewardId"  orm:"reward_id"  description:"奖励ID"`     // 奖励ID
	KidId     uint64      `json:"kidId"     orm:"kid_id"     description:"儿童成员ID"`   // 儿童成员ID
	UserId    uint64      `json:"userId"    orm:"user_id"    description:"兑换用户ID"`   // 兑换用户ID
	Title     string      `json:"title"     orm:"title"      description:"奖励标题"`     // 奖励标题
	Icon      string      `json:"icon"      orm:"icon"       description:"图标标识"`     // 图标标识
	ImageUrl  string      `json:"imageUrl"  orm:"image_url"  description:"奖励图片地址"`   // 奖励图片地址
	StarCost  int         `json:"starCost"  orm:"star_cost"  description:"消耗星星数量"`   // 消耗星星数量
	Remark    string      `json:"remark"    orm:"remark"     description:"兑换备注"`     // 兑换备注
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`     // 创建时间
}
