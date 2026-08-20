// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsNotificationOutbox is the golang structure for table kids_notification_outbox.
type KidsNotificationOutbox struct {
	Id             uint64      `json:"id"             orm:"id"              description:"主键"`       // 主键
	NotificationId string      `json:"notificationId" orm:"notification_id" description:"接口通知标识"`   // 接口通知标识
	CircleId       string      `json:"circleId"       orm:"circle_id"       description:"接口圈子标识"`   // 接口圈子标识
	AccountId      string      `json:"accountId"      orm:"account_id"      description:"接收账号标识"`   // 接收账号标识
	ExchangeId     string      `json:"exchangeId"     orm:"exchange_id"     description:"接口兑换标识"`   // 接口兑换标识
	EventType      string      `json:"eventType"      orm:"event_type"      description:"通知事件类型"`   // 通知事件类型
	Payload        string      `json:"payload"        orm:"payload"         description:"通知载荷"`     // 通知载荷
	CommitSequence uint64      `json:"commitSequence" orm:"commit_sequence" description:"提交序列"`     // 提交序列
	Status         string      `json:"status"         orm:"status"          description:"投递状态"`     // 投递状态
	AttemptCount   uint        `json:"attemptCount"   orm:"attempt_count"   description:"已尝试投递次数"`  // 已尝试投递次数
	NextAttemptAt  *gtime.Time `json:"nextAttemptAt"  orm:"next_attempt_at" description:"下次允许投递时间"` // 下次允许投递时间
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:"创建时间"`     // 创建时间
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:"更新时间"`     // 更新时间
	Version        uint64      `json:"version"        orm:"version"         description:"版本号"`      // 版本号
}
