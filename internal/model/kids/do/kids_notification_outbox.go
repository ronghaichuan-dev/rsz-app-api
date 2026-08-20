// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsNotificationOutbox is the golang structure of table kids_notification_outbox for DAO operations like Where/Data.
type KidsNotificationOutbox struct {
	g.Meta         `orm:"table:kids_notification_outbox, do:true"`
	Id             any         // 主键
	NotificationId any         // 接口通知标识
	CircleId       any         // 接口圈子标识
	AccountId      any         // 接收账号标识
	ExchangeId     any         // 接口兑换标识
	EventType      any         // 通知事件类型
	Payload        any         // 通知载荷
	CommitSequence any         // 提交序列
	Status         any         // 投递状态
	AttemptCount   any         // 已尝试投递次数
	NextAttemptAt  *gtime.Time // 下次允许投递时间
	CreatedAt      *gtime.Time // 创建时间
	UpdatedAt      *gtime.Time // 更新时间
	Version        any         // 版本号
}
