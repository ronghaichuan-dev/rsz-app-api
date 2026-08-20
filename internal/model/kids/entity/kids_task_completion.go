// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsTaskCompletion is the golang structure for table kids_task_completion.
type KidsTaskCompletion struct {
	Id             uint64      `json:"id"             orm:"id"              description:"主键"`      // 主键
	CompletionId   string      `json:"completionId"   orm:"completion_id"   description:"接口完成标识"`  // 接口完成标识
	CircleId       string      `json:"circleId"       orm:"circle_id"       description:"接口圈子标识"`  // 接口圈子标识
	TaskId         string      `json:"taskId"         orm:"task_id"         description:"接口任务标识"`  // 接口任务标识
	MemberId       string      `json:"memberId"       orm:"member_id"       description:"接口成员标识"`  // 接口成员标识
	ScheduledDate  *gtime.Time `json:"scheduledDate"  orm:"scheduled_date"  description:"预定日"`     // 预定日
	ZoneId         string      `json:"zoneId"         orm:"zone_id"         description:"时区"`      // 时区
	ProofAssetId   string      `json:"proofAssetId"   orm:"proof_asset_id"  description:"凭证资产标识"`  // 凭证资产标识
	TitleSnapshot  string      `json:"titleSnapshot"  orm:"title_snapshot"  description:"标题快照"`    // 标题快照
	StarsSnapshot  uint        `json:"starsSnapshot"  orm:"stars_snapshot"  description:"星星快照"`    // 星星快照
	CompletedBy    string      `json:"completedBy"    orm:"completed_by"    description:"完成操作者快照"` // 完成操作者快照
	CompletedAt    *gtime.Time `json:"completedAt"    orm:"completed_at"    description:"完成时间"`    // 完成时间
	CommitSequence uint64      `json:"commitSequence" orm:"commit_sequence" description:"提交序列"`    // 提交序列
	Version        uint64      `json:"version"        orm:"version"         description:"版本号"`     // 版本号
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`    // 创建时间
}
