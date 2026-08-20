// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTask is the golang structure for table kids_task.
type KidsTask struct {
	Id                    uint64      `json:"id"                    orm:"id"                      description:"任务ID"`                                 // 任务ID
	Title                 string      `json:"title"                 orm:"title"                   description:"任务标题"`                                 // 任务标题
	Icon                  string      `json:"icon"                  orm:"icon"                    description:"图标标识"`                                 // 图标标识
	Note                  string      `json:"note"                  orm:"note"                    description:"任务说明"`                                 // 任务说明
	Star                  int         `json:"star"                  orm:"star"                    description:"星星奖励数量"`                               // 星星奖励数量
	TaskDate              *gtime.Time `json:"taskDate"              orm:"task_date"               description:"计划日期"`                                 // 计划日期
	CompletionMode        string      `json:"completionMode"        orm:"completion_mode"         description:"完成模式：single/rotation/anyone/everyone"` // 完成模式：single/rotation/anyone/everyone
	RepeatRule            string      `json:"repeatRule"            orm:"repeat_rule"             description:"重复规则"`                                 // 重复规则
	RepeatEndType         string      `json:"repeatEndType"         orm:"repeat_end_type"         description:"重复结束类型：never/date/count"`              // 重复结束类型：never/date/count
	RepeatEndDate         *gtime.Time `json:"repeatEndDate"         orm:"repeat_end_date"         description:"重复结束日期"`                               // 重复结束日期
	RepeatEndCount        int         `json:"repeatEndCount"        orm:"repeat_end_count"        description:"重复结束次数"`                               // 重复结束次数
	TimeLimitType         string      `json:"timeLimitType"         orm:"time_limit_type"         description:"时间限制类型：all_day/range"`                 // 时间限制类型：all_day/range
	TimeLimitStart        string      `json:"timeLimitStart"        orm:"time_limit_start"        description:"开始时间，格式HH:mm"`                         // 开始时间，格式HH:mm
	TimeLimitEnd          string      `json:"timeLimitEnd"          orm:"time_limit_end"          description:"结束时间，格式HH:mm"`                         // 结束时间，格式HH:mm
	ReminderType          string      `json:"reminderType"          orm:"reminder_type"           description:"提醒类型：none/at_time/before_start"`       // 提醒类型：none/at_time/before_start
	ReminderAt            string      `json:"reminderAt"            orm:"reminder_at"             description:"提醒时间，格式HH:mm"`                         // 提醒时间，格式HH:mm
	ReminderOffsetMinutes int         `json:"reminderOffsetMinutes" orm:"reminder_offset_minutes" description:"提前提醒分钟数"`                              // 提前提醒分钟数
	NeedPhotoProof        uint        `json:"needPhotoProof"        orm:"need_photo_proof"        description:"是否需要照片凭证"`                             // 是否需要照片凭证
	TagId                 uint64      `json:"tagId"                 orm:"tag_id"                  description:"标签ID"`                                 // 标签ID
	Completed             uint        `json:"completed"             orm:"completed"               description:"是否已完成"`                                // 是否已完成
	CompletedBy           uint64      `json:"completedBy"           orm:"completed_by"            description:"完成任务的儿童成员ID"`                          // 完成任务的儿童成员ID
	PhotoUrl              string      `json:"photoUrl"              orm:"photo_url"               description:"照片凭证地址"`                               // 照片凭证地址
	CompletedAt           *gtime.Time `json:"completedAt"           orm:"completed_at"            description:"完成时间"`                                 // 完成时间
	DeletedAt             *gtime.Time `json:"deletedAt"             orm:"deleted_at"              description:"删除时间"`                                 // 删除时间
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              description:"创建时间"`                                 // 创建时间
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"              description:"更新时间"`                                 // 更新时间
}
