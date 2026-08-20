// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskDefinition is the golang structure of table kids_task_definition for DAO operations like Where/Data.
type KidsTaskDefinition struct {
	g.Meta                  `orm:"table:kids_task_definition, do:true"`
	Id                      any         // 主键
	TaskId                  any         // 接口任务标识
	CircleId                any         // 接口圈子标识
	Title                   any         // 任务标题
	Notes                   any         // 任务备注
	Emoji                   any         // 任务图标
	Stars                   any         // 星星数量
	StartDate               *gtime.Time // 系列开始日期
	ZoneId                  any         // 时区
	RepeatRule              any         // 重复规则
	EndRule                 any         // 结束规则
	TimeLimitMinuteOfDay    any         // 时限分钟
	ReminderConfig          any         // 提醒配置
	PhotoRequired           any         // 是否要求图片凭证
	TaskTagId               any         // 接口任务标签标识
	SeriesRevision          any         // 系列修订号
	FutureEffectiveFromDate *gtime.Time // 当前未来生效日
	Status                  any         // 任务状态
	Version                 any         // 版本号
	DeletedAt               *gtime.Time // 删除时间
	CreatedAt               *gtime.Time // 创建时间
	UpdatedAt               *gtime.Time // 更新时间
}
