// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractSequence is the golang structure for table kids_contract_sequence.
type KidsContractSequence struct {
	Id                 uint        `json:"id"                 orm:"id"                   description:"固定序列标识"`  // 固定序列标识
	NextCommitSequence uint64      `json:"nextCommitSequence" orm:"next_commit_sequence" description:"下一个提交序列"` // 下一个提交序列
	UpdatedAt          *gtime.Time `json:"updatedAt"          orm:"updated_at"           description:"更新时间"`    // 更新时间
}
