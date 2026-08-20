// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsNotification is the golang structure of table kids_notification for DAO operations like Where/Data.
type KidsNotification struct {
	g.Meta           `orm:"table:kids_notification, do:true"`
	Id               any         // 通知ID
	MemberId         any         // 目标成员ID，0表示全家庭
	NotificationType any         // 通知类型
	Title            any         // 通知标题
	Content          any         // 通知内容
	IsRead           any         // 是否已读
	CreatedAt        *gtime.Time // 创建时间
	UpdatedAt        *gtime.Time // 更新时间
}
