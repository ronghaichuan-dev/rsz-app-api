// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsContractTaskOccurrence is the golang structure for table kids_contract_task_occurrence.
type KidsContractTaskOccurrence struct {
	Id                    uint64      `json:"id"                    orm:"id"                      description:"内部主键"`          // 内部主键
	CircleId              string      `json:"circleId"              orm:"circle_id"               description:"合同圈子标识"`        // 合同圈子标识
	TaskId                string      `json:"taskId"                orm:"task_id"                 description:"合同任务标识"`        // 合同任务标识
	MemberId              string      `json:"memberId"              orm:"member_id"               description:"合同成员标识"`        // 合同成员标识
	ScheduledDate         *gtime.Time `json:"scheduledDate"         orm:"scheduled_date"          description:"预定日"`           // 预定日
	ZoneId                string      `json:"zoneId"                orm:"zone_id"                 description:"时区"`            // 时区
	DefinitionRevision    uint64      `json:"definitionRevision"    orm:"definition_revision"     description:"定义修订号"`         // 定义修订号
	TitleSnapshot         string      `json:"titleSnapshot"         orm:"title_snapshot"          description:"标题快照"`          // 标题快照
	NotesSnapshot         string      `json:"notesSnapshot"         orm:"notes_snapshot"          description:"备注快照"`          // 备注快照
	EmojiSnapshot         string      `json:"emojiSnapshot"         orm:"emoji_snapshot"          description:"图标快照"`          // 图标快照
	StarsSnapshot         uint        `json:"starsSnapshot"         orm:"stars_snapshot"          description:"星星快照"`          // 星星快照
	PhotoRequiredSnapshot uint        `json:"photoRequiredSnapshot" orm:"photo_required_snapshot" description:"图片要求快照"`        // 图片要求快照
	TaskTagIdSnapshot     string      `json:"taskTagIdSnapshot"     orm:"task_tag_id_snapshot"    description:"标签快照"`          // 标签快照
	State                 string      `json:"state"                 orm:"state"                   description:"occurrence 状态"` // occurrence 状态
	CompletionId          string      `json:"completionId"          orm:"completion_id"           description:"完成事实标识"`        // 完成事实标识
	Version               uint64      `json:"version"               orm:"version"                 description:"版本号"`           // 版本号
	CreatedAt             *gtime.Time `json:"createdAt"             orm:"created_at"              description:"创建时间"`          // 创建时间
	UpdatedAt             *gtime.Time `json:"updatedAt"             orm:"updated_at"              description:"更新时间"`          // 更新时间
}
