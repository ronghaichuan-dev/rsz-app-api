// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskAssignee is the golang structure of table kids_task_assignee for DAO operations like Where/Data.
type KidsTaskAssignee struct {
	g.Meta        `orm:"table:kids_task_assignee, do:true"`
	Id            any         // 任务分配ID
	TaskId        any         // 任务ID
	KidId         any         // 儿童成员ID
	AssigneeOrder any         // 分配顺序，用于轮流模式
	Completed     any         // 该儿童是否已完成
	PhotoUrl      any         // 该儿童照片凭证地址
	CompletedAt   *gtime.Time // 该儿童完成时间
	CreatedAt     *gtime.Time // 创建时间
	UpdatedAt     *gtime.Time // 更新时间
}
