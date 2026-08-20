// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractReceipt is the golang structure of table kids_contract_receipt for DAO operations like Where/Data.
type KidsContractReceipt struct {
	g.Meta      `orm:"table:kids_contract_receipt, do:true"`
	Id          any         // 主键
	ReceiptId   any         // 合同回执标识
	CommitId    any         // 提交标识
	OperationId any         // 合同操作标识
	ResultKind  any         // 结果类型
	CommittedAt *gtime.Time // 提交时间
	CreatedAt   *gtime.Time // 创建时间
}
