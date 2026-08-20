// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsStarLedger is the golang structure of table kids_star_ledger for DAO operations like Where/Data.
type KidsStarLedger struct {
	g.Meta             `orm:"table:kids_star_ledger, do:true"`
	Id                 any         // 主键
	LedgerId           any         // 接口流水标识
	CircleId           any         // 接口圈子标识
	MemberId           any         // 接口成员标识
	Source             any         // 来源快照
	Delta              any         // 余额变化
	Reason             any         // 原因
	Actor              any         // 操作者快照
	ReversalOfLedgerId any         // 原流水标识
	CommitSequence     any         // 提交序列
	CreatedAt          *gtime.Time // 创建时间
}
