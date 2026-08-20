// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsNotification is the golang structure for table kids_notification.
type KidsNotification struct {
	Id               uint64      `json:"id"               orm:"id"                description:"通知ID"`          // 通知ID
	MemberId         uint64      `json:"memberId"         orm:"member_id"         description:"目标成员ID，0表示全家庭"` // 目标成员ID，0表示全家庭
	NotificationType string      `json:"notificationType" orm:"notification_type" description:"通知类型"`          // 通知类型
	Title            string      `json:"title"            orm:"title"             description:"通知标题"`          // 通知标题
	Content          string      `json:"content"          orm:"content"           description:"通知内容"`          // 通知内容
	IsRead           uint        `json:"isRead"           orm:"is_read"           description:"是否已读"`          // 是否已读
	CreatedAt        *gtime.Time `json:"createdAt"        orm:"created_at"        description:"创建时间"`          // 创建时间
	UpdatedAt        *gtime.Time `json:"updatedAt"        orm:"updated_at"        description:"更新时间"`          // 更新时间
}
