// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsStarLedger is the golang structure for table kids_star_ledger.
type KidsStarLedger struct {
	Id                 uint64      `json:"id"                 orm:"id"                    description:"主键"`     // 主键
	LedgerId           string      `json:"ledgerId"           orm:"ledger_id"             description:"接口流水标识"` // 接口流水标识
	CircleId           string      `json:"circleId"           orm:"circle_id"             description:"接口圈子标识"` // 接口圈子标识
	MemberId           string      `json:"memberId"           orm:"member_id"             description:"接口成员标识"` // 接口成员标识
	Source             string      `json:"source"             orm:"source"                description:"来源快照"`   // 来源快照
	Delta              int         `json:"delta"              orm:"delta"                 description:"余额变化"`   // 余额变化
	Reason             string      `json:"reason"             orm:"reason"                description:"原因"`     // 原因
	Actor              string      `json:"actor"              orm:"actor"                 description:"操作者快照"`  // 操作者快照
	ReversalOfLedgerId string      `json:"reversalOfLedgerId" orm:"reversal_of_ledger_id" description:"原流水标识"`  // 原流水标识
	CommitSequence     uint64      `json:"commitSequence"     orm:"commit_sequence"       description:"提交序列"`   // 提交序列
	CreatedAt          *gtime.Time `json:"createdAt"          orm:"created_at"            description:"创建时间"`   // 创建时间
}
