// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsMutationReceipt is the golang structure for table kids_mutation_receipt.
type KidsMutationReceipt struct {
	Id          uint64      `json:"id"          orm:"id"           description:"主键"`     // 主键
	ReceiptId   string      `json:"receiptId"   orm:"receipt_id"   description:"接口回执标识"` // 接口回执标识
	CommitId    string      `json:"commitId"    orm:"commit_id"    description:"提交标识"`   // 提交标识
	OperationId string      `json:"operationId" orm:"operation_id" description:"接口操作标识"` // 接口操作标识
	ResultKind  string      `json:"resultKind"  orm:"result_kind"  description:"结果类型"`   // 结果类型
	CommittedAt *gtime.Time `json:"committedAt" orm:"committed_at" description:"提交时间"`   // 提交时间
	CreatedAt   *gtime.Time `json:"createdAt"   orm:"created_at"   description:"创建时间"`   // 创建时间
}
