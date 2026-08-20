// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsCircleSelection is the golang structure for table kids_circle_selection.
type KidsCircleSelection struct {
	Id              uint64      `json:"id"              orm:"id"                description:"主键"`     // 主键
	SelectionId     string      `json:"selectionId"     orm:"selection_id"      description:"接口选择标识"` // 接口选择标识
	AccountId       string      `json:"accountId"       orm:"account_id"        description:"接口账号标识"` // 接口账号标识
	CurrentCircleId string      `json:"currentCircleId" orm:"current_circle_id" description:"当前圈子标识"` // 当前圈子标识
	Version         uint64      `json:"version"         orm:"version"           description:"版本号"`    // 版本号
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"        description:"创建时间"`   // 创建时间
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"        description:"更新时间"`   // 更新时间
}
