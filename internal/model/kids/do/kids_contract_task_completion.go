// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractTaskCompletion is the golang structure of table kids_contract_task_completion for DAO operations like Where/Data.
type KidsContractTaskCompletion struct {
	g.Meta         `orm:"table:kids_contract_task_completion, do:true"`
	Id             any         // 主键
	CompletionId   any         // 合同完成标识
	CircleId       any         // 合同圈子标识
	TaskId         any         // 合同任务标识
	MemberId       any         // 合同成员标识
	ScheduledDate  *gtime.Time // 预定日
	ZoneId         any         // 时区
	ProofAssetId   any         // 凭证资产标识
	TitleSnapshot  any         // 标题快照
	StarsSnapshot  any         // 星星快照
	CompletedBy    any         // 完成操作者快照
	CompletedAt    *gtime.Time // 完成时间
	CommitSequence any         // 提交序列
	Version        any         // 版本号
	CreatedAt      *gtime.Time // 创建时间
}
