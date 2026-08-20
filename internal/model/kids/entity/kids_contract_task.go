// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractTask is the golang structure for table kids_contract_task.
type KidsContractTask struct {
	Id                      uint64      `json:"id"                      orm:"id"                         description:"主键"`       // 主键
	TaskId                  string      `json:"taskId"                  orm:"task_id"                    description:"合同任务标识"`   // 合同任务标识
	CircleId                string      `json:"circleId"                orm:"circle_id"                  description:"合同圈子标识"`   // 合同圈子标识
	Title                   string      `json:"title"                   orm:"title"                      description:"任务标题"`     // 任务标题
	Notes                   string      `json:"notes"                   orm:"notes"                      description:"任务备注"`     // 任务备注
	Emoji                   string      `json:"emoji"                   orm:"emoji"                      description:"任务图标"`     // 任务图标
	Stars                   uint        `json:"stars"                   orm:"stars"                      description:"星星数量"`     // 星星数量
	StartDate               *gtime.Time `json:"startDate"               orm:"start_date"                 description:"系列开始日期"`   // 系列开始日期
	ZoneId                  string      `json:"zoneId"                  orm:"zone_id"                    description:"时区"`       // 时区
	RepeatRule              string      `json:"repeatRule"              orm:"repeat_rule"                description:"重复规则"`     // 重复规则
	EndRule                 string      `json:"endRule"                 orm:"end_rule"                   description:"结束规则"`     // 结束规则
	TimeLimitMinuteOfDay    uint        `json:"timeLimitMinuteOfDay"    orm:"time_limit_minute_of_day"   description:"时限分钟"`     // 时限分钟
	ReminderConfig          string      `json:"reminderConfig"          orm:"reminder_config"            description:"提醒配置"`     // 提醒配置
	PhotoRequired           uint        `json:"photoRequired"           orm:"photo_required"             description:"是否要求图片凭证"` // 是否要求图片凭证
	TaskTagId               string      `json:"taskTagId"               orm:"task_tag_id"                description:"合同任务标签标识"` // 合同任务标签标识
	SeriesRevision          uint64      `json:"seriesRevision"          orm:"series_revision"            description:"系列修订号"`    // 系列修订号
	FutureEffectiveFromDate *gtime.Time `json:"futureEffectiveFromDate" orm:"future_effective_from_date" description:"当前未来生效日"`  // 当前未来生效日
	Status                  string      `json:"status"                  orm:"status"                     description:"任务状态"`     // 任务状态
	Version                 uint64      `json:"version"                 orm:"version"                    description:"版本号"`      // 版本号
	DeletedAt               *gtime.Time `json:"deletedAt"               orm:"deleted_at"                 description:"删除时间"`     // 删除时间
	CreatedAt               *gtime.Time `json:"createdAt"               orm:"created_at"                 description:"创建时间"`     // 创建时间
	UpdatedAt               *gtime.Time `json:"updatedAt"               orm:"updated_at"                 description:"更新时间"`     // 更新时间
}
