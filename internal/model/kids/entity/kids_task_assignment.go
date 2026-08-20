// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskAssignment is the golang structure for table kids_task_assignment.
type KidsTaskAssignment struct {
	Id        uint64      `json:"id"        orm:"id"         description:"主键"`     // 主键
	TaskId    string      `json:"taskId"    orm:"task_id"    description:"接口任务标识"` // 接口任务标识
	MemberId  string      `json:"memberId"  orm:"member_id"  description:"接口成员标识"` // 接口成员标识
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:"创建时间"`   // 创建时间
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:"更新时间"`   // 更新时间
}
