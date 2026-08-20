// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsSyncCommit is the golang structure for table kids_sync_commit.
type KidsSyncCommit struct {
	Id             uint64      `json:"id"             orm:"id"              description:"主键"`     // 主键
	CommitId       string      `json:"commitId"       orm:"commit_id"       description:"接口提交标识"` // 接口提交标识
	CircleId       string      `json:"circleId"       orm:"circle_id"       description:"圈子标识"`   // 圈子标识
	CommitSequence uint64      `json:"commitSequence" orm:"commit_sequence" description:"单调提交序列"` // 单调提交序列
	ChangePayload  string      `json:"changePayload"  orm:"change_payload"  description:"完整变更集合"` // 完整变更集合
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`   // 创建时间
}
