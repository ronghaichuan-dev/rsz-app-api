// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskCancellation is the golang structure of table kids_task_cancellation for DAO operations like Where/Data.
type KidsTaskCancellation struct {
	g.Meta         `orm:"table:kids_task_cancellation, do:true"`
	Id             any         // 主键
	CancellationId any         // 接口取消标识
	CompletionId   any         // 接口完成标识
	ReasonCode     any         // 取消原因
	CancelledBy    any         // 取消操作者快照
	CancelledAt    *gtime.Time // 取消时间
	CommitSequence any         // 提交序列
	CreatedAt      *gtime.Time // 创建时间
}
