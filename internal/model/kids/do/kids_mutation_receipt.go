// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsMutationReceipt is the golang structure of table kids_mutation_receipt for DAO operations like Where/Data.
type KidsMutationReceipt struct {
	g.Meta      `orm:"table:kids_mutation_receipt, do:true"`
	Id          any         // 主键
	ReceiptId   any         // 接口回执标识
	CommitId    any         // 提交标识
	OperationId any         // 接口操作标识
	ResultKind  any         // 结果类型
	CommittedAt *gtime.Time // 提交时间
	CreatedAt   *gtime.Time // 创建时间
}
