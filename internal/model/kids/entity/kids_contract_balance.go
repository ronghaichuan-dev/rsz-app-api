// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractBalance is the golang structure for table kids_contract_balance.
type KidsContractBalance struct {
	Id                   uint64      `json:"id"                   orm:"id"                     description:"主键"`     // 主键
	CircleId             string      `json:"circleId"             orm:"circle_id"              description:"合同圈子标识"` // 合同圈子标识
	MemberId             string      `json:"memberId"             orm:"member_id"              description:"合同成员标识"` // 合同成员标识
	Balance              int         `json:"balance"              orm:"balance"                description:"星星余额"`   // 星星余额
	Version              uint64      `json:"version"              orm:"version"                description:"版本号"`    // 版本号
	SourceCommitId       string      `json:"sourceCommitId"       orm:"source_commit_id"       description:"来源提交标识"` // 来源提交标识
	SourceCommitSequence uint64      `json:"sourceCommitSequence" orm:"source_commit_sequence" description:"来源提交序列"` // 来源提交序列
	UpdatedAt            *gtime.Time `json:"updatedAt"            orm:"updated_at"             description:"更新时间"`   // 更新时间
}
