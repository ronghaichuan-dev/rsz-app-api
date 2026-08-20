// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTask is the golang structure of table kids_task for DAO operations like Where/Data.
type KidsTask struct {
	g.Meta                `orm:"table:kids_task, do:true"`
	Id                    any         // 任务ID
	Title                 any         // 任务标题
	Icon                  any         // 图标标识
	Note                  any         // 任务说明
	Star                  any         // 星星奖励数量
	TaskDate              *gtime.Time // 计划日期
	CompletionMode        any         // 完成模式：single/rotation/anyone/everyone
	RepeatRule            any         // 重复规则
	RepeatEndType         any         // 重复结束类型：never/date/count
	RepeatEndDate         *gtime.Time // 重复结束日期
	RepeatEndCount        any         // 重复结束次数
	TimeLimitType         any         // 时间限制类型：all_day/range
	TimeLimitStart        any         // 开始时间，格式HH:mm
	TimeLimitEnd          any         // 结束时间，格式HH:mm
	ReminderType          any         // 提醒类型：none/at_time/before_start
	ReminderAt            any         // 提醒时间，格式HH:mm
	ReminderOffsetMinutes any         // 提前提醒分钟数
	NeedPhotoProof        any         // 是否需要照片凭证
	TagId                 any         // 标签ID
	Completed             any         // 是否已完成
	CompletedBy           any         // 完成任务的儿童成员ID
	PhotoUrl              any         // 照片凭证地址
	CompletedAt           *gtime.Time // 完成时间
	DeletedAt             *gtime.Time // 删除时间
	CreatedAt             *gtime.Time // 创建时间
	UpdatedAt             *gtime.Time // 更新时间
}
