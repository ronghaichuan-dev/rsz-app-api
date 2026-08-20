// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// KidsDeviceNotification is the golang structure for table kids_device_notification.
type KidsDeviceNotification struct {
	Id            uint64      `json:"id"            orm:"id"             description:"设备通知ID"`             // 设备通知ID
	UserId        uint64      `json:"userId"        orm:"user_id"        description:"用户ID"`               // 用户ID
	DeviceId      string      `json:"deviceId"      orm:"device_id"      description:"设备ID"`               // 设备ID
	Platform      string      `json:"platform"      orm:"platform"       description:"平台：ios/android/web"` // 平台：ios/android/web
	DeviceToken   string      `json:"deviceToken"   orm:"device_token"   description:"推送设备令牌"`             // 推送设备令牌
	Authorized    uint        `json:"authorized"    orm:"authorized"     description:"是否已授权通知"`            // 是否已授权通知
	TaskEnabled   uint        `json:"taskEnabled"   orm:"task_enabled"   description:"是否开启任务提醒"`           // 是否开启任务提醒
	RewardEnabled uint        `json:"rewardEnabled" orm:"reward_enabled" description:"是否开启奖励提醒"`           // 是否开启奖励提醒
	MemberEnabled uint        `json:"memberEnabled" orm:"member_enabled" description:"是否开启成员提醒"`           // 是否开启成员提醒
	CreatedAt     *gtime.Time `json:"createdAt"     orm:"created_at"     description:"创建时间"`               // 创建时间
	UpdatedAt     *gtime.Time `json:"updatedAt"     orm:"updated_at"     description:"更新时间"`               // 更新时间
}
