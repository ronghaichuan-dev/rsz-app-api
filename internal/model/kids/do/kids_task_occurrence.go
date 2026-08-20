// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskOccurrence is the golang structure of table kids_task_occurrence for DAO operations like Where/Data.
type KidsTaskOccurrence struct {
	g.Meta                `orm:"table:kids_task_occurrence, do:true"`
	Id                    any         // 内部主键
	CircleId              any         // 接口圈子标识
	TaskId                any         // 接口任务标识
	MemberId              any         // 接口成员标识
	ScheduledDate         *gtime.Time // 预定日
	ZoneId                any         // 时区
	DefinitionRevision    any         // 定义修订号
	TitleSnapshot         any         // 标题快照
	NotesSnapshot         any         // 备注快照
	EmojiSnapshot         any         // 图标快照
	StarsSnapshot         any         // 星星快照
	PhotoRequiredSnapshot any         // 图片要求快照
	TaskTagIdSnapshot     any         // 标签快照
	State                 any         // occurrence 状态
	CompletionId          any         // 完成事实标识
	Version               any         // 版本号
	CreatedAt             *gtime.Time // 创建时间
	UpdatedAt             *gtime.Time // 更新时间
}
