// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractTaskCancellation is the golang structure for table kids_contract_task_cancellation.
type KidsContractTaskCancellation struct {
	Id             uint64      `json:"id"             orm:"id"              description:"主键"`      // 主键
	CancellationId string      `json:"cancellationId" orm:"cancellation_id" description:"合同取消标识"`  // 合同取消标识
	CompletionId   string      `json:"completionId"   orm:"completion_id"   description:"合同完成标识"`  // 合同完成标识
	ReasonCode     string      `json:"reasonCode"     orm:"reason_code"     description:"取消原因"`    // 取消原因
	CancelledBy    string      `json:"cancelledBy"    orm:"cancelled_by"    description:"取消操作者快照"` // 取消操作者快照
	CancelledAt    *gtime.Time `json:"cancelledAt"    orm:"cancelled_at"    description:"取消时间"`    // 取消时间
	CommitSequence uint64      `json:"commitSequence" orm:"commit_sequence" description:"提交序列"`    // 提交序列
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`    // 创建时间
}
