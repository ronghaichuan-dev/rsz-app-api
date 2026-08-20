// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskAssignee is the golang structure for table kids_task_assignee.
type KidsTaskAssignee struct {
	Id            uint64      `json:"id"            orm:"id"             description:"任务分配ID"`      // 任务分配ID
	TaskId        uint64      `json:"taskId"        orm:"task_id"        description:"任务ID"`        // 任务ID
	KidId         uint64      `json:"kidId"         orm:"kid_id"         description:"儿童成员ID"`      // 儿童成员ID
	AssigneeOrder int         `json:"assigneeOrder" orm:"assignee_order" description:"分配顺序，用于轮流模式"` // 分配顺序，用于轮流模式
	Completed     uint        `json:"completed"     orm:"completed"      description:"该儿童是否已完成"`    // 该儿童是否已完成
	PhotoUrl      string      `json:"photoUrl"      orm:"photo_url"      description:"该儿童照片凭证地址"`   // 该儿童照片凭证地址
	CompletedAt   *gtime.Time `json:"completedAt"   orm:"completed_at"   description:"该儿童完成时间"`     // 该儿童完成时间
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`        // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"更新时间"`        // 更新时间
}
