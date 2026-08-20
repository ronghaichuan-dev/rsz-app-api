// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractTaskAssignment is the golang structure of table kids_contract_task_assignment for DAO operations like Where/Data.
type KidsContractTaskAssignment struct {
	g.Meta    `orm:"table:kids_contract_task_assignment, do:true"`
	Id        any         // 主键
	TaskId    any         // 合同任务标识
	MemberId  any         // 合同成员标识
	CreatedAt *gtime.Time // 创建时间
	UpdatedAt *gtime.Time // 更新时间
}
